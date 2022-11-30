package v2

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/creachadair/taskgroup"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/clist"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

var _ mempool.Mempool = (*TxMempool)(nil)

const (
	evictedTxCacheSize = 100
)

// TxMempoolOption sets an optional parameter on the TxMempool.
type TxMempoolOption func(*TxMempool)

// TxMempool implemements the Mempool interface and allows the application to
// set priority values on transactions in the CheckTx response. When selecting
// transactions to include in a block, higher-priority transactions are chosen
// first.  When evicting transactions from the mempool for size constraints,
// lower-priority transactions are evicted sooner.
//
// Within the mempool, transactions are ordered by time of arrival, and are
// gossiped to the rest of the network based on that order (gossip order does
// not take priority into account).
type TxMempool struct {
	// Immutable fields
	logger          log.Logger
	config          *config.MempoolConfig
	proxyAppConn    proxy.AppConnMempool
	metrics         *mempool.Metrics
	rejectedTxCache *LRUTxCache // rejected Txs by CheckTx

	// Atomically-updated fields
	txsBytes int64 // atomic: the total size of all transactions in the mempool, in bytes

	// These fields are not synchronized. They are modified in `Update` which should never
	// be called concurrently.
	notifiedTxsAvailable bool
	txsAvailable         chan struct{} // one value sent per height when mempool is not empty
	preCheck             mempool.PreCheckFunc
	postCheck            mempool.PostCheckFunc
	height               int64 // the latest height passed to Update

	// Concurrent list of valid transactions (passed CheckTx)
	txs *clist.CList

	// Synchronized fields, protected by mtx.
	mtx        *sync.RWMutex
	evictedTxs map[types.TxKey]*evictedTxInfo  // txs that were evicted but passed CheckTx
	txByKey    map[types.TxKey]*clist.CElement // used as a lookup table
	txBySender map[string]*clist.CElement      // for sender != ""
}

// NewTxMempool constructs a new, empty priority mempool at the specified
// initial height and using the given config and options.
func NewTxMempool(
	logger log.Logger,
	cfg *config.MempoolConfig,
	proxyAppConn proxy.AppConnMempool,
	height int64,
	options ...TxMempoolOption,
) *TxMempool {

	txmp := &TxMempool{
		logger:          logger,
		config:          cfg,
		proxyAppConn:    proxyAppConn,
		metrics:         mempool.NopMetrics(),
		rejectedTxCache: NewLRUTxCache(cfg.CacheSize),
		evictedTxs:      make(map[types.TxKey]*evictedTxInfo),
		txs:             clist.New(),
		mtx:             new(sync.RWMutex),
		height:          height,
		txByKey:         make(map[types.TxKey]*clist.CElement),
		txBySender:      make(map[string]*clist.CElement),
		preCheck:        func(_ types.Tx) error { return nil },
		postCheck:       func(_ types.Tx, _ *abci.ResponseCheckTx) error { return nil },
	}

	for _, opt := range options {
		opt(txmp)
	}

	return txmp
}

// WithPreCheck sets a filter for the mempool to reject a transaction if f(tx)
// returns an error. This is executed before CheckTx. It only applies to the
// first created block. After that, Update() overwrites the existing value.
func WithPreCheck(f mempool.PreCheckFunc) TxMempoolOption {
	return func(txmp *TxMempool) { txmp.preCheck = f }
}

// WithPostCheck sets a filter for the mempool to reject a transaction if
// f(tx, resp) returns an error. This is executed after CheckTx. It only applies
// to the first created block. After that, Update overwrites the existing value.
func WithPostCheck(f mempool.PostCheckFunc) TxMempoolOption {
	return func(txmp *TxMempool) { txmp.postCheck = f }
}

// WithMetrics sets the mempool's metrics collector.
func WithMetrics(metrics *mempool.Metrics) TxMempoolOption {
	return func(txmp *TxMempool) { txmp.metrics = metrics }
}

// Lock obtains a write-lock on the mempool. A caller must be sure to explicitly
// release the lock when finished. No transactions will be added or removed
// until the lock is released
func (txmp *TxMempool) Lock() { txmp.mtx.Lock() }

// Unlock releases a write-lock on the mempool.
func (txmp *TxMempool) Unlock() { txmp.mtx.Unlock() }

// Size returns the number of valid transactions in the mempool. It is
// thread-safe.
func (txmp *TxMempool) Size() int { return txmp.txs.Len() }

