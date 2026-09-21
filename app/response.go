package app

import "net/http"

// Response is Sitekit's response boundary around net/http.ResponseWriter.
//
// It preserves the standard http.ResponseWriter contract while tracking the
// observable response lifecycle for the application layer.
type Response struct {
	writer    http.ResponseWriter
	status    int
	bytes     int64
	committed bool
}

// NewResponse creates a response boundary around w.
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{writer: w}
}

// Header returns the underlying response headers.
func (r *Response) Header() http.Header {
	return r.writer.Header()
}

// WriteHeader commits the response with status.
//
// A response can be committed only once. Subsequent calls are ignored,
// matching the net/http ResponseWriter contract.
func (r *Response) WriteHeader(status int) {
	if r.committed {
		return
	}

	r.writer.WriteHeader(status)
	r.status = status
	r.committed = true
}

// Write writes response data.
//
// The first write implicitly commits the response with status 200, matching
// the net/http ResponseWriter contract.
func (r *Response) Write(p []byte) (int, error) {
	if !r.committed {
		r.WriteHeader(http.StatusOK)
	}

	n, err := r.writer.Write(p)
	r.bytes += int64(n)
	return n, err
}

// Status returns the committed HTTP status code.
//
// It returns 0 when the response has not been committed yet.
func (r *Response) Status() int {
	return r.status
}

// BytesWritten returns the number of bytes successfully written to the
// underlying response writer.
func (r *Response) BytesWritten() int64 {
	return r.bytes
}

// Committed reports whether the response headers have been committed.
func (r *Response) Committed() bool {
	return r.committed
}

// Unwrap returns the underlying net/http response writer.
//
// It is the explicit escape hatch for APIs that require a concrete or
// optional net/http.ResponseWriter implementation.
func (r *Response) Unwrap() http.ResponseWriter {
	return r.writer
}
