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
	if capacity <= 0 {
		capacity = 0
	}
	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
	}
}
func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	getElem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.ll.MoveToFront(getElem)

	return getElem.Value.(*entry[K, V]).value, true
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity == 0 {
		return
	}
	updEntry, ok := c.items[key]
	if ok {
		updEntry.Value.(*entry[K, V]).value = value
		return
	}

	if c.ll.Len() == c.capacity {
		link_del := c.ll.Back()
		delete(c.items, link_del.Value.(*entry[K, V]).key)
		c.ll.Remove(link_del)
	}

	newEntry := &entry[K, V]{
		key:   key,
		value: value,
	}
	//вставляю элемент в начало list и дублирую ссылку на него в мапы
	c.items[key] = c.ll.PushFront(newEntry)
}
