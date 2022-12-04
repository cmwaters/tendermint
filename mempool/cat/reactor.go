package cat

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/p2p"
	protomem "github.com/tendermint/tendermint/proto/tendermint/mempool"
	"github.com/tendermint/tendermint/types"
)

const (
	maxStateChannelSize = tmhash.Size + tmhash.TruncatedSize + 10 // tx_key + node_id + buffer (for proto encoding)
	MempoolStateChannel = byte(0x31)
)

// Reactor handles mempool tx broadcasting amongst peers.
// It maintains a map from peer ID to counter, to prevent gossiping txs to the
// peers you received it from.
type Reactor struct {
	p2p.BaseReactor
	config  *cfg.MempoolConfig
	mempool *TxPool
	ids     *mempoolIDs
}

// NewReactor returns a new Reactor with the given config and mempool.
func NewReactor(config *cfg.MempoolConfig, mempool *TxPool) *Reactor {
	memR := &Reactor{
		config:  config,
		mempool: mempool,
		ids:     newMempoolIDs(),
	}
	memR.BaseReactor = *p2p.NewBaseReactor("Mempool", memR)
	return memR
}

// InitPeer implements Reactor by creating a state for the peer.
func (memR *Reactor) InitPeer(peer p2p.Peer) p2p.Peer {
	memR.ids.ReserveForPeer(peer)
	return peer
}

// SetLogger sets the Logger on the reactor and the underlying mempool.
func (memR *Reactor) SetLogger(l log.Logger) {
	memR.Logger = l
}

// OnStart implements p2p.BaseReactor.
func (memR *Reactor) OnStart() error {
	if !memR.config.Broadcast {
		memR.Logger.Info("Tx broadcasting is disabled")
	}
	go func() {
		for {
			select {
			case <-memR.Quit():
				return

			// any newly verified tx via RFC, we immediately broadcast to everyone
			case tx := <-memR.mempool.outboundTxs:
				memR.broadcastNewTx(tx)
			}
		}
	}()
	return nil
}

// GetChannels implements Reactor by returning the list of channels for this
// reactor.
func (memR *Reactor) GetChannels() []*p2p.ChannelDescriptor {
	largestTx := make([]byte, memR.config.MaxTxBytes)
	batchMsg := protomem.Message{
		Sum: &protomem.Message_Txs{
			Txs: &protomem.Txs{Txs: [][]byte{largestTx}},
		},
	}

	return []*p2p.ChannelDescriptor{
		{
			ID:                  mempool.MempoolChannel,
			Priority:            6,
			RecvMessageCapacity: batchMsg.Size(),
			MessageType:         &protomem.Message{},
		},
		{
			ID:                  MempoolStateChannel,
			Priority:            5,
			RecvMessageCapacity: maxStateChannelSize,
			MessageType:         &protomem.Message{},
		},
	}
}

// AddPeer implements Reactor.
// It starts a broadcast routine ensuring all txs are forwarded to the given peer.
func (memR *Reactor) AddPeer(peer p2p.Peer) {
	if memR.config.Broadcast {
		memR.sendAllTxKeys(peer)
	}
}

// RemovePeer implements Reactor.
func (memR *Reactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	memR.ids.Reclaim(peer)
	// broadcast routine checks if peer is gone and returns
}

// Receive implements Reactor.
// It adds any received transactions to the mempool.
func (memR *Reactor) ReceiveEnvelope(e p2p.Envelope) {
	memR.Logger.Debug("Receive", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)
	switch msg := e.Message.(type) {
	case *protomem.Txs:
		protoTxs := msg.GetTxs()
		if len(protoTxs) == 0 {
			memR.Logger.Error("received tmpty txs from peer", "src", e.Src)
			return
		}
		txInfo := mempool.TxInfo{SenderID: memR.ids.GetForPeer(e.Src)}
		if e.Src != nil {
			txInfo.SenderP2PID = e.Src.ID()
		}

		var err error
		for _, tx := range protoTxs {
			ntx := types.Tx(tx)
			key := ntx.Key()
			if memR.mempool.IsRejectedTx(key) {
				// The peer has sent us a transaction that we have already rejected. Since `CheckTx` can
				// be non-deterministic, we don't punish the peer but instead just ignore the msg
				continue
			}
			if memR.mempool.WasRecentlyEvicted(key) {
				// the transaction was recently evicted. If true, we attempt to re-add it to the mempool
				// skipping check tx.
				err := memR.mempool.TryReinsertEvictedTx(key, ntx, txInfo.SenderID)
				if err != nil {
					memR.Logger.Info("Unable to readd evicted tx", "tx_key", key, "err", err)
				}
				continue
			}
			if memR.mempool.Has(key) {
				// We have already received this transaction. We mark the peer that sent the message
				// as already seeing the transaction as well and then we finish
				memR.mempool.PeerHasTx(txInfo.SenderID, key)
				continue
			}
			memR.broadcastSeenTx(key)
			_, err = memR.mempool.TryAddNewTx(ntx, key, txInfo)
			if err != nil {
				memR.Logger.Info("Could not add tx", "tx_key", key, "err", err)
			}
		}
	case *protomem.SeenTx:
		if len(msg.TxKey) != types.TxKeySize {
			memR.Logger.Error("Peer sent SeenTx with incorrect key size", "len", len(msg.TxKey))
			return
		}
		var txKey [types.TxKeySize]byte
		copy(txKey[:], msg.TxKey)
		memR.mempool.PeerHasTx(memR.ids.GetForPeer(e.Src), types.TxKey(txKey))

	default:
		memR.Logger.Error("unknown message type", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)
		memR.Switch.StopPeerForError(e.Src, fmt.Errorf("mempool cannot handle message of type: %T", e.Message))
		return
	}

	// broadcasting happens from go routines per peer
}

func (memR *Reactor) Receive(chID byte, peer p2p.Peer, msgBytes []byte) {
	msg := &protomem.Message{}
	err := proto.Unmarshal(msgBytes, msg)
	if err != nil {
		panic(err)
	}
	uw, err := msg.Unwrap()
	if err != nil {
		panic(err)
	}
	memR.ReceiveEnvelope(p2p.Envelope{
		ChannelID: chID,
		Src:       peer,
		Message:   uw,
	})
}

// PeerState describes the state of a peer.
type PeerState interface {
	GetHeight() int64
}

func (memR *Reactor) broadcastSeenTx(txKey types.TxKey) {
	memR.Logger.Debug("broadcasting seen tx", "tx_key", txKey)
	memR.Switch.BroadcastEnvelope(p2p.Envelope{
		ChannelID: MempoolStateChannel,
		Message:   &protomem.SeenTx{TxKey: txKey[:]},
	})
}

func (memR *Reactor) broadcastNewTx(tx types.Tx) {
	memR.Logger.Debug("broadcasting new tx to all peers")
	memR.Switch.BroadcastEnvelope(p2p.Envelope{
		ChannelID: mempool.MempoolChannel,
		Message:   &protomem.Txs{Txs: [][]byte{tx}},
	})
}

// sendAllTxKeys loops through all txs currently in the mempool and iteratively
// sends a `SeenTx` message to the peer. This is added to a queue and will block
// when the queue becomes full.
func (memR *Reactor) sendAllTxKeys(peer p2p.Peer) {
	txKeys := memR.mempool.store.getAllKeys()
	for _, txKey := range txKeys {
		p2p.SendEnvelopeShim(peer, p2p.Envelope{ //nolint: staticcheck
			ChannelID: mempool.MempoolChannel,
			Message:   &protomem.SeenTx{TxKey: txKey[:]},
		}, memR.Logger)
	}
}
