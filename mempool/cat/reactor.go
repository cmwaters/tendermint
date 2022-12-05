package cat

import (
	"fmt"
	"time"

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
	// tx_key + node_id + buffer (for proto encoding)
	maxStateChannelSize = tmhash.Size + tmhash.TruncatedSize + 10

	defaultGossipDelay = 100 * time.Millisecond

	// Content Addressable Tx Pool gossips state based messages (SeenTx and WantTx) on a separate channel
	// for cross compatibility
	MempoolStateChannel = byte(0x31)
)

// Reactor handles mempool tx broadcasting logic amongst peers. For the main
// logic behind the protocol, refer to `ReceiveEnvelope` or to the english
// spec under /.spec.md
type Reactor struct {
	p2p.BaseReactor
	opts     *ReactorOptions
	mempool  *TxPool
	ids      *mempoolIDs
	requests *requestScheduler
}

type ReactorOptions struct {
	// ListenOnly means that the node will never broadcast any of the transactions that
	// it receives. This is useful for keeping transactions private
	ListenOnly bool

	// MaxTxSize is the maximum size of a transaction that can be received
	MaxTxSize int

	// MaxGossipDelay is the maximum allotted time that the reactor expects a transaction to
	// arrive before issuing a new request to a different peer
	MaxGossipDelay time.Duration
}

func (opts *ReactorOptions) VerifyAndComplete() error {
	if opts.MaxTxSize == 0 {
		opts.MaxTxSize = cfg.DefaultMempoolConfig().MaxTxBytes
	}

	if opts.MaxGossipDelay == 0 {
		opts.MaxGossipDelay = defaultGossipDelay
	}

	if opts.MaxTxSize < 0 {
		return fmt.Errorf("max tx size (%d) cannot be negative", opts.MaxTxSize)
	}

	if opts.MaxGossipDelay < 0 {
		return fmt.Errorf("max gossip delay (%d) cannot be negative", opts.MaxGossipDelay)
	}

	return nil
}

// NewReactor returns a new Reactor with the given config and mempool.
func NewReactor(mempool *TxPool, opts *ReactorOptions) (*Reactor, error) {
	err := opts.VerifyAndComplete()
	if err != nil {
		return nil, err
	}
	memR := &Reactor{
		opts:     opts,
		mempool:  mempool,
		ids:      newMempoolIDs(),
		requests: newRequestScheduler(opts.MaxGossipDelay),
	}
	memR.BaseReactor = *p2p.NewBaseReactor("Mempool", memR)
	return memR, nil
}

// SetLogger sets the Logger on the reactor and the underlying mempool.
func (memR *Reactor) SetLogger(l log.Logger) {
	memR.Logger = l
}

