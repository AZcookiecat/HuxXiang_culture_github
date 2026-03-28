package app

import (
	"strings"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, ttl time.Duration)
	Delete(key string)
	DeletePrefix(prefix string)
}

type cacheItem struct {
	value     any
	expiresAt time.Time
}

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{items: make(map[string]cacheItem)}
}

func (c *InMemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.Delete(key)
		return nil, false
	}
	return item.value, true
}

func (c *InMemoryCache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := cacheItem{value: value}
	if ttl > 0 {
		item.expiresAt = time.Now().Add(ttl)
	}
	c.items[key] = item
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *InMemoryCache) DeletePrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.items {
		if strings.HasPrefix(key, prefix) {
			delete(c.items, key)
		}
	}
}