// SizeBytes return the total sum in bytes of all the valid transactions in the
// mempool. It is thread-safe.
func (txmp *TxMempool) SizeBytes() int64 { return atomic.LoadInt64(&txmp.txsBytes) }

// FlushAppConn executes FlushSync on the mempool's proxyAppConn.
//
// The caller must hold an exclusive mempool lock (by calling txmp.Lock) before
// calling FlushAppConn.
func (txmp *TxMempool) FlushAppConn() error {
	// N.B.: We have to issue the call outside the lock so that its callback can
	// fire.  It's safe to do this, the flush will block until complete.
	//
	// We could just not require the caller to hold the lock at all, but the
	// semantics of the Mempool interface require the caller to hold it, and we
	// can't change that without disrupting existing use.
	txmp.mtx.Unlock()
	defer txmp.mtx.Lock()

	return txmp.proxyAppConn.FlushSync()
}

// EnableTxsAvailable enables the mempool to trigger events when transactions
// are available on a block by block basis.
func (txmp *TxMempool) EnableTxsAvailable() {
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()

	txmp.txsAvailable = make(chan struct{}, 1)
}

// TxsAvailable returns a channel which fires once for every height, and only
// when transactions are available in the mempool. It is thread-safe.
func (txmp *TxMempool) TxsAvailable() <-chan struct{} { return txmp.txsAvailable }

func (txmp *TxMempool) SeenTx(txKey types.TxKey) bool {
	txmp.mtx.RLock()
	defer txmp.mtx.RUnlock()
	return txmp.seenTx(txKey)
}

func (txmp *TxMempool) seenTx(txKey types.TxKey) bool {
	_, exists := txmp.txByKey[txKey]
	return exists
}

func (txmp *TxMempool) IsRejectedTx(txKey types.TxKey) bool {
	return txmp.rejectedTxCache.Has(txKey)
}

func (txmp *TxMempool) WasRecentlyEvicted(txKey types.TxKey, tx types.Tx) bool {
	txmp.mtx.RLock()
	defer txmp.mtx.RUnlock()
	_, exists := txmp.evictedTxs[txKey]
	return exists
}

func (txmp *TxMempool) TryReinsertEvictedTx(txKey types.TxKey, tx types.Tx, peer uint16) error {
	txmp.mtx.RLock()
	info, exists := txmp.evictedTxs[txKey]
	txmp.mtx.RUnlock()
	if !exists {
		return nil
	}

	wtx := &wrappedTx{
		tx:        tx,
		key:       tx.Key(),
		timestamp: time.Now().UTC(),
		height:    txmp.height,
		peers:     info.peers,
		priority:  info.priority,
		gasWanted: info.gasWanted,
		sender:    info.sender,
	}
	checkTxResp := &abci.ResponseCheckTx{
		Code:      abci.CodeTypeOK,
		Priority:  info.priority,
		Sender:    info.sender,
		GasWanted: info.gasWanted,
	}
	wtx.SetPeer(peer)
	return txmp.addNewTransaction(wtx, checkTxResp)
}

// CheckTx adds the given transaction to the mempool if it fits and passes the
// application's ABCI CheckTx method.
//
// CheckTx reports an error without adding tx if:
//
// - The size of tx exceeds the configured maximum transaction size.
// - The pre-check hook reports an error for tx.
// - The transaction already exists in the transaction clist or in the rejectedTxCache.
//
// If tx passes all of the above conditions, `TryAddNewTx` is called
func (txmp *TxMempool) CheckTx(tx types.Tx, cb func(*abci.Response), txInfo mempool.TxInfo) error {
	// Reject transactions in excess of the configured maximum transaction size.
	if len(tx) > txmp.config.MaxTxBytes {
		return mempool.ErrTxTooLarge{Max: txmp.config.MaxTxBytes, Actual: len(tx)}
	}

	// If a precheck hook is defined, call it before invoking the application.
	if err := txmp.preCheck(tx); err != nil {
		return mempool.ErrPreCheck{Reason: err}
	}

	key := tx.Key()

	if txmp.IsRejectedTx(key) {
		// The peer has sent us a transaction that we have marked as invalid. Since `CheckTx` can
		// be non-deterministic, we don't punish the peer but instead just ignore the msg
		return mempool.ErrTxInCache
	}

	if txmp.WasRecentlyEvicted(key, tx) {
		// the transaction was recently evicted. If true, we attempt to re-add it to the mempool
		// skipping check tx.
		return nil
	}

	// This is a new transaction that we haven't seen before. Verify it against the app and attempt
	// to add it to the transaction pool.
	rsp, err := txmp.TryAddNewTx(tx, tx.Key(), txInfo)
	if err != nil {
		return err
	}

	// call the callback if it is set
	if cb != nil {
		cb(&abci.Response{Value: &abci.Response_CheckTx{CheckTx: rsp}})
	}
	return nil
}

