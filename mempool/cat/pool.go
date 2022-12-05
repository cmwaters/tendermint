package cat

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/creachadair/taskgroup"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

var _ mempool.Mempool = (*TxPool)(nil)

const (
	evictedTxCacheSize = 100
	seenByPeerSetSize  = 200
)

// TxPoolOption sets an optional parameter on the TxPool.
type TxPoolOption func(*TxPool)

// TxPool implemements the Mempool interface and allows the application to
// set priority values on transactions in the CheckTx response. When selecting
// transactions to include in a block, higher-priority transactions are chosen
// first.  When evicting transactions from the mempool for size constraints,
// lower-priority transactions are evicted sooner.
//
// Within the txpool, transactions are ordered by time of arrival, and are
// gossiped to the rest of the network based on that order (gossip order does
// not take priority into account).
type TxPool struct {
	// Immutable fields
	logger       log.Logger
	config       *config.MempoolConfig
	proxyAppConn proxy.AppConnMempool
	metrics      *mempool.Metrics

	// These fields are not synchronized. They are modified in `Update` which should never
	// be called concurrently.
	notifiedTxsAvailable bool
	txsAvailable         chan struct{} // one value sent per height when mempool is not empty
	preCheck             mempool.PreCheckFunc
	postCheck            mempool.PostCheckFunc
	height               int64 // the latest height passed to Update

	// new transactions that need to be broadcasted to all peers
	outboundTxs chan types.Tx
	// Thread-safe cache of rejected transactions for quick look-up
	rejectedTxCache *LRUTxCache
	// Thread-safe cache of valid txs that were evicted
	evictedTxs *EvictedTxCache
	// Thread-safe list of transactions peers have seen that we have not yet seen
	seenByPeersSet *SeenTxSet

	// Store of wrapped transactions
	store *store
}

// NewTxPool constructs a new, empty priority mempool at the specified
// initial height and using the given config and options.
func NewTxPool(
	logger log.Logger,
	cfg *config.MempoolConfig,
	proxyAppConn proxy.AppConnMempool,
	height int64,
	options ...TxPoolOption,
) *TxPool {

	txmp := &TxPool{
		logger:          logger,
		config:          cfg,
		proxyAppConn:    proxyAppConn,
		metrics:         mempool.NopMetrics(),
		rejectedTxCache: NewLRUTxCache(cfg.CacheSize),
		evictedTxs:      NewEvictedTxCache(evictedTxCacheSize),
		seenByPeersSet:  NewSeenTxSet(seenByPeerSetSize),
		height:          height,
		preCheck:        func(_ types.Tx) error { return nil },
		postCheck:       func(_ types.Tx, _ *abci.ResponseCheckTx) error { return nil },
		store:           newStore(),
		outboundTxs:     make(chan types.Tx, 1),
	}

	for _, opt := range options {
		opt(txmp)
	}

	return txmp
}

// WithPreCheck sets a filter for the mempool to reject a transaction if f(tx)
// returns an error. This is executed before CheckTx. It only applies to the
// first created block. After that, Update() overwrites the existing value.
func WithPreCheck(f mempool.PreCheckFunc) TxPoolOption {
	return func(txmp *TxPool) { txmp.preCheck = f }
}

// WithPostCheck sets a filter for the mempool to reject a transaction if
// f(tx, resp) returns an error. This is executed after CheckTx. It only applies
// to the first created block. After that, Update overwrites the existing value.
func WithPostCheck(f mempool.PostCheckFunc) TxPoolOption {
	return func(txmp *TxPool) { txmp.postCheck = f }
}

// WithMetrics sets the mempool's metrics collector.
func WithMetrics(metrics *mempool.Metrics) TxPoolOption {
	return func(txmp *TxPool) { txmp.metrics = metrics }
}

// Lock obtains a write-lock on the mempool. A caller must be sure to explicitly
// release the lock when finished. No transactions will be added or removed
// until the lock is released
func (txmp *TxPool) Lock() {}

// Unlock releases a write-lock on the mempool.
func (txmp *TxPool) Unlock() {}

