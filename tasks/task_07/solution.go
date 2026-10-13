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
	lruCachePtr := &LRUCache[K, V]{capacity: capacity}
	if capacity > 0 {
		lruCachePtr.items = make(map[K]*list.Element, capacity)
	}
	return lruCachePtr
}

func (lruCachePtr *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	lruCachePtr.mu.Lock()
	defer lruCachePtr.mu.Unlock()

	el, found := lruCachePtr.items[key]
	if !found {
		return value, false
	}
	lruCachePtr.ll.MoveToFront(el)
	return el.Value.(*entry[K, V]).value, true
}

func (lruCachePtr *LRUCache[K, V]) Set(key K, value V) {
	lruCachePtr.mu.Lock()
	defer lruCachePtr.mu.Unlock()

	if lruCachePtr.capacity <= 0 {
		return
	}

	if el, found := lruCachePtr.items[key]; found {
		el.Value.(*entry[K, V]).value = value
		return
	}

	if lruCachePtr.ll.Len() >= lruCachePtr.capacity {
		if backEl := lruCachePtr.ll.Back(); backEl != nil {
			delete(lruCachePtr.items, backEl.Value.(*entry[K, V]).key)
			lruCachePtr.ll.Remove(backEl)
		}
	}

	el := lruCachePtr.ll.PushFront(&entry[K, V]{key: key, value: value})
	lruCachePtr.items[key] = el
}
