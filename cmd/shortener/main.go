package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"url_shortener/internal/logger"
	"url_shortener/internal/storage"
)

func main() {
	ctx := context.Background()

	log := logger.New()
	_ = godotenv.Load()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "3000"
	}

	host := os.Getenv("HTTP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	pg_dsn := os.Getenv("PG_DSN")
	if host == "" {
		log.Error("not found pg_dsn in env")
	}

	// database initialization
	db, err := storage.NewDB(ctx, pg_dsn)
	if err != nil {
		log.Error("failed to connect ot database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Get("/healtz", func(w http.ResponseWriter, r *http.Request) {
		err = db.Ping(ctx)
		if err != nil {
			http.Error(w, "database unavialible", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("server healtz OK!"))
	})

	addr := host + ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("http server started", "addr", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "err", err)
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Error("graceful shutdown error", "error", err)
	}

	log.Info("server exited")

}
