package otel

import (
	"reflect"
	"sync"
)

type cache[T any] struct {
	mu    sync.RWMutex
	cache map[reflect.Type]T
}

func newCache[T any]() *cache[T] {
	return &cache[T]{
		cache: make(map[reflect.Type]T),
	}
}

func (c *cache[T]) get(t reflect.Type) (T, bool) {
	c.mu.RLock()
	v, ok := c.cache[t]
	c.mu.RUnlock()
	return v, ok
}

func (c *cache[T]) set(t reflect.Type, v T) {
	c.mu.Lock()
	c.cache[t] = v
	c.mu.Unlock()
}
