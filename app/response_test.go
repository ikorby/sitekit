package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	if r.Committed() {
		t.Fatal("new response is committed")
	}
	if got := r.Status(); got != 0 {
		t.Fatalf("initial status = %d, want 0", got)
	}

	r.Header().Set("X-Test", "ok")
	r.WriteHeader(http.StatusCreated)
	r.WriteHeader(http.StatusNoContent)

	if !r.Committed() {
		t.Fatal("response is not committed")
	}
	if got := r.Status(); got != http.StatusCreated {
		t.Fatalf("status = %d, want %d", got, http.StatusCreated)
	}
	if got := w.Code; got != http.StatusCreated {
		t.Fatalf("recorder status = %d, want %d", got, http.StatusCreated)
	}
	if got := r.Header().Get("X-Test"); got != "ok" {
		t.Fatalf("header = %q, want %q", got, "ok")
	}
}

func TestResponseWriteImplicitlyCommitsOK(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	n, err := r.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 5 {
		t.Fatalf("Write() bytes = %d, want 5", n)
	}
	if !r.Committed() {
		t.Fatal("response is not committed")
	}
	if got := r.Status(); got != http.StatusOK {
		t.Fatalf("status = %d, want %d", got, http.StatusOK)
	}
	if got := r.BytesWritten(); got != 5 {
		t.Fatalf("bytes = %d, want 5", got)
	}
	if got := w.Body.String(); got != "hello" {
		t.Fatalf("body = %q, want %q", got, "hello")
	}
}

func TestResponseCountsWrittenBytesOnWriteError(t *testing.T) {
	w := &failingWriter{err: errWrite}
	r := NewResponse(w)

	n, err := r.Write([]byte("hello"))
	if err != errWrite {
		t.Fatalf("Write() error = %v, want %v", err, errWrite)
	}
	if n != 2 {
		t.Fatalf("Write() bytes = %d, want 2", n)
	}
	if got := r.BytesWritten(); got != 2 {
		t.Fatalf("bytes = %d, want 2", got)
	}
}

func TestResponseUnwrap(t *testing.T) {
	w := httptest.NewRecorder()
	r := NewResponse(w)

	if got := r.Unwrap(); got != w {
		t.Fatalf("Unwrap() = %T, want original writer", got)
	}
}

type failingWriter struct {
	err error
}

func (w *failingWriter) Header() http.Header {
	return make(http.Header)
}

func (w *failingWriter) WriteHeader(int) {}

func (w *failingWriter) Write(p []byte) (int, error) {
	return 2, w.err
}

var errWrite = testError("write failed")

type testError string

func (e testError) Error() string { return string(e) }
