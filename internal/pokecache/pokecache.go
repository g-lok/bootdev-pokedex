package pokecache

import (
	"errors"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

func NewCache(interval time.Duration) (*Cache, error) {
	if interval <= 0 {
		return nil, errors.New("interval must be higher than 0")
	}
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	// reapLoop()
	go cache.reapLoop()

	return cache, nil
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{time.Now(), val}
	c.entries[key] = entry
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, ok := c.entries[key]
	defer c.mu.Unlock()
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop() // not necessary in go>=v1.23 but play it safe

	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}
