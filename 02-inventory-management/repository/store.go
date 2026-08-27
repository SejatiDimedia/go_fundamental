package repository

import "sync"

type MemoryStore[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		data: make(map[string]T),
	}
}

func (s *MemoryStore[T]) Save(id string, item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = item
	return nil
}

func (s *MemoryStore[T]) FindByID(id string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, exists := s.data[id]
	return item, exists
}

func (s *MemoryStore[T]) GetAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]T, 0, len(s.data))
	for _, item := range s.data {
		items = append(items, item)
	}
	return items
}

func (s *MemoryStore[T]) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
}
