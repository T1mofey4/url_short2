package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"url_shortener/internal/storage"
)

func TestHealtzHandler(t *testing.T) {
	ctx := context.Background()
	_ = godotenv.Load()

	pg_dsn := os.Getenv("PG_DSN")
	if pg_dsn == "" {
		t.Fatalf("not found pg_dsn in env")
	}

	// database initialization
	db, err := storage.NewDB(ctx, pg_dsn)
	if err != nil {
		t.Fatalf("failed to connect ot database: %v", err)
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

	req := httptest.NewRequest(http.MethodGet, "/healtz", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}

	want := `server healtz OK!`
	if got := rec.Body.String(); got != want {
		t.Errorf("expected body %q, got %q", want, got)
	}
}
