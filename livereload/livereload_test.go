package livereload_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ikorby/sitekit/config"
	"github.com/ikorby/sitekit/livereload"
)

func TestMiddleware_InjectsScriptInDevelopment(t *testing.T) {
	cfg := &config.Config{Env: config.Development}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>Hello World</body></html>"))
	})

	handler := livereload.Middleware(cfg)(next)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "<!-- Sitekit Live Reload -->") {
		t.Fatalf("expected script to be injected in development, got:\n%s", body)
	}
	if !strings.Contains(body, "new EventSource") {
		t.Fatalf("expected EventSource initialization, got:\n%s", body)
	}
}

func TestMiddleware_SkipsInjectionForNonHTML(t *testing.T) {
	cfg := &config.Config{Env: config.Development}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message": "Hello </body>"}`)) // Специально добавляем тег, чтобы проверить ложное срабатывание
	})

	handler := livereload.Middleware(cfg)(next)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, "<!-- Sitekit Live Reload -->") {
		t.Fatalf("script should not be injected into JSON responses, got:\n%s", body)
	}
}

func TestMiddleware_SkipsInjectionInProduction(t *testing.T) {
	cfg := &config.Config{Env: config.Production}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>Production Data</body></html>"))
	})

	handler := livereload.Middleware(cfg)(next)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, "<!-- Sitekit Live Reload -->") {
		t.Fatalf("script should NEVER be injected in production, got:\n%s", body)
	}
}

func TestSSEHandler_Headers(t *testing.T) {
	reloadCh := make(chan struct{})
	handler := livereload.SSEHandler(reloadCh)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/__sitekit/livereload", nil)

	// Отменяем контекст, чтобы горутина обработчика не зависла навечно
	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("Cache-Control") != "no-cache" {
		t.Errorf("expected Cache-Control no-cache, got %s", rec.Header().Get("Cache-Control"))
	}
}
