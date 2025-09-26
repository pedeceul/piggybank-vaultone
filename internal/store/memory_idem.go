package store

import (
	"context"
	"sync"
)

type memoryIdem struct {
	mu sync.Mutex
	m  map[string]IdempotencyRecord
}

func NewMemoryIdempotencyStore() IdempotencyStore {
	return &memoryIdem{m: make(map[string]IdempotencyRecord)}
}

func (s *memoryIdem) Get(ctx context.Context, key string) (IdempotencyRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.m[key]
	return rec, ok, nil
}

func (s *memoryIdem) Set(ctx context.Context, key string, rec IdempotencyRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = rec
	return nil
}
