package lru

import (
	"sync"
	"testing"
)

// TestLRUBasic ensures Put/Get and eviction ordering works.
func TestLRUBasic(t *testing.T) {
	cache, err := NewLRU[int, string](2) // Fix: Initialize cache with capacity
	if err != nil {
		t.Fatal(err)
	}

	cache.Put(1, "one")
	cache.Put(2, "two")

	if val, ok := cache.Get(1); !ok || val != "one" {
		t.Errorf("expected one, got %v", val)
	}

	cache.Put(3, "three") // should evict key 2 (least recently used)

	if _, ok := cache.Get(2); ok {
		t.Errorf("expected key 2 to be evicted")
	}
}

// TestEvictionCallback verifies eviction callback works correctly.
func TestEvictionCallback(t *testing.T) {
	cache, _ := NewLRU[int, string](1) // Fix: Initialize cache with capacity
	evicted := false

	cache.SetEvictionCallback(func(k int, v string) {
		evicted = true
		if k != 1 || v != "one" {
			t.Errorf("unexpected eviction: %d -> %s", k, v)
		}
	})

	cache.Put(1, "one")
	cache.Put(2, "two") // should trigger eviction of key 1

	if !evicted {
		t.Errorf("eviction callback not triggered")
	}
}

// TestConcurrency ensures thread safety under parallel load.
func TestConcurrency(t *testing.T) {
	cache, _ := NewLRU[int, int](100) // Fix: Initialize cache with capacity
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Put(i, i)
			cache.Get(i)
		}(i)
	}

	wg.Wait()

	if cache.Len() > 100 {
		t.Errorf("cache size exceeded capacity: got %d", cache.Len())
	}
}

// TestOverwrite verifies that overwriting an existing key updates the value.
func TestOverwrite(t *testing.T) {
	cache, _ := NewLRU[int, string](2)

	cache.Put(1, "one")
	cache.Put(1, "uno") // Overwrite key 1

	if val, ok := cache.Get(1); !ok || val != "uno" {
		t.Errorf("expected uno, got %v", val)
	}
}

// TestEvictionOrder verifies that the least recently used item is evicted.
func TestEvictionOrder(t *testing.T) {
	cache, _ := NewLRU[int, string](2)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Get(1)          // Access key 1 to make it recently used
	cache.Put(3, "three") // Should evict key 2

	if _, ok := cache.Get(2); ok {
		t.Errorf("expected key 2 to be evicted")
	}

	if val, ok := cache.Get(1); !ok || val != "one" {
		t.Errorf("expected one, got %v", val)
	}
}

// TestEmptyCache ensures Get on an empty cache returns false.
func TestEmptyCache(t *testing.T) {
	cache, _ := NewLRU[int, string](2)

	if _, ok := cache.Get(1); ok {
		t.Errorf("expected cache miss for key 1")
	}
}

// TestZeroCapacity ensures creating a cache with zero capacity fails.
func TestZeroCapacity(t *testing.T) {
	if _, err := NewLRU[int, string](0); err == nil {
		t.Errorf("expected error for zero capacity cache")
	}
}