// OnStart implements p2p.BaseReactor.
func (memR *Reactor) OnStart() error {
	if memR.opts.ListenOnly {
		memR.Logger.Info("Tx broadcasting is disabled")
		return nil
	}
	go func() {
		for {
			select {
			case <-memR.Quit():
				return

			// listen in for any newly verified tx via RFC, then immediately
			// broadcast it to peers.
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
	largestTx := make([]byte, memR.opts.MaxTxSize)
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

// InitPeer implements Reactor by creating a state for the peer.
func (memR *Reactor) InitPeer(peer p2p.Peer) p2p.Peer {
	memR.ids.ReserveForPeer(peer)
	return peer
}

// AddPeer implements Reactor.
// It starts a broadcast routine ensuring all txs are forwarded to the given peer.
func (memR *Reactor) AddPeer(peer p2p.Peer) {
	if !memR.opts.ListenOnly {
		memR.sendAllTxKeys(peer)
	}
}

// RemovePeer implements Reactor.
func (memR *Reactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	memR.ids.Reclaim(peer.ID())
	// broadcast routine checks if peer is gone and returns
}

// ReceiveEnvelope implements Reactor.
// It adds any received transactions to the mempool.
func (memR *Reactor) ReceiveEnvelope(e p2p.Envelope) {
	memR.Logger.Debug("Receive", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)
	switch msg := e.Message.(type) {

	// A peer has sent us one or more transactions. This could be either because we requested them
	// or because the peer received a new transaction and is broadcasting it to us.
	// NOTE: This setup also means that we can support older mempool implementations that simply
	// flooded the network with transactions.
	case *protomem.Txs:
		protoTxs := msg.GetTxs()
		if len(protoTxs) == 0 {
			memR.Logger.Error("received tmpty txs from peer", "src", e.Src)
			return
		}
		peerID := memR.ids.GetIDForPeer(e.Src.ID())
		txInfo := mempool.TxInfo{SenderID: peerID}
		if e.Src != nil {
			txInfo.SenderP2PID = e.Src.ID()
		}

		var err error
		for _, tx := range protoTxs {
			ntx := types.Tx(tx)
			key := ntx.Key()
			// If we requested the transaction we mark it as received.
			if memR.requests.From(peerID).Includes(key) {
				memR.requests.MarkReceived(peerID, key)
			} else if memR.mempool.IsRejectedTx(key) { // we will never request a transaction we've already rejected.
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
			memR.mempool.PeerHasTx(txInfo.SenderID, key)
			if memR.mempool.Has(key) {
				// We have already received this transaction. We mark the peer that sent the message
				// as already seeing the transaction as well and then we finish
				continue
			}
			if !memR.opts.ListenOnly {
				memR.broadcastSeenTx(key)
			}
			_, err = memR.mempool.TryAddNewTx(ntx, key, txInfo)
			if err != nil {
				memR.Logger.Info("Could not add tx", "tx_key", key, "err", err)
			}
		}

	// A peer has indicated to us that it has a transaction. We first verify the txkey and
	// mark that peer as having the transaction. Then we proceed with the following logic:
	//
	// 1. If we have the transaction, we do nothing.
	// 2. If we don't have the transaction, we check if the original sender is a peer we are
	// connected to. The original recipients of a transaction will immediately broadcast it
	// to everyone so if we haven't received it yet we will likely receive it soon. Therefore,
	// we set a timer to check after a certain amount of time if we still don't have the transaction.
	// 3. If we're not connected to the original sender, or we exceed the timeout without
	// receiving the transaction we request it from this peer.
	case *protomem.SeenTx:
		txKey, err := types.TxKeyFromBytes(msg.TxKey)
		if err != nil {
			memR.Logger.Error("Peer sent SeenTx with incorrect key size", "err", err)
			memR.Switch.StopPeerForError(e.Src, err)
			return
		}
		memR.mempool.PeerHasTx(memR.ids.GetIDForPeer(e.Src.ID()), txKey)
		// Check if we don't already have the transaction and that it was recently rejected
		if !memR.mempool.Has(txKey) && !memR.mempool.IsRejectedTx(txKey) {
			// If we are already requesting that tx, then we don't need to go any further.
			if memR.requests.ForTx(txKey) {
				return
			}

			// Check if `From` is specified and we are connected to the peer that originally sent the transaction
			from := ""
			if msg.XFrom.(*protomem.SeenTx_From) != nil {
				from = msg.XFrom.(*protomem.SeenTx_From).From
			}
			if from != "" && memR.ids.GetIDForPeer(p2p.ID(from)) == 0 {
				// We are not connected to the peer that originally sent the transaction
				// so we request it from the peer that sent us the SeenTx message
				memR.requestTx(txKey, e.Src)
			} else {
				// We are connected to the peer that originally sent the transaction so we
				// assume there's a high probability that the original sender will also
				// send us the transaction. We set a timeout in case this is not true.
				time.AfterFunc(memR.opts.MaxGossipDelay, func() {
					// If we still don't have the transaction after the timeout, we find a new peer to request the tx
					if !memR.mempool.Has(txKey) {
						// NOTE: During this period, the peer may, for some reason have disconnected from us.
						if memR.ids.GetIDForPeer(e.Src.ID()) == 0 {
							// Get the first peer from the set
							memR.findNewPeerToSendTx(txKey)
						} else {
							// We're still connected to the peer that sent us the SeenTx so request
							// the transaction from them
							memR.requestTx(txKey, e.Src)
						}

					}
				})
			}
		}

	// A peer is requesting a transaction that we have claimed to have. Find the specified
	// transaction and broadcast it to the peer. We may no longer have the transaction
	case *protomem.WantTx:
		txKey, err := types.TxKeyFromBytes(msg.TxKey)
		if err != nil {
			memR.Logger.Error("Peer sent WantTx with incorrect key size", "err", err)
			memR.Switch.StopPeerForError(e.Src, err)
			return
		}
		tx, has := memR.mempool.Get(types.TxKey(txKey))
		if has && !memR.opts.ListenOnly {
			p2p.SendEnvelopeShim(e.Src, p2p.Envelope{ //nolint: staticcheck
				ChannelID: mempool.MempoolChannel,
				Message:   &protomem.Txs{Txs: [][]byte{tx}},
			}, memR.Logger)
		}
		memR.mempool.PeerHasTx(memR.ids.GetIDForPeer(e.Src.ID()), types.TxKey(txKey))

	default:
		memR.Logger.Error("unknown message type", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)
		memR.Switch.StopPeerForError(e.Src, fmt.Errorf("mempool cannot handle message of type: %T", e.Message))
		return
	}
}

// Receive wraps ReceiveEnvelope supporting backwards compatibility.
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
	alreadySeenTx := memR.mempool.seenByPeersSet.Get(txKey)
	for _, peer := range memR.Switch.Peers().List() {
		peerID := memR.ids.GetIDForPeer(peer.ID())
		if _, ok := alreadySeenTx[peerID]; ok {
			continue
		}
		p2p.SendEnvelopeShim(peer, p2p.Envelope{ //nolint: staticcheck
			ChannelID: mempool.MempoolChannel,
			Message:   &protomem.SeenTx{TxKey: txKey[:]},
		}, memR.Logger)
	}
}

func (memR *Reactor) broadcastNewTx(tx types.Tx) {
	memR.Logger.Debug("broadcasting new tx to all peers")
	memR.Switch.BroadcastEnvelope(p2p.Envelope{
		ChannelID: mempool.MempoolChannel,
		Message:   &protomem.Txs{Txs: [][]byte{tx}},
	})
}

func (memR *Reactor) requestTx(txKey types.TxKey, peer p2p.Peer) {
	memR.Logger.Debug("requesting tx", "tx_key", txKey, "peer", peer)
	success := p2p.SendEnvelopeShim(peer, p2p.Envelope{ //nolint: staticcheck
		ChannelID: MempoolStateChannel,
		Message:   &protomem.WantTx{TxKey: txKey[:]},
	}, memR.Logger)
	if success {
		memR.requests.Add(txKey, memR.ids.GetIDForPeer(peer.ID()), memR.findNewPeerToSendTx)
	}
}

func (memR *Reactor) findNewPeerToSendTx(txKey types.TxKey) {
	// pop the next peer in the list of remaining peers that have seen the tx
	peerID := memR.mempool.seenByPeersSet.Pop(txKey)
	if peerID == 0 {
		// No other peer has the transaction we are looking for.
		// We give up 🤷‍♂️
		return
	}
	peer := memR.ids.GetPeer(peerID)
	memR.requestTx(txKey, peer)
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
