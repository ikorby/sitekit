package livereload

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ikorby/sitekit/config"
)

func StartWatcher(dir string, interval time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		var lastMod time.Time
		for {
			var currentMod time.Time
			_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					info, _ := d.Info()
					if info.ModTime().After(currentMod) {
						currentMod = info.ModTime()
					}
				}
				return nil
			})

			if lastMod.IsZero() {
				lastMod = currentMod
			} else if currentMod.After(lastMod) {
				lastMod = currentMod
				ch <- struct{}{}
			}
			time.Sleep(interval)
		}
	}()
	return ch
}

func SSEHandler(reloadCh <-chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-reloadCh:
				_, err := fmt.Fprintf(w, "data: reload\n\n")
				if err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}

func Middleware(cfg *config.Config) func(http.Handler) http.Handler {
	script := []byte(`
<!-- Sitekit Live Reload -->
<script>
	const evtSource = new EventSource("/__sitekit/livereload");
	evtSource.onmessage = function(event) {
		if (event.data === "reload") {
			window.location.reload();
		}
	};
</script>
</body>`)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg == nil || cfg.IsProduction() {
				next.ServeHTTP(w, r)
				return
			}

			rec := &responseInterceptor{ResponseWriter: w, script: script}
			next.ServeHTTP(rec, r)
		})
	}
}

type responseInterceptor struct {
	http.ResponseWriter
	script []byte
}

func (w *responseInterceptor) Write(b []byte) (int, error) {
	if strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		b = bytes.Replace(b, []byte("</body>"), w.script, 1)
	}
	return w.ResponseWriter.Write(b)
}
