package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ikorby/sitekit/page"
	"github.com/ikorby/sitekit/render"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	renderer *render.Renderer
}

func newContext(w http.ResponseWriter, r *http.Request, renderer *render.Renderer) *Context {
	return &Context{Writer: w, Request: r, renderer: renderer}
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}

type ctxKey string

func WithValue(r *http.Request, key string, val any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ctxKey(key), val))
}

func (c *Context) Value(key string) any {
	return c.Request.Context().Value(ctxKey(key))
}

func (c *Context) Param(name string) string {
	return c.Request.PathValue(name)
}

func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

func (c *Context) Render(status int, p page.Page) error {
	if c.renderer == nil {
		return errNoRenderer
	}
	return c.renderer.Render(c.Writer, status, p)
}

func (c *Context) JSON(status int, v any) error {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(status)
	return json.NewEncoder(c.Writer).Encode(v)
}

func (c *Context) Redirect(url string, status int) {
	http.Redirect(c.Writer, c.Request, url, status)
}
