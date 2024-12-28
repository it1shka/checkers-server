package multiplayer

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

func (d *SafeDict[K, V]) HasKey(key K) bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	_, exists := d.storage[key]
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

func (d *SafeDict[K, V]) IfExists(key K, action func(value V)) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	value, exists := d.storage[key]
	if !exists {
		return
	}
	action(value)
}