// TryAddNewTx attempts to add a tx that has not already been seen before. It first marks it as seen
// to avoid races with the same tx. It then call `CheckTx` so that the application can validate it.
// If it passes `CheckTx`, the new transaction is added to the mempool solong as it has
// sufficient priority and space else if evicted it will return an error
func (txmp *TxMempool) TryAddNewTx(tx types.Tx, key types.TxKey, txInfo mempool.TxInfo) (*abci.ResponseCheckTx, error) {
	// reserve the key
	if !txmp.reserveTx(key) {
		return nil, errors.New("tx already added")
	}

	// Early exit if the proxy connection has an error.
	if err := txmp.proxyAppConn.Error(); err != nil {
		txmp.unreserveTx(key)
		return nil, err
	}

	// Invoke an ABCI CheckTx for this transaction.
	rsp, err := txmp.proxyAppConn.CheckTxSync(abci.RequestCheckTx{Tx: tx})
	if err != nil {
		txmp.unreserveTx(key)
		return nil, err
	}
	wtx := &wrappedTx{
		tx:        tx,
		key:       tx.Key(),
		timestamp: time.Now().UTC(),
		height:    txmp.height,
		priority:  rsp.Priority,
		gasWanted: rsp.GasWanted,
		sender:    rsp.Sender,
		peers:     map[uint16]bool{txInfo.SenderID: true},
	}
	if err := txmp.addNewTransaction(wtx, rsp); err != nil {
		txmp.unreserveTx(key)
		return nil, err
	}
	return rsp, nil
}

func (txmp *TxMempool) reserveTx(key types.TxKey) bool {
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()
	if _, ok := txmp.txByKey[key]; ok {
		return false // already reserved
	} else {
		txmp.txByKey[key] = &clist.CElement{}
	}
	return true
}

func (txmp *TxMempool) unreserveTx(key types.TxKey) {
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()
	value, ok := txmp.txByKey[key]
	if ok && value.Value == nil {
		delete(txmp.txByKey, key)
	}
}

// RemoveTxByKey removes the transaction with the specified key from the
// mempool. It reports an error if no such transaction exists. This operation
// does not remove the transaction from the rejectedTxCache.
func (txmp *TxMempool) RemoveTxByKey(txKey types.TxKey) error {
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()
	return txmp.removeTxByKey(txKey)
}

// removeTxByKey removes the specified transaction key from the mempool.
// The caller must hold txmp.mtx excluxively.
func (txmp *TxMempool) removeTxByKey(key types.TxKey) error {
	if elt, ok := txmp.txByKey[key]; ok {
		w := elt.Value.(*wrappedTx)
		delete(txmp.txByKey, key)
		delete(txmp.txBySender, w.sender)
		txmp.txs.Remove(elt)
		elt.DetachPrev()
		elt.DetachNext()
		atomic.AddInt64(&txmp.txsBytes, -w.Size())
		return nil
	}
	return fmt.Errorf("transaction %x not found", key)
}

// removeTxByElement removes the specified transaction element from the mempool.
// The caller must hold txmp.mtx exclusively.
func (txmp *TxMempool) removeTxByElement(elt *clist.CElement) {
	w := elt.Value.(*wrappedTx)
	delete(txmp.txByKey, w.tx.Key())
	delete(txmp.txBySender, w.sender)
	txmp.txs.Remove(elt)
	elt.DetachPrev()
	elt.DetachNext()
	atomic.AddInt64(&txmp.txsBytes, -w.Size())
}

// Flush purges the contents of the mempool and the cache, leaving both empty.
// The current height is not modified by this operation.
func (txmp *TxMempool) Flush() {
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()

	// Remove all the transactions in the list explicitly, so that the sizes
	// and indexes get updated properly.
	cur := txmp.txs.Front()
	for cur != nil {
		next := cur.Next()
		txmp.removeTxByElement(cur)
		cur = next
	}
	txmp.rejectedTxCache.Reset()
}

