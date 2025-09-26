package middlewarex

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/vaultone/api/internal/store"
)

// Simple in-memory idempotency cache for dev; Phase 3 will get PG-backed store.
type cachedResponse struct {
	StatusCode int
	Body       []byte
	Expiry     time.Time
	ReqHash    string
}

type idemState struct {
	mu    sync.Mutex
	items map[string]cachedResponse
}

func IdempotencyWithStore(ttl time.Duration, st store.IdempotencyStore) func(http.Handler) http.Handler {
	state := &idemState{items: make(map[string]cachedResponse)} // fallback for nil store
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			// read body for hashing and replay
			bodyBytes, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			// restore body for downstream
			r.Body = io.NopCloser(bytesReader(bodyBytes))

			h := sha256.Sum256(bodyBytes)
			reqHash := hex.EncodeToString(h[:])

			state.mu.Lock()
			if st != nil {
				if cr, ok, _ := st.Get(r.Context(), key); ok && time.Now().Before(cr.Expiry) {
					if cr.ReqHash != reqHash {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusConflict)
						_, _ = w.Write([]byte(`{"code":"idempotency_conflict","message":"request payload hash mismatch for key"}`))
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(cr.StatusCode)
					_, _ = w.Write(cr.Body)
					return
				}
			}
			if cr, ok := state.items[key]; ok && time.Now().Before(cr.Expiry) {
				// if hash mismatch, return 409 per spec
				if cr.ReqHash != reqHash {
					state.mu.Unlock()
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusConflict)
					_, _ = w.Write([]byte(`{"code":"idempotency_conflict","message":"request payload hash mismatch for key"}`))
					return
				}
				// replay cached response
				state.mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(cr.StatusCode)
				_, _ = w.Write(cr.Body)
				return
			}
			state.mu.Unlock()

			rec := &responseRecorder{ResponseWriter: w}
			next.ServeHTTP(rec, r)

			// cache successful (2xx) or 409 responses
			if (rec.status >= 200 && rec.status < 300) || rec.status == http.StatusConflict {
				if st != nil {
					_ = st.Set(context.Background(), key, store.IdempotencyRecord{StatusCode: rec.status, Body: rec.buf, Expiry: time.Now().Add(ttl), ReqHash: reqHash})
				} else {
					state.mu.Lock()
					state.items[key] = cachedResponse{StatusCode: rec.status, Body: rec.buf, Expiry: time.Now().Add(ttl), ReqHash: reqHash}
					state.mu.Unlock()
				}
			}
		})
	}
}

// Idempotency is a convenience wrapper using in-memory store.
func Idempotency(ttl time.Duration) func(http.Handler) http.Handler {
	return IdempotencyWithStore(ttl, nil)
}

type responseRecorder struct {
	http.ResponseWriter
	buf    []byte
	status int
}

func (r *responseRecorder) WriteHeader(code int) { r.status = code; r.ResponseWriter.WriteHeader(code) }
func (r *responseRecorder) Write(p []byte) (int, error) {
	r.buf = append(r.buf, p...)
	return r.ResponseWriter.Write(p)
}

// bytesReader wraps a byte slice to satisfy io.Reader without copying.
func bytesReader(b []byte) *bytesReaderT { return &bytesReaderT{b: b} }

type bytesReaderT struct{ b []byte }

func (r *bytesReaderT) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}
