package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}
type entry[K comparable, V any] struct {
	key   K
	value V
}
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		items: make(map[K]*list.Element),
	}
}
func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	i, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.ll.MoveToBack(i)
	res := i.Value.(entry[K, V])
	return res.value, true
}
func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity <= 0 {
		return
	}

	if i, ok := c.items[key]; ok {
		res := i.Value.(entry[K, V])
		res.value = value
		i.Value = res
		return
	}

	if c.ll.Len() >= c.capacity {
		first := c.ll.Front()
		res := first.Value.(entry[K, V])
		delete(c.items, res.key)
		c.ll.Remove(first)
	}
	
	newItem := c.ll.PushBack(entry[K, V]{key: key, value: value})
	c.items[key] = newItem		
}

