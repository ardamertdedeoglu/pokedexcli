package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]CacheEntry
	mux          sync.Mutex
	interval     time.Duration
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		cacheEntries: make(map[string]CacheEntry),
		mux:          sync.Mutex{},
		interval:     interval,
	}
	go c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	c.cacheEntries[key] = CacheEntry{createdAt: time.Now(), val: val}
	c.mux.Unlock()
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	res, ok := c.cacheEntries[key]
	c.mux.Unlock()
	if !ok {
		return nil, ok
	}

	return res.val, ok

}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mux.Lock()
		now := time.Now()
		for key, entry := range c.cacheEntries {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.cacheEntries, key)
			}
		}
		c.mux.Unlock()
	}
}
