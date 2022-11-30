package v2

import (
	"container/list"
	"time"

	tmsync "github.com/tendermint/tendermint/libs/sync"
	"github.com/tendermint/tendermint/types"
)

// LRUTxCache maintains a thread-safe LRU cache of raw transactions. The cache
// only stores the hash of the raw transaction.
type LRUTxCache struct {
	mtx      tmsync.Mutex
	size     int
	cacheMap map[types.TxKey]*list.Element
	list     *list.List
}

func NewLRUTxCache(cacheSize int) *LRUTxCache {
	return &LRUTxCache{
		size:     cacheSize,
		cacheMap: make(map[types.TxKey]*list.Element, cacheSize),
		list:     list.New(),
	}
}

// GetList returns the underlying linked-list that backs the LRU cache. Note,
// this should be used for testing purposes only!
func (c *LRUTxCache) GetList() *list.List {
	return c.list
}

func (c *LRUTxCache) Reset() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cacheMap = make(map[types.TxKey]*list.Element, c.size)
	c.list.Init()
}

func (c *LRUTxCache) Push(txKey types.TxKey) bool {
	if c.size == 0 {
		return true
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	moved, ok := c.cacheMap[txKey]
	if ok {
		c.list.MoveToBack(moved)
		return false
	}

	if c.list.Len() >= c.size {
		front := c.list.Front()
		if front != nil {
			frontKey := front.Value.(types.TxKey)
			delete(c.cacheMap, frontKey)
			c.list.Remove(front)
		}
	}

	e := c.list.PushBack(txKey)
	c.cacheMap[txKey] = e

	return true
}

func (c *LRUTxCache) Remove(txKey types.TxKey) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	e := c.cacheMap[txKey]
	delete(c.cacheMap, txKey)

	if e != nil {
		c.list.Remove(e)
	}
}

func (c *LRUTxCache) Has(txKey types.TxKey) bool {
	if c.size == 0 {
		return false
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	_, ok := c.cacheMap[txKey]
	return ok
}

type evictedTxInfo struct {
	timeEvicted time.Time
	priority    int64
	gasWanted   int64
	sender      string
	peers       map[uint16]bool
}
