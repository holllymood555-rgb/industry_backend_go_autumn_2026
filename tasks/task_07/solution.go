package main

import (
    "container/list"
    "sync"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	ll       list.List
	items    map[K]*list.Element
	mu sync.Mutex
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
    c.mu.Lock()
    defer c.mu.Unlock()
	return &LRUCache[K, V]{
		capacity: capacity,
		items: make(map[K]*list.Element),
		ll: list.New(),
	}
}
	
func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
	var a *list.Element
	a, ok = c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.ll.MoveToBack(a)
	ent := a.Value.(entry[K, V])
	return ent.value, true
}
	
func (c *LRUCache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.capacity <= 0 {
        return
    }

    if a, ok := c.items[key]; ok {
        a.Value = entry[K, V]{key: key, value: value}
        return
    }

    if c.ll.Len() >= c.capacity {
        oldest := c.ll.Front()
        ent := oldest.Value.(entry[K, V])
        delete(c.items, ent.key)
        c.ll.Remove(oldest)
    }

    a := c.ll.PushBack(entry[K, V]{key: key, value: value})
    c.items[key] = a
}

type LRU[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
}
