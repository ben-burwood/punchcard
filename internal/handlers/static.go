package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SPA(staticDir string) http.Handler {
	if info, err := os.Stat(staticDir); err != nil || !info.IsDir() {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	abs, err := filepath.Abs(staticDir)
	if err != nil {
		abs = staticDir
	}
	fs := http.FileServer(http.Dir(abs))
	indexPath := filepath.Join(abs, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		full := filepath.Join(abs, clean)

		if !strings.HasPrefix(full, abs) {
			http.NotFound(w, r)
			return
		}

		if r.URL.Path == "/" {
			http.ServeFile(w, r, indexPath)
			return
		}

		if info, err := os.Stat(full); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, indexPath)
	})
}
