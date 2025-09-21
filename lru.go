// Package lru implements a thread-safe Least Recently Used (LRU) cache.
// It supports generic key-value pairs and provides O(1) Get and Put operations.
package lru

import (
	"container/list"
	"errors"
	"sync"
)

// LRU is a thread-safe Least Recently Used cache with O(1) Get and Put.
type LRU[K comparable, V any] struct {
	cap     int
	mu      sync.RWMutex
	list    *list.List // holds *entry[K,V]
	idx     map[K]*list.Element
	onEvict func(key K, value V) // optional eviction callback
}

type entry[K comparable, V any] struct {
	key K
	val V
}

// NewLRU creates a new LRU cache with the specified capacity.
// Returns an error if capacity <= 0.
func NewLRU[K comparable, V any](capacity int) (*LRU[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	return &LRU[K, V]{
		cap:  capacity,
		list: list.New(),
		idx:  make(map[K]*list.Element, capacity),
	}, nil
}

// Get retrieves the value for the given key if present.
// Moves the accessed item to the front of the cache.
func (c *LRU[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var zero V
	if el, ok := c.idx[key]; ok {
		// Move to front under write lock
		c.mu.RUnlock()
		c.mu.Lock()
		c.list.MoveToFront(el)
		c.mu.Unlock()
		c.mu.RLock()
		return el.Value.(*entry[K, V]).val, true
	}
	return zero, false
}

// Put inserts or updates the value for the given key.
// If capacity is exceeded, evicts the least recently used item.
func (c *LRU[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.idx[key]; ok {
		el.Value.(*entry[K, V]).val = val
		c.list.MoveToFront(el)
		return
	}
	el := c.list.PushFront(&entry[K, V]{key: key, val: val})
	c.idx[key] = el
	if c.list.Len() > c.cap {
		tail := c.list.Back()
		if tail != nil {
			c.list.Remove(tail)
			kv := tail.Value.(*entry[K, V])
			delete(c.idx, kv.key)
			if c.onEvict != nil {
				c.onEvict(kv.key, kv.val)
			}
		}
	}
}

// Len returns the current number of items in the cache.
func (c *LRU[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}

// SetEvictionCallback sets the callback to be called when an item is evicted.
func (c *LRU[K, V]) SetEvictionCallback(fn func(key K, value V)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvict = fn
}
