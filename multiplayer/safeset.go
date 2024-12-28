package multiplayer

import "sync"

type SafeSet[T comparable] struct {
	mutex   sync.RWMutex
	storage map[T]bool
}

func NewSafeSet[T comparable]() *SafeSet[T] {
	return &SafeSet[T]{
		mutex:   sync.RWMutex{},
		storage: make(map[T]bool),
	}
}

func (s *SafeSet[T]) Add(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.storage[value] = true
}

func (s *SafeSet[T]) Remove(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.storage, value)
}

func (s *SafeSet[T]) Has(value T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, exists := s.storage[value]
	return exists
}

func (s *SafeSet[T]) Values() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	values := make([]T, len(s.storage))
	index := 0
	for key := range s.storage {
		values[index] = key
		index++
	}
	return values
}

// Like Values(), but clears collection afterwards
func (s *SafeSet[T]) Eject() []T {
	s.mutex.Lock()
	defer func() {
		clear(s.storage)
		s.mutex.Unlock()
	}()
	values := make([]T, len(s.storage))
	index := 0
	for key := range s.storage {
		values[index] = key
		index++
	}
	return values
}

func (s *SafeSet[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	clear(s.storage)
}
