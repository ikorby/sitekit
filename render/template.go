package render

import (
	"fmt"
	"html/template"
	"path"
)

func (r *Renderer) templateFor(layout, tmpl string) (*template.Template, error) {
	key := layout + "|" + tmpl

	if !r.reload {
		r.mu.RLock()
		t, ok := r.cache[key]
		r.mu.RUnlock()
		if ok {
			return t, nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.reload {
		if t, ok := r.cache[key]; ok {
			return t, nil
		}
	}

	layoutPath := path.Join(layoutsDir, layout)
	pagePath := path.Join(pagesDir, tmpl)

	t, err := template.New(path.Base(layoutPath)).
		Funcs(r.funcs).
		ParseFS(r.fsys, layoutPath, pagePath)
	if err != nil {
		return nil, fmt.Errorf("parse templates (layout=%q, page=%q): %w", layout, tmpl, err)
	}

	if !r.reload {
		r.cache[key] = t
	}

	return t, nil
}
