package auth

import (
	"crypto/subtle"
	"net/http"
)

func RequireAPIKey(expected string) func(http.Handler) http.Handler {
	expBytes := []byte(expected)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := []byte(r.Header.Get("X-API-Key"))
			if len(got) != len(expBytes) || subtle.ConstantTimeCompare(got, expBytes) != 1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"Invalid or missing API key"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
