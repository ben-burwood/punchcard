package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ben-burwood/punchcard/internal/auth"
	"github.com/ben-burwood/punchcard/internal/config"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Status() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})
}

func Login(cfg config.Config, sessions *auth.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body loginRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}
		if body.Username != cfg.DashUser || body.Password != cfg.DashPassword {
			writeError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		token := sessions.Create()
		auth.WriteCookie(w, token)
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})
}

func Logout(sessions *auth.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := auth.TokenFromRequest(r); token != "" {
			sessions.Delete(token)
		}
		auth.ClearCookie(w)
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})
}