// PeerHasTx marks that the transaction has been seen by a peer
func (txmp *TxMempool) PeerHasTx(peer uint16, txKey types.TxKey) {
	txmp.mtx.RLock()
	defer txmp.mtx.RUnlock()
	if el, exists := txmp.txByKey[txKey]; exists {
		wtx := el.Value.(*wrappedTx)
		wtx.SetPeer(peer)
	}
}

// allEntriesSorted returns a slice of all the transactions currently in the
// mempool, sorted in nonincreasing order by priority with ties broken by
// increasing order of arrival time.
func (txmp *TxMempool) allEntriesSorted() []*wrappedTx {
	txmp.mtx.RLock()
	defer txmp.mtx.RUnlock()

	all := make([]*wrappedTx, len(txmp.txByKey))
	idx := 0
	for _, tx := range txmp.txByKey {
		all[idx] = tx.Value.(*wrappedTx)
		idx++
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].priority == all[j].priority {
			return all[i].timestamp.Before(all[j].timestamp)
		}
		return all[i].priority > all[j].priority // N.B. higher priorities first
	})
	return all
}

// ReapMaxBytesMaxGas returns a slice of valid transactions that fit within the
// size and gas constraints. The results are ordered by nonincreasing priority,
// with ties broken by increasing order of arrival.  Reaping transactions does
// not remove them from the mempool.add
//
// If maxBytes < 0, no limit is set on the total size in bytes.
// If maxGas < 0, no limit is set on the total gas cost.
//
// If the mempool is empty or has no transactions fitting within the given
// constraints, the result will also be empty.
func (txmp *TxMempool) ReapMaxBytesMaxGas(maxBytes, maxGas int64) types.Txs {
	var totalGas, totalBytes int64

	var keep []types.Tx //nolint:prealloc
	for _, w := range txmp.allEntriesSorted() {
		// N.B. When computing byte size, we need to include the overhead for
		// encoding as protobuf to send to the application.
		totalGas += w.gasWanted
		totalBytes += types.ComputeProtoSizeForTxs([]types.Tx{w.tx})
		if (maxGas >= 0 && totalGas > maxGas) || (maxBytes >= 0 && totalBytes > maxBytes) {
			break
		}
		keep = append(keep, w.tx)
	}
	return keep
}

// TxsWaitChan returns a channel that is closed when there is at least one
// transaction available to be gossiped.
func (txmp *TxMempool) TxsWaitChan() <-chan struct{} { return txmp.txs.WaitChan() }

// TxsFront returns the frontmost element of the pending transaction list.
// It will be nil if the mempool is empty.
func (txmp *TxMempool) TxsFront() *clist.CElement { return txmp.txs.Front() }

// ReapMaxTxs returns up to max transactions from the mempool. The results are
// ordered by nonincreasing priority with ties broken by increasing order of
// arrival. Reaping transactions does not remove them from the mempool.
//
// If max < 0, all transactions in the mempool are reaped.
//
// The result may have fewer than max elements (possibly zero) if the mempool
// does not have that many transactions available.
func (txmp *TxMempool) ReapMaxTxs(max int) types.Txs {
	var keep []types.Tx //nolint:prealloc

	for _, w := range txmp.allEntriesSorted() {
		if max >= 0 && len(keep) >= max {
			break
		}
		keep = append(keep, w.tx)
	}
	return keep
}

