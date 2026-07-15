package seo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ikorby/sitekit/config"
)

type RobotsRule struct {
	UserAgent string
	Disallow  []string
	Allow     []string
}

func Robots(cfg *config.Config, rules ...RobotsRule) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder

		if cfg.IsDevelopment() {
			_, _ = b.WriteString("User-agent: *\nDisallow: /\n")
		} else {
			if len(rules) == 0 {
				rules = []RobotsRule{{UserAgent: "*"}}
			}
			for _, rule := range rules {
				ua := rule.UserAgent
				if ua == "" {
					ua = "*"
				}
				_, _ = fmt.Fprintf(&b, "User-agent: %s\n", ua)
				for _, d := range rule.Disallow {
					_, _ = fmt.Fprintf(&b, "Disallow: %s\n", d)
				}
				for _, a := range rule.Allow {
					_, _ = fmt.Fprintf(&b, "Allow: %s\n", a)
				}
				_, _ = b.WriteString("\n")
			}
			if base := strings.TrimRight(cfg.BaseURL, "/"); base != "" {
				_, _ = fmt.Fprintf(&b, "Sitemap: %s/sitemap.xml\n", base)
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(b.String()))
	})
}
