package cache

import (
	"sync"
	"time"
)

// Cache is a simple in-memory cache with TTL support.
type Cache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// New creates a new Cache with the given default TTL.
func New() *Cache {
	return &Cache{
		items: make(map[string]*cacheItem),
	}
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		c.Delete(key)
		return nil, false
	}
	return item.value, true
}

// Set stores a value in the cache with the given TTL.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheItem)
}