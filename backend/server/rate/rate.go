package rate

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/httprate"
	"github.com/gobwas/glob"
)

type Limiter struct {
	scopes         []*scopedLimiter
	defaultLimiter func(http.Handler) http.Handler

	windowLength time.Duration
	options      []httprate.Option

	wildcardScopes []*scopedLimiter
}

type scopedLimiter struct {
	Method, Pattern string

	glob    glob.Glob
	handler func(http.Handler) http.Handler
}

func NewLimiter(defaultLimit int, windowLength time.Duration, options ...httprate.Option) *Limiter {
	return &Limiter{
		windowLength: windowLength,
		options:      options,
	}
}

func (l *Limiter) Scope(method, pattern string, requestLimit int) error {
	handler := httprate.Limit(requestLimit, l.windowLength, l.options...)

	g, err := glob.Compile("/v*"+pattern, '/')
	if err != nil {
		return err
	}

	if method == "*" {
		l.wildcardScopes = append(l.wildcardScopes, &scopedLimiter{method, pattern, g, handler})
	} else {
		l.scopes = append(l.scopes, &scopedLimiter{method, pattern, g, handler})
	}
	return nil
}

func (l *Limiter) Handler() func(http.Handler) http.Handler {
	sort.Slice(l.scopes, func(i, j int) bool {
		len1 := len(strings.Split(l.scopes[i].Pattern, "/"))
		len2 := len(strings.Split(l.scopes[j].Pattern, "/"))

		return len1 > len2
	})
	l.scopes = append(l.scopes, l.wildcardScopes...)

	return l.handle
}

func (l *Limiter) handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, s := range l.scopes {
			if (r.Method == s.Method || s.Method == "*") && s.glob.Match(r.URL.Path) {
				bucket := s.Pattern
				if s.Method != "*" {
					bucket = s.Method + " " + s.Pattern
				}
				w.Header().Set("X-RateLimit-Bucket", bucket)

				s.handler(next).ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("X-RateLimit-Bucket", "/")
		l.defaultLimiter(next).ServeHTTP(w, r)
	})
}
