# go-lru-cache
Production-ready thread-safe LRU cache implementation in Go

# Go LRU Cache

A **generic**, **thread-safe**, and **high-performance** Least Recently Used (LRU) cache implementation in Go.

This package offers an LRU cache with:
- O(1) time complexity for `Get` and `Put`.
- Thread safety using `sync.RWMutex`.
- Optional eviction callback.
- Generic support for any comparable key and value types.

## 🚀 Features
✔ Efficient LRU cache using `container/list` and `map`.  
✔ Supports concurrent reads and writes.  
✔ Custom eviction handler.  
✔ Easy to integrate in any Go application.

## 📦 Installation

```bash
go get github.com/Prajwal306/go-lru-cache
```
To use in your code:  
```bash
import "github.com/Prajwal306/go-lru-cache/lru"
```
📖 Usage Example
```
package main

import (
	"fmt"
	"log"

	"github.com/Prajwal306/go-lru-cache/lru"
)

func main() {
	cache, err := lru.NewLRU 
	if err != nil {
		log.Fatal(err)
	}

	cache.SetEvictionCallback(func(key string, value int) {
		fmt.Printf("Evicted: %s -> %d\n", key, value)
	})

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	fmt.Println(cache.Get("a")) // 1, true

	cache.Put("d", 4) // Evicts "b"

	if _, ok := cache.Get("b"); !ok {
		fmt.Println("b has been evicted")
	}
}
```
