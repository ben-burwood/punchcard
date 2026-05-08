package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ben-burwood/punchcard/internal/auth"
	"github.com/ben-burwood/punchcard/internal/config"
	"github.com/ben-burwood/punchcard/internal/handlers"
	"github.com/ben-burwood/punchcard/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	sessions := auth.NewStore()
	runs := store.NewJobRuns(db)

	mux := http.NewServeMux()

	mux.Handle("POST /api/punch", auth.RequireAPIKey(cfg.APIKey)(handlers.Punch(runs)))

	mux.Handle("GET /web/auth/status", auth.RequireSession(sessions)(handlers.Status()))
	mux.Handle("POST /web/auth/login", handlers.Login(cfg, sessions))
	mux.Handle("POST /web/auth/logout", handlers.Logout(sessions))

	mux.Handle("GET /web/jobs/running", auth.RequireSession(sessions)(handlers.Running(runs)))
	mux.Handle("GET /web/jobs/history", auth.RequireSession(sessions)(handlers.History(runs)))

	mux.Handle("GET /health", handlers.Health())

	mux.Handle("/", handlers.SPA(cfg.StaticDir))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