// Update removes all the given transactions from the mempool and the cache,
// and updates the current block height. The blockTxs and deliverTxResponses
// must have the same length with each response corresponding to the tx at the
// same offset.
//
// If the configuration enables recheck, Update sends each remaining
// transaction after removing blockTxs to the ABCI CheckTx method.  Any
// transactions marked as invalid during recheck are also removed.
//
// The caller must hold an exclusive mempool lock (by calling txmp.Lock) before
// calling Update.
func (txmp *TxMempool) Update(
	blockHeight int64,
	blockTxs types.Txs,
	deliverTxResponses []*abci.ResponseDeliverTx,
	newPreFn mempool.PreCheckFunc,
	newPostFn mempool.PostCheckFunc,
) error {
	// Safety check: Transactions and responses must match in number.
	if len(blockTxs) != len(deliverTxResponses) {
		panic(fmt.Sprintf("mempool: got %d transactions but %d DeliverTx responses",
			len(blockTxs), len(deliverTxResponses)))
	}

	txmp.height = blockHeight
	txmp.notifiedTxsAvailable = false

	if newPreFn != nil {
		txmp.preCheck = newPreFn
	}
	if newPostFn != nil {
		txmp.postCheck = newPostFn
	}

	for _, tx := range blockTxs {
		// Regardless of success, remove the transaction from the mempool.
		_ = txmp.removeTxByKey(tx.Key())
	}

	txmp.purgeExpiredTxs(blockHeight)

	// If there any uncommitted transactions left in the mempool, we either
	// initiate re-CheckTx per remaining transaction or notify that remaining
	// transactions are left.
	size := txmp.Size()
	txmp.metrics.Size.Set(float64(size))
	if size > 0 {
		if txmp.config.Recheck {
			txmp.recheckTransactions()
		} else {
			txmp.notifyTxsAvailable()
		}
	}
	return nil
}

// addNewTransaction handles the ABCI CheckTx response for the first time a
// transaction is added to the mempool.  A recheck after a block is committed
// goes to handleRecheckResult.
//
// If either the application rejected the transaction or a post-check hook is
// defined and rejects the transaction, it is discarded.
//
// Otherwise, if the mempool is full, check for lower-priority transactions
// that can be evicted to make room for the new one. If no such transactions
// exist, this transaction is logged and dropped; otherwise the selected
// transactions are evicted.
//
// Finally, the new transaction is added and size stats updated.
func (txmp *TxMempool) addNewTransaction(wtx *wrappedTx, checkTxRes *abci.ResponseCheckTx) error {
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()

	var err error
	if txmp.postCheck != nil {
		err = txmp.postCheck(wtx.tx, checkTxRes)
	}

	if err != nil || checkTxRes.Code != abci.CodeTypeOK {
		txmp.logger.Info(
			"rejected bad transaction",
			"priority", wtx.priority,
			"tx", fmt.Sprintf("%X", wtx.tx.Hash()),
			"peer_id", wtx.peers,
			"code", checkTxRes.Code,
			"post_check_err", err,
		)

		txmp.metrics.RejectedTxs.Add(1)

		// If there was a post-check error, record its text in the result for
		// debugging purposes.
		if err != nil {
			checkTxRes.MempoolError = err.Error()
		}
		return err
	}

	// At this point the application has ruled the transaction valid, but the
	// mempool might be full. If so, find the lowest-priority items with lower
	// priority than the application assigned to this new one, and evict as many
	// of them as necessary to make room for tx. If no such items exist, we
	// discard tx.

	if err := txmp.canAddTx(wtx); err != nil {
		var victims []*clist.CElement // eligible transactions for eviction
		var victimBytes int64         // total size of victims
		for cur := txmp.txs.Front(); cur != nil; cur = cur.Next() {
			cw := cur.Value.(*wrappedTx)
			if cw.priority < wtx.priority {
				victims = append(victims, cur)
				victimBytes += cw.Size()
			}
		}

		// If there are no suitable eviction candidates, or the total size of
		// those candidates is not enough to make room for the new transaction,
		// drop the new one.
		if len(victims) == 0 || victimBytes < wtx.Size() {
			txmp.metrics.EvictedTxs.Add(1)
			txmp.cacheEvictedTx(wtx)
			checkTxRes.MempoolError =
				fmt.Sprintf("rejected valid incoming transaction; mempool is full (%X)",
					wtx.tx.Hash())
			return fmt.Errorf("rejected valid incoming transaction; mempool is full (%X)",
				wtx.tx.Hash())
		}

		txmp.logger.Debug("evicting lower-priority transactions",
			"new_tx", fmt.Sprintf("%X", wtx.tx.Hash()),
			"new_priority", wtx.priority,
		)

		// Sort lowest priority items first so they will be evicted first.  Break
		// ties in favor of newer items (to maintain FIFO semantics in a group).
		sort.Slice(victims, func(i, j int) bool {
			iw := victims[i].Value.(*wrappedTx)
			jw := victims[j].Value.(*wrappedTx)
			if iw.priority == jw.priority {
				return iw.timestamp.After(jw.timestamp)
			}
			return iw.priority < jw.priority
		})

		// Evict as many of the victims as necessary to make room.
		var evictedBytes int64
		for _, vic := range victims {
			w := vic.Value.(*wrappedTx)
			txmp.evictTx(w)

			// We may not need to evict all the eligible transactions.  Bail out
			// early if we have made enough room.
			evictedBytes += w.Size()
			if evictedBytes >= wtx.Size() {
				break
			}
		}
	}

	txmp.insertTx(wtx)

	txmp.metrics.TxSizeBytes.Observe(float64(wtx.Size()))
	txmp.metrics.Size.Set(float64(txmp.Size()))
	txmp.logger.Debug(
		"inserted new valid transaction",
		"priority", wtx.priority,
		"tx", fmt.Sprintf("%X", wtx.tx.Hash()),
		"height", txmp.height,
		"num_txs", txmp.Size(),
	)
	txmp.notifyTxsAvailable()
	return nil
}

