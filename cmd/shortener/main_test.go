package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHealtzHandler(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/healtz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	req := httptest.NewRequest(http.MethodGet, "/healtz", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}

	want := `{"status":"ok"}`
	if got := rec.Body.String(); got != want {
		t.Errorf("expected body %q, got %q", want, got)
	}
}
