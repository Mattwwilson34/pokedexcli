package internal

import (
	"sync"
	"time"
)

type Cache struct {
	entriesMu sync.Mutex
	entries   map[string]cacheEntry
	interval  time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{interval: interval, entries: make(map[string]cacheEntry)}
	go cache.reapLoop()
	return cache
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.entriesMu.Lock()
		for key, val := range c.entries {
			age := time.Since(val.createdAt)
			if age > c.interval {
				delete(c.entries, key)
			}
		}
		c.entriesMu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) error {
	c.entriesMu.Lock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.entriesMu.Unlock()
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.entriesMu.Lock()
	entry, ok := c.entries[key]
	c.entriesMu.Unlock()
	if !ok {
		return nil, false
	}
	return entry.val, true
}
