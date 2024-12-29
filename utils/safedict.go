package utils

import "sync"

type SafeDict[K comparable, V any] struct {
	mutex   sync.RWMutex
	storage map[K]V
}

func NewSafeDict[K comparable, V any]() *SafeDict[K, V] {
	return &SafeDict[K, V]{
		mutex:   sync.RWMutex{},
		storage: make(map[K]V),
	}
}

func (d *SafeDict[K, V]) Get(key K) (V, bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	value, ok := d.storage[key]
	return value, ok
}

func (d *SafeDict[K, V]) GetOrDefault(key K, defaultValue V) V {
	value, exists := d.Get(key)
	if !exists {
		return defaultValue
	}
	return value
}

func (d *SafeDict[K, V]) GetOrEmpty(key K) V {
	value, exists := d.Get(key)
	if !exists {
		var empty V
		return empty
	}
	return value
}

func (d *SafeDict[K, V]) HasKey(key K) bool {
	_, exists := d.Get(key)
	return exists
}

func (d *SafeDict[K, V]) Put(key K, value V) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.storage[key] = value
}

func (d *SafeDict[K, V]) Delete(key K) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.storage, key)
}

func (d *SafeDict[K, V]) IfExists(key K, do func(value V)) {
	value, exists := d.Get(key)
	if exists {
		do(value)
	}
}

func (d *SafeDict[K, V]) Clear() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	clear(d.storage)
}

func (d *SafeDict[K, V]) Keys() []K {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	keys := make([]K, len(d.storage))
	index := 0
	for key := range d.storage {
		keys[index] = key
		index++
	}
	return keys
}

func (d *SafeDict[K, V]) Values() []V {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	values := make([]V, len(d.storage))
	index := 0
	for _, value := range d.storage {
		values[index] = value
		index++
	}
	return values
}
