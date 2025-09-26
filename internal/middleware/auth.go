package middlewarex

import (
	"net/http"
	"strings"
)

// APIKeyAuth returns a middleware that enforces a static API key when enabled.
// Exempts /healthz and /readyz. Optionally exempts webhook if sharedSecret matches X-Webhook-Secret header.
func APIKeyAuth(enabled bool, apiKey string, webhookSharedSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !enabled {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := r.URL.Path
			if p == "/healthz" || p == "/readyz" {
				next.ServeHTTP(w, r)
				return
			}
			if p == "/v1/webhooks/payment_event" && webhookSharedSecret != "" {
				if r.Header.Get("X-Webhook-Secret") == webhookSharedSecret {
					next.ServeHTTP(w, r)
					return
				}
			}
			provided := extractAPIKey(r)
			if provided == "" || provided != apiKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"code":"unauthorized","message":"invalid or missing API key"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractAPIKey(r *http.Request) string {
	if v := r.Header.Get("X-API-Key"); v != "" {
		return v
	}
	if v := r.Header.Get("Api-Key"); v != "" {
		return v
	}
	if v := r.Header.Get("Authorization"); v != "" {
		if strings.HasPrefix(strings.ToLower(v), "bearer ") {
			return strings.TrimSpace(v[7:])
		}
		return v
	}
	return ""
}
