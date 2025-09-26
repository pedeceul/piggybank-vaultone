package store

import (
	"context"
	"time"
)

// IdempotencyRecord represents a cached response for an idempotent POST.
type IdempotencyRecord struct {
	StatusCode int
	Body       []byte
	Expiry     time.Time
	ReqHash    string
}

// IdempotencyStore is the contract for caching idempotent POST responses.
type IdempotencyStore interface {
	Get(ctx context.Context, key string) (IdempotencyRecord, bool, error)
	Set(ctx context.Context, key string, rec IdempotencyRecord) error
}