// Size returns the number of valid transactions in the mempool. It is
// thread-safe.
func (txmp *TxPool) Size() int { return txmp.store.size() }

// SizeBytes return the total sum in bytes of all the valid transactions in the
// mempool. It is thread-safe.
func (txmp *TxPool) SizeBytes() int64 { return txmp.store.totalBytes() }

// FlushAppConn executes FlushSync on the mempool's proxyAppConn.
//
// The caller must hold an exclusive mempool lock (by calling txmp.Lock) before
// calling FlushAppConn.
func (txmp *TxPool) FlushAppConn() error {
	return txmp.proxyAppConn.FlushSync()
}

// EnableTxsAvailable enables the mempool to trigger events when transactions
// are available on a block by block basis.
func (txmp *TxPool) EnableTxsAvailable() {
	txmp.txsAvailable = make(chan struct{}, 1)
}

// TxsAvailable returns a channel which fires once for every height, and only
// when transactions are available in the mempool. It is thread-safe.
func (txmp *TxPool) TxsAvailable() <-chan struct{} { return txmp.txsAvailable }

func (txmp *TxPool) Has(txKey types.TxKey) bool {
	return txmp.store.has(txKey)
}

func (txmp *TxPool) Get(txKey types.TxKey) (types.Tx, bool) {
	wtx := txmp.store.get(txKey)
	if wtx != nil {
		return wtx.tx, true
	}
	return types.Tx{}, false
}

func (txmp *TxPool) IsRejectedTx(txKey types.TxKey) bool {
	return txmp.rejectedTxCache.Has(txKey)
}

func (txmp *TxPool) WasRecentlyEvicted(txKey types.TxKey) bool {
	return txmp.evictedTxs.Has(txKey)
}

func (txmp *TxPool) TryReinsertEvictedTx(txKey types.TxKey, tx types.Tx, peer uint16) error {
	info := txmp.evictedTxs.Pop(txKey)
	if info == nil {
		return nil
	}
	txmp.logger.Debug("attempting to reinsert evicted tx", "txKey", fmt.Sprintf("%X", txKey))
	wtx := newWrappedTx(
		tx, txKey, txmp.height, info.gasWanted, info.priority, info.sender,
	)
	checkTxResp := &abci.ResponseCheckTx{
		Code:      abci.CodeTypeOK,
		Priority:  info.priority,
		Sender:    info.sender,
		GasWanted: info.gasWanted,
	}
	return txmp.addNewTransaction(wtx, checkTxResp)
}

// CheckTx implements the Mempool interface and wraps around `VerifyAndAddTxFromClient`
func (txmp *TxPool) CheckTx(tx types.Tx, cb func(*abci.Response), txInfo mempool.TxInfo) error {
	return txmp.VerifyAndAddTxFromClient(tx, cb, txInfo)
}

// VerifyAndAddTxFromClient adds the given transaction to the mempool if it fits and passes the
// application's ABCI CheckTx method. This should be viewed as the entry method for new transactions
// into the network. In practice this happens via an RPC endpoint
//
// VerifyAndAddTxFromClient reports an error without adding tx if:
//
// - The size of tx exceeds the configured maximum transaction size.
// - The pre-check hook reports an error for tx.
// - The transaction already exists in the transaction clist or in the rejectedTxCache.
//
// If tx passes all of the above conditions, `TryAddNewTx` is called.
func (txmp *TxPool) VerifyAndAddTxFromClient(tx types.Tx, cb func(*abci.Response), txInfo mempool.TxInfo) error {
	// Reject transactions in excess of the configured maximum transaction size.
	if len(tx) > txmp.config.MaxTxBytes {
		return mempool.ErrTxTooLarge{Max: txmp.config.MaxTxBytes, Actual: len(tx)}
	}

	key := tx.Key()

	if txmp.IsRejectedTx(key) {
		// The peer has sent us a transaction that we have marked as invalid. Since `CheckTx` can
		// be non-deterministic, we don't punish the peer but instead just ignore the msg
		return mempool.ErrTxInCache
	}

	if txmp.WasRecentlyEvicted(key) {
		// the transaction was recently evicted. If true, we attempt to re-add it to the mempool
		// skipping check tx.
		return txmp.TryReinsertEvictedTx(key, tx, txInfo.SenderID)
	}

	// This is a new transaction that we haven't seen before. Verify it against the app and attempt
	// to add it to the transaction pool.
	rsp, err := txmp.TryAddNewTx(tx, key, txInfo)
	if err != nil {
		return err
	}

	// send this to the reactor to be broadcast to all peers
	txmp.outboundTxs <- tx

	// call the callback if it is set
	if cb != nil {
		cb(&abci.Response{Value: &abci.Response_CheckTx{CheckTx: rsp}})
	}
	return nil
}

