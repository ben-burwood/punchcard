package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

const (
	SessionCookie = "punchcard_session"
	SessionTTL    = 24 * time.Hour
)

type Store struct {
	mu       sync.Mutex
	sessions map[string]time.Time
}

func NewStore() *Store {
	return &Store{sessions: map[string]time.Time{}}
}

func (s *Store) Create() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for tok, exp := range s.sessions {
		if exp.Before(now) {
			delete(s.sessions, tok)
		}
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	token := hex.EncodeToString(buf)
	s.sessions[token] = now.Add(SessionTTL)
	return token
}

func (s *Store) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

func (s *Store) Validate(token string) bool {
	if token == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.sessions[token]
	if !ok {
		return false
	}
	if exp.Before(time.Now()) {
		delete(s.sessions, token)
		return false
	}
	return true
}

func TokenFromRequest(r *http.Request) string {
	c, err := r.Cookie(SessionCookie)
	if err != nil {
		return ""
	}
	return c.Value
}

func WriteCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(SessionTTL.Seconds()),
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   0,
	})
}

func RequireSession(s *Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.Validate(TokenFromRequest(r)) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"Unauthorized"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
