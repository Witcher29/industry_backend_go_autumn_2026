package main

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	cachePtr := &Cache[K, V]{capacity: capacity}
	if capacity > 0 {
		cachePtr.items = make(map[K]V, capacity)
	}
	return cachePtr
}

func (cachePtr *Cache[K, V]) Get(k K) (v V, ok bool) {
	v, ok = cachePtr.items[k]
	return v, ok
}

func (cachePtr *Cache[K, V]) Set(k K, v V) bool {
	if cachePtr.capacity <= 0 {
		return false
	}
	if _, exists := cachePtr.items[k]; exists {
		cachePtr.items[k] = v
		return true
	}
	if len(cachePtr.items) >= cachePtr.capacity {
		return false
	}
	cachePtr.items[k] = v
	return true
}
