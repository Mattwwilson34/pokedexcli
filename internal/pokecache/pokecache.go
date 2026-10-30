package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu       sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
	hits     int
	misses   int
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
		c.mu.Lock()
		for key, val := range c.entries {
			age := time.Since(val.createdAt)
			if age > c.interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		c.misses += 1
		return nil, false
	}
	c.hits += 1
	return entry.val, true
}

func (c *Cache) PrintCacheStats() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println()
	fmt.Printf("Number of Cache Entries: %v", len(c.entries))
	fmt.Println()
	fmt.Printf("Number of Cache Hits: %v", c.hits)
	fmt.Println()
	fmt.Printf("Number of Cache Misses: %v", c.misses)
	fmt.Println()

	return nil
}