func (txmp *TxMempool) insertTx(wtx *wrappedTx) {
	elt := txmp.txs.PushBack(wtx)
	txmp.txByKey[wtx.tx.Key()] = elt
	if s := wtx.sender; s != "" {
		txmp.txBySender[s] = elt
	}
	// if we're reinserting an evicted transaction
	// remove it from the map
	delete(txmp.evictedTxs, wtx.key)

	atomic.AddInt64(&txmp.txsBytes, wtx.Size())
}

func (txmp *TxMempool) evictTx(wtx *wrappedTx) {
	txmp.removeTxByKey(wtx.key)
	txmp.metrics.EvictedTxs.Add(1)
	txmp.cacheEvictedTx(wtx)
	txmp.logger.Debug(
		"evicted valid existing transaction; mempool full",
		"old_tx", fmt.Sprintf("%X", wtx.key),
		"old_priority", wtx.priority,
	)
}

func (txmp *TxMempool) cacheEvictedTx(wtx *wrappedTx) {
	txmp.evictedTxs[wtx.key] = &evictedTxInfo{
		timeEvicted: time.Now().UTC(),
		priority:    wtx.priority,
		gasWanted:   wtx.gasWanted,
		sender:      wtx.sender,
		peers:       wtx.peers,
	}
	// if cache too large, remove the oldest entry
	if len(txmp.evictedTxs) > evictedTxCacheSize {
		oldestTxKey := wtx.key
		oldestTxTime := time.Now().UTC()
		for key, info := range txmp.evictedTxs {
			if info.timeEvicted.Before(oldestTxTime) {
				oldestTxTime = info.timeEvicted
				oldestTxKey = key
			}
		}
		delete(txmp.evictedTxs, oldestTxKey)
	}
}

// handleRecheckResult handles the responses from ABCI CheckTx calls issued
// during the recheck phase of a block Update.  It removes any transactions
// invalidated by the application.
//
// This method is NOT executed for the initial CheckTx on a new transaction;
// that case is handled by addNewTransaction instead.
func (txmp *TxMempool) handleRecheckResult(tx types.Tx, checkTxRes *abci.ResponseCheckTx) {
	txmp.metrics.RecheckTimes.Add(1)
	txmp.mtx.Lock()
	defer txmp.mtx.Unlock()

	// Find the transaction reported by the ABCI callback. It is possible the
	// transaction was evicted during the recheck, in which case the transaction
	// will be gone.
	elt, ok := txmp.txByKey[tx.Key()]
	if !ok {
		return
	}
	wtx := elt.Value.(*wrappedTx)

	// If a postcheck hook is defined, call it before checking the result.
	var err error
	if txmp.postCheck != nil {
		err = txmp.postCheck(tx, checkTxRes)
	}

	if checkTxRes.Code == abci.CodeTypeOK && err == nil {
		// Note that we do not update the transaction with any of the values returned in
		// recheck tx
		return // N.B. Size of mempool did not change
	}

	txmp.logger.Debug(
		"existing transaction no longer valid; failed re-CheckTx callback",
		"priority", wtx.priority,
		"tx", fmt.Sprintf("%X", wtx.key),
		"err", err,
		"code", checkTxRes.Code,
	)
	txmp.removeTxByElement(elt)
	txmp.metrics.FailedTxs.Add(1)
	txmp.metrics.Size.Set(float64(txmp.Size()))
}

