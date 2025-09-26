package middlewarex

import (
	"net/http"
	"sync"
	"time"
)

// tokenBucket provides a simple in-memory token bucket per IP (dev only).
type tokenBucket struct {
	capacity     int
	refillTokens int
	refillEvery  time.Duration
	tokens       int
	lastRefill   time.Time
}

type rlState struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

// RateLimit creates a naive per-remote-IP token bucket.
func RateLimit(capacity, refillTokens int, refillEvery time.Duration) func(http.Handler) http.Handler {
	state := &rlState{buckets: make(map[string]*tokenBucket)}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			state.mu.Lock()
			b, ok := state.buckets[ip]
			if !ok {
				b = &tokenBucket{capacity: capacity, refillTokens: refillTokens, refillEvery: refillEvery, tokens: capacity, lastRefill: time.Now()}
				state.buckets[ip] = b
			}
			now := time.Now()
			if d := now.Sub(b.lastRefill); d >= b.refillEvery {
				refills := int(d / b.refillEvery)
				b.tokens += refills * b.refillTokens
				if b.tokens > b.capacity {
					b.tokens = b.capacity
				}
				b.lastRefill = now
			}
			if b.tokens <= 0 {
				state.mu.Unlock()
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"code":"rate_limited","message":"too many requests"}`))
				return
			}
			b.tokens--
			state.mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
