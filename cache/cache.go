package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, duration time.Duration)
	Delete(key string)
}

type MemCache struct {
	cache map[string]*cacheItem
	mu    *sync.RWMutex
}

type cacheItem struct {
	value  interface{}
	expiry time.Time
}

func New(mu *sync.RWMutex) *MemCache {
	return &MemCache{cache: make(map[string]*cacheItem), mu: mu}
}

func (c *MemCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("cache error: key '%s' not found", key)
	}

	// Check if the value is expired
	if time.Now().After(item.expiry) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.cache, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, fmt.Errorf("cache error: key '%s' is expired", key)
	}

	return item.value, nil
}

func (c *MemCache) Set(key string, val interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &cacheItem{
		value:  val,
		expiry: time.Now().Add(duration),
	}
}

func (c *MemCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}
