package storage

import (
	"fmt"
	"sync"
)

type Memory[T any] struct {
	mu    sync.RWMutex
	store map[string]T
}

func NewMemory[T any]() Operations[T] {
	return &Memory[T]{
		mu:    sync.RWMutex{},
		store: make(map[string]T),
	}
}

func (m *Memory[T]) Store(k string, v T) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		return fmt.Errorf("store is not initialized")
	}

	m.store[k] = v
	return nil
}

func (m *Memory[T]) Get(k string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var v T
	if m.store == nil {
		return v, fmt.Errorf("store is not initialized")
	}

	v, ok := m.store[k]
	if !ok {
		return v, fmt.Errorf("key not found")
	}

	return v, nil
}

func (m *Memory[T]) Delete(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		return fmt.Errorf("store is not initialized")
	}

	delete(m.store, k)
	return nil
}

func (m *Memory[T]) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		return nil
	}

	for k := range m.store {
		delete(m.store, k)
	}

	m.store = nil
	return nil
}
