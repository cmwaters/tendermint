package cat

import (
	"sync"
	"time"

	"github.com/tendermint/tendermint/types"
)

// requestScheduler tracks the lifecycle of outbound transaction requests.
type requestScheduler struct {
	mtx sync.Mutex

	// responseTime is the time the scheduler
	// waits for a response from a peer before
	// invoking the callback
	responseTime time.Duration

	// requestsByPeer is a lookup table of requests by peer.
	// Multiple tranasctions can be requested by a single peer at one
	requestsByPeer map[uint16]RequestSet

	// requestsByTx is a lookup table for requested txs.
	// There can only be one request per tx.
	requestsByTx map[types.TxKey]uint16
}

func newRequestScheduler(responseTime time.Duration) *requestScheduler {
	return &requestScheduler{
		responseTime:   responseTime,
		requestsByPeer: make(map[uint16]RequestSet),
		requestsByTx:   make(map[types.TxKey]uint16),
	}
}

func (r *requestScheduler) Add(key types.TxKey, peer uint16, onTimeout func(key types.TxKey)) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	timer := time.AfterFunc(r.responseTime, func() {
		// cleanup resources
		r.MarkReceived(peer, key)
		// trigger callback
		onTimeout(key)
	})
	if _, ok := r.requestsByPeer[peer]; !ok {
		r.requestsByPeer[peer] = RequestSet{key: timer}
	} else {
		r.requestsByPeer[peer][key] = timer
	}
	r.requestsByTx[key] = peer
}

func (r *requestScheduler) ForTx(key types.TxKey) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	_, ok := r.requestsByTx[key]
	return ok
}

func (r *requestScheduler) From(peer uint16) RequestSet {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	requestSet, ok := r.requestsByPeer[peer]
	if !ok {
		return RequestSet{}
	}
	return requestSet
}

func (r *requestScheduler) MarkReceived(peer uint16, key types.TxKey) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if timer, ok := r.requestsByPeer[peer][key]; ok {
		timer.Stop()
	}

	delete(r.requestsByPeer[peer], key)
	delete(r.requestsByTx, key)
}

type RequestSet map[types.TxKey]*time.Timer

func (rs RequestSet) Includes(key types.TxKey) bool {
	_, ok := rs[key]
	return ok
}