func (txmp *TxPool) OutboundTxs() <-chan types.Tx {
	return txmp.outboundTxs
}

// TryAddNewTx attempts to add a tx that has not already been seen before. It first marks it as seen
// to avoid races with the same tx. It then call `CheckTx` so that the application can validate it.
// If it passes `CheckTx`, the new transaction is added to the mempool solong as it has
// sufficient priority and space else if evicted it will return an error
func (txmp *TxPool) TryAddNewTx(tx types.Tx, key types.TxKey, txInfo mempool.TxInfo) (*abci.ResponseCheckTx, error) {
	// reserve the key
	if !txmp.store.reserve(key) {
		txmp.logger.Debug("mempool already attempting to verify and add transaction", "txKey", fmt.Sprintf("%X", key))
		txmp.PeerHasTx(txInfo.SenderID, key)
		return nil, errors.New("tx already added")
	}

	resp, err := txmp.tryAddNewTx(tx, key, txInfo)
	if err != nil {
		// remove the reservation if adding failed
		txmp.store.release(key)
	}
	return resp, err
}

// TryAddNewTx attempts to add a tx that has not already been seen before. It first marks it as seen
// to avoid races with the same tx. It then call `CheckTx` so that the application can validate it.
// If it passes `CheckTx`, the new transaction is added to the mempool solong as it has
// sufficient priority and space else if evicted it will return an error
func (txmp *TxPool) tryAddNewTx(tx types.Tx, key types.TxKey, txInfo mempool.TxInfo) (*abci.ResponseCheckTx, error) {
	// Reject transactions in excess of the configured maximum transaction size.
	if len(tx) > txmp.config.MaxTxBytes {
		return nil, mempool.ErrTxTooLarge{Max: txmp.config.MaxTxBytes, Actual: len(tx)}
	}

	// reserve the key
	if !txmp.store.reserve(key) {
		txmp.logger.Debug("mempool already attempting to verify and add transaction", "txKey", fmt.Sprintf("%X", key))
		txmp.PeerHasTx(txInfo.SenderID, key)
		return nil, errors.New("tx already added")
	}
	defer txmp.store.release(key)

	// If a precheck hook is defined, call it before invoking the application.
	if err := txmp.preCheck(tx); err != nil {
		return nil, mempool.ErrPreCheck{Reason: err}
	}

	// Early exit if the proxy connection has an error.
	if err := txmp.proxyAppConn.Error(); err != nil {
		return nil, err
	}

	// Invoke an ABCI CheckTx for this transaction.
	rsp, err := txmp.proxyAppConn.CheckTxSync(abci.RequestCheckTx{Tx: tx})
	if err != nil {
		return rsp, err
	}
	if rsp.Code != abci.CodeTypeOK {
		txmp.metrics.RejectedTxs.Add(1)
		return rsp, fmt.Errorf("application rejected transaction with code %d", rsp.Code)
	}

	// Create wrapped tx
	wtx := newWrappedTx(
		tx, key, txmp.height, rsp.GasWanted, rsp.Priority, rsp.Sender,
	)

	// Perform the post check
	err = txmp.postCheck(wtx.tx, rsp)
	if err != nil {
		txmp.metrics.RejectedTxs.Add(1)
		return rsp, fmt.Errorf("rejected bad transaction after post check: %w", err)
	}

	// Now we consider the transaction to be valid. Once a transaction is valid, it
	// can only become invalid if recheckTx is enabled and RecheckTx returns a non zero code
	if err := txmp.addNewTransaction(wtx, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

// RemoveTxByKey removes the transaction with the specified key from the
// mempool. It reports an error if no such transaction exists. This operation
// does not remove the transaction from the rejectedTxCache.
func (txmp *TxPool) RemoveTxByKey(txKey types.TxKey) error {
	txmp.store.remove(txKey)
	return nil
}

// Flush purges the contents of the mempool and the cache, leaving both empty.
// The current height is not modified by this operation.
func (txmp *TxPool) Flush() {
	// Remove all the transactions in the list explicitly, so that the sizes
	// and indexes get updated properly.
	txmp.store = newStore()
	txmp.seenByPeersSet.Reset()
	txmp.evictedTxs.Reset()
	txmp.rejectedTxCache.Reset()
}

// PeerHasTx marks that the transaction has been seen by a peer.
// It returns true if the mempool has the transaction and has recorded the
// peer and false if the mempool has not yet seen the transaction that the
// peer has
func (txmp *TxPool) PeerHasTx(peer uint16, txKey types.TxKey) {
	txmp.seenByPeersSet.Add(txKey, peer)
}

// allEntriesSorted returns a slice of all the transactions currently in the
// mempool, sorted in nonincreasing order by priority with ties broken by
// increasing order of arrival time.
func (txmp *TxPool) allEntriesSorted() []*wrappedTx {
	txs := txmp.store.getAllTxs()
	sort.Slice(txs, func(i, j int) bool {
		if txs[i].priority == txs[j].priority {
			return txs[i].timestamp.Before(txs[j].timestamp)
		}
		return txs[i].priority > txs[j].priority // N.B. higher priorities first
	})
	return txs
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
func (txmp *TxPool) ReapMaxBytesMaxGas(maxBytes, maxGas int64) types.Txs {
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

// ReapMaxTxs returns up to max transactions from the mempool. The results are
// ordered by nonincreasing priority with ties broken by increasing order of
// arrival. Reaping transactions does not remove them from the mempool.
//
// If max < 0, all transactions in the mempool are reaped.
//
// The result may have fewer than max elements (possibly zero) if the mempool
// does not have that many transactions available.
func (txmp *TxPool) ReapMaxTxs(max int) types.Txs {
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
func (txmp *TxPool) Update(
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
		key := tx.Key()
		_ = txmp.store.remove(key)
		_ = txmp.evictedTxs.Pop(key)
		_ = txmp.seenByPeersSet.Pop(key)
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
func (txmp *TxPool) addNewTransaction(wtx *wrappedTx, checkTxRes *abci.ResponseCheckTx) error {
	// At this point the application has ruled the transaction valid, but the
	// mempool might be full. If so, find the lowest-priority items with lower
	// priority than the application assigned to this new one, and evict as many
	// of them as necessary to make room for tx. If no such items exist, we
	// discard tx.
	if !txmp.canAddTx(wtx) {
		victims, victimBytes := txmp.store.getTxsBelowPriority(wtx.priority)

		// If there are no suitable eviction candidates, or the total size of
		// those candidates is not enough to make room for the new transaction,
		// drop the new one.
		if len(victims) == 0 || victimBytes < wtx.size() {
			txmp.metrics.EvictedTxs.Add(1)
			txmp.evictedTxs.Push(wtx)
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
			iw := victims[i]
			jw := victims[j]
			if iw.priority == jw.priority {
				return iw.timestamp.After(jw.timestamp)
			}
			return iw.priority < jw.priority
		})

		// Evict as many of the victims as necessary to make room.
		var evictedBytes int64
		for _, tx := range victims {
			txmp.evictTx(tx)

			// We may not need to evict all the eligible transactions.  Bail out
			// early if we have made enough room.
			evictedBytes += tx.size()
			if evictedBytes >= wtx.size() {
				break
			}
		}
	}

	txmp.store.set(wtx)

	txmp.metrics.TxSizeBytes.Observe(float64(wtx.size()))
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

func (txmp *TxPool) evictTx(wtx *wrappedTx) {
	txmp.store.remove(wtx.key)
	txmp.metrics.EvictedTxs.Add(1)
	txmp.evictedTxs.Push(wtx)
	txmp.logger.Debug(
		"evicted valid existing transaction; mempool full",
		"old_tx", fmt.Sprintf("%X", wtx.key),
		"old_priority", wtx.priority,
	)
}

// handleRecheckResult handles the responses from ABCI CheckTx calls issued
// during the recheck phase of a block Update.  It removes any transactions
// invalidated by the application.
//
// This method is NOT executed for the initial CheckTx on a new transaction;
// that case is handled by addNewTransaction instead.
func (txmp *TxPool) handleRecheckResult(wtx *wrappedTx, checkTxRes *abci.ResponseCheckTx) {
	txmp.metrics.RecheckTimes.Add(1)

	// If a postcheck hook is defined, call it before checking the result.
	var err error
	if txmp.postCheck != nil {
		err = txmp.postCheck(wtx.tx, checkTxRes)
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
	txmp.store.remove(wtx.key)
	txmp.metrics.FailedTxs.Add(1)
	txmp.metrics.Size.Set(float64(txmp.Size()))
}

// recheckTransactions initiates re-CheckTx ABCI calls for all the transactions
// currently in the mempool. It reports the number of recheck calls that were
// successfully initiated.
//
// Precondition: The mempool is not empty.
// The caller must hold txmp.mtx exclusively.
func (txmp *TxPool) recheckTransactions() {
	if txmp.Size() == 0 {
		panic("mempool: cannot run recheck on an empty mempool")
	}
	txmp.logger.Debug(
		"executing re-CheckTx for all remaining transactions",
		"num_txs", txmp.Size(),
		"height", txmp.height,
	)

	// Collect transactions currently in the mempool requiring recheck.
	wtxs := txmp.store.getAllTxs()

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
					txmp.handleRecheckResult(wtx, rsp)
				}
				return nil
			})
		}
		_ = txmp.proxyAppConn.FlushAsync()

		// When recheck is complete, trigger a notification for more transactions.
		_ = g.Wait()
		txmp.notifyTxsAvailable()
	}()
}

// canAddTx returns an error if we cannot insert the provided *wrappedTx into
// the mempool due to mempool configured constraints. Otherwise, nil is
// returned and the transaction can be inserted into the mempool.
func (txmp *TxPool) canAddTx(wtx *wrappedTx) bool {
	numTxs := txmp.Size()
	txBytes := txmp.SizeBytes()

	if numTxs >= txmp.config.Size || wtx.size()+txBytes > txmp.config.MaxTxsBytes {
		return false
	}

	return true
}

// purgeExpiredTxs removes all transactions from the mempool that have exceeded
// their respective height or time-based limits as of the given blockHeight.
// Transactions removed by this operation are not removed from the rejectedTxCache.
//
// The caller must hold txmp.mtx exclusively.
func (txmp *TxPool) purgeExpiredTxs(blockHeight int64) {
	if txmp.config.TTLNumBlocks == 0 && txmp.config.TTLDuration == 0 {
		return // nothing to do
	}

	expirationHeight := blockHeight - txmp.config.TTLNumBlocks
	if txmp.config.TTLNumBlocks == 0 {
		expirationHeight = 0
	}

	now := time.Now()
	expirationAge := now.Add(-txmp.config.TTLDuration)
	if txmp.config.TTLDuration == 0 {
		expirationAge = time.Time{}
	}

	txmp.store.purgeExpiredTxs(expirationHeight, expirationAge)

	// purge old evicted transactions
	if txmp.config.TTLDuration > 0 {
		limit := now.Add(-txmp.config.TTLDuration)
		txmp.evictedTxs.Prune(limit)
	}
}

func (txmp *TxPool) notifyTxsAvailable() {
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
