package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	apiTSLayout = "2006-01-02T15:04:05.000000"
)

func formatWithZ(t time.Time) string {
	return t.UTC().Format(apiTSLayout) + "Z"
}

func formatNoZ(t time.Time) string {
	return t.UTC().Format(apiTSLayout)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
