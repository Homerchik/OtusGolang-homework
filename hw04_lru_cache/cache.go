package hw04lrucache

type Key string

type Value struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	storedValue := Value{key: key, value: value}
	if item, exists := cache.items[key]; exists {
		item.Value = storedValue
		cache.queue.MoveToFront(item)
		return true
	}
	if cache.queue.Len() == cache.capacity {
		item := cache.queue.Back()
		cache.queue.Remove(item)
		delete(cache.items, item.Value.(Value).key)
	}
	newItem := cache.queue.PushFront(storedValue)
	cache.items[key] = newItem
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := cache.items[key]; exists {
		cache.queue.MoveToFront(item)
		return item.Value.(Value).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}
