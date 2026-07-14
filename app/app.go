package app

import (
	"errors"
	"net/http"
)

type HandlerFunc func(c *Context) error

type ErrorHandler func(c *Context, err error)

var errNoRenderer = errors.New("app: renderer is not configured")

func defaultErrorHandler(c *Context, err error) {
	if c.W.Header().Get("Content-Type") != "" {
		return
	}
	http.Error(c.W, "Internal Server Error", http.StatusInternalServerError)
}

func (a *App) Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := newContext(w, r, a.Renderer)

		if err := h(c); err != nil {
			eh := a.ErrorHandler
			if eh == nil {
				eh = defaultErrorHandler
			}
			eh(c, err)
		}
	}
}
