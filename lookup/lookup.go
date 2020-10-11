// description: lookup
//
// @author: xwc1125
// @date: 2019/11/20
package lookup

import (
	"github.com/allegro/bigcache/v2"
	"sync"
	"time"
)

// hash ==>
type Lookup struct {
	mu         sync.RWMutex
	cache      *bigcache.BigCache
	allKeyList map[string]struct{}
}

func NewLookup(txConfig *PoolConfig, OnRemoveWithReason func(key string, entry []byte, reason bigcache.RemoveReason)) *Lookup {
	cache, _ := bigcache.NewBigCache(getCacheConfig(txConfig, OnRemoveWithReason))
	return &Lookup{
		cache:      cache,
		allKeyList: map[string]struct{}{},
	}
}

func getCacheConfig(txConfig *PoolConfig, OnRemoveWithReason func(key string, entry []byte, reason bigcache.RemoveReason)) bigcache.Config {
	//bigcache.DefaultConfig(10 * time.Minute)
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 16384,

		// time after which entry can be evicted
		LifeWindow: txConfig.TxLifeTime * time.Minute, // 过期时间，10分钟后，自动删除

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		//CleanWindow: 5 * time.Minute,
		CleanWindow: txConfig.TxTaskTime * time.Second,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 1024,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: OnRemoveWithReason,
	}
	return config
}

func (tl *Lookup) Len() uint64 {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return uint64(len(tl.allKeyList))
}

func (tl *Lookup) Exist(hash string) bool {
	//tl.mu.RLock()
	//defer tl.mu.RUnlock()
	_, err := tl.cache.Get(hash)
	if err != nil {
		return false
	}
	return true
}
func (tl *Lookup) Get(hash string) ([]byte, error) {
	//tl.mu.RLock()
	//defer tl.mu.RUnlock()
	return tl.cache.Get(hash)
}

// 添加hash==> 节点peer
func (tl *Lookup) Add(hash string, data []byte) {
	tl.mu.Lock()
	tl.allKeyList[hash] = struct{}{}
	tl.mu.Unlock()
	tl.cache.Set(hash, data)
}

func (tl *Lookup) Del(hash string) {
	flag := false
	tl.mu.Lock()
	if _, ok := tl.allKeyList[hash]; ok {
		flag = ok
	}
	if flag {
		delete(tl.allKeyList, hash)
	}
	tl.mu.Unlock()
	if flag {
		tl.cache.Delete(hash)
	}
}

func (tl *Lookup) GetAllKeys() []string {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	var allKey []string
	for k, _ := range tl.allKeyList {
		allKey = append(allKey, k)
	}
	return allKey
}
