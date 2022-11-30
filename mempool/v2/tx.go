package v2

import (
	"time"

	"github.com/tendermint/tendermint/types"
)

// wrappedTx defines a wrapper around a raw transaction with additional metadata
// that is used for indexing. With the exception of the map of peers who have
// seen this transaction, this struct should never be modified
type wrappedTx struct {
	tx        types.Tx        // the original transaction data
	key       types.TxKey     // the transaction hash
	height    int64           // height when this transaction was initially checked (for expiry)
	timestamp time.Time       // time when transaction was entered (for TTL)
	gasWanted int64           // app: gas required to execute this transaction
	priority  int64           // app: priority value for this transaction
	sender    string          // app: assigned sender label
	peers     map[uint16]bool // peer IDs who have sent us this transaction
}

// Size reports the size of the raw transaction in bytes.
func (w *wrappedTx) Size() int64 { return int64(len(w.tx)) }

// SetPeer adds the specified peer ID as a sender of w.
func (w *wrappedTx) SetPeer(id uint16) {
	if w.peers == nil {
		w.peers = map[uint16]bool{id: true}
	} else {
		w.peers[id] = true
	}
}

// HasPeer reports whether the specified peer ID is a sender of w.
func (w *wrappedTx) HasPeer(id uint16) bool {
	_, ok := w.peers[id]
	return ok
}