// recheckTransactions initiates re-CheckTx ABCI calls for all the transactions
// currently in the mempool. It reports the number of recheck calls that were
// successfully initiated.
//
// Precondition: The mempool is not empty.
// The caller must hold txmp.mtx exclusively.
func (txmp *TxMempool) recheckTransactions() {
	if txmp.Size() == 0 {
		panic("mempool: cannot run recheck on an empty mempool")
	}
	txmp.logger.Debug(
		"executing re-CheckTx for all remaining transactions",
		"num_txs", txmp.Size(),
		"height", txmp.height,
	)

	// Collect transactions currently in the mempool requiring recheck.
	wtxs := make([]*wrappedTx, 0, txmp.txs.Len())
	for e := txmp.txs.Front(); e != nil; e = e.Next() {
		wtxs = append(wtxs, e.Value.(*wrappedTx))
	}

	// Issue CheckTx calls for each remaining transaction, and when all the
	// rechecks are complete signal watchers that transactions may be available.
	go func() {
		g, start := taskgroup.New(nil).Limit(2 * runtime.NumCPU())

		for _, wtx := range wtxs {
			wtx := wtx
			start(func() error {
				// The response for this CheckTx is handled by the default recheckTxCallback.
				rsp, err := txmp.proxyAppConn.CheckTxSync(abci.RequestCheckTx{
					Tx:   wtx.tx,
					Type: abci.CheckTxType_Recheck,
				})
				if err != nil {
					txmp.logger.Error("failed to execute CheckTx during recheck",
						"err", err, "hash", fmt.Sprintf("%x", wtx.tx.Hash()))
				} else {
					txmp.handleRecheckResult(wtx.tx, rsp)
				}
				return nil
			})
		}
		_ = txmp.proxyAppConn.FlushAsync()

		// When recheck is complete, trigger a notification for more transactions.
		_ = g.Wait()
		txmp.mtx.Lock()
		defer txmp.mtx.Unlock()
		txmp.notifyTxsAvailable()
	}()
}

// canAddTx returns an error if we cannot insert the provided *wrappedTx into
// the mempool due to mempool configured constraints. Otherwise, nil is
// returned and the transaction can be inserted into the mempool.
func (txmp *TxMempool) canAddTx(wtx *wrappedTx) error {
	numTxs := txmp.Size()
	txBytes := txmp.SizeBytes()

	if numTxs >= txmp.config.Size || wtx.Size()+txBytes > txmp.config.MaxTxsBytes {
		return mempool.ErrMempoolIsFull{
			NumTxs:      numTxs,
			MaxTxs:      txmp.config.Size,
			TxsBytes:    txBytes,
			MaxTxsBytes: txmp.config.MaxTxsBytes,
		}
	}

	return nil
}

// purgeExpiredTxs removes all transactions from the mempool that have exceeded
// their respective height or time-based limits as of the given blockHeight.
// Transactions removed by this operation are not removed from the rejectedTxCache.
//
// The caller must hold txmp.mtx exclusively.
func (txmp *TxMempool) purgeExpiredTxs(blockHeight int64) {
	if txmp.config.TTLNumBlocks == 0 && txmp.config.TTLDuration == 0 {
		return // nothing to do
	}

	now := time.Now()
	cur := txmp.txs.Front()
	for cur != nil {
		// N.B. Grab the next element first, since if we remove cur its successor
		// will be invalidated.
		next := cur.Next()

		w := cur.Value.(*wrappedTx)
		if txmp.config.TTLNumBlocks > 0 && (blockHeight-w.height) > txmp.config.TTLNumBlocks {
			txmp.removeTxByElement(cur)
			txmp.metrics.EvictedTxs.Add(1)
		} else if txmp.config.TTLDuration > 0 && now.Sub(w.timestamp) > txmp.config.TTLDuration {
			txmp.removeTxByElement(cur)
			txmp.metrics.EvictedTxs.Add(1)
		}
		cur = next
	}

	// purge old evicted transactions
	if txmp.config.TTLDuration > 0 {
		for key, info := range txmp.evictedTxs {
			if now.Sub(info.timeEvicted) > txmp.config.TTLDuration {
				delete(txmp.evictedTxs, key)
			}
		}
	}
}

func (txmp *TxMempool) notifyTxsAvailable() {
	if txmp.Size() == 0 {
		return // nothing to do
	}

	if txmp.txsAvailable != nil && !txmp.notifiedTxsAvailable {
		// channel cap is 1, so this will send once
		txmp.notifiedTxsAvailable = true

		select {
		case txmp.txsAvailable <- struct{}{}:
		default:
		}
	}
}
