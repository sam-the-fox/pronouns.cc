package server

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/server/auth"
	"codeberg.org/pronounscc/pronouns.cc/backend/server/rate"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	chiprometheus "github.com/toshi0607/chi-prometheus"
)

// Revision is the git commit, filled at build time
var (
	Revision = "[unknown]"
	Tag      = "[unknown]"
)

// Repository is the URL of the git repository
const Repository = "https://codeberg.org/pronounscc/pronouns.cc"

type Server struct {
	Router *chi.Mux

	DB   *db.DB
	Auth *auth.Verifier
}

func New() (*Server, error) {
	db, err := db.New()
	if err != nil {
		return nil, err
	}

	s := &Server{
		Router: chi.NewMux(),

		DB:   db,
		Auth: auth.New(),
	}

	if os.Getenv("DEBUG") == "true" {
		s.Router.Use(middleware.Logger)
	}
	s.Router.Use(middleware.Recoverer)
	// add Sentry tracing handler
	s.Router.Use(s.sentry)

	// add CORS
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// Allow all methods normally used by the API
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// enable request latency tracking
	os.Setenv(chiprometheus.EnvChiPrometheusLatencyBuckets, "10,25,50,100,300,500,1000,5000")
	prom := chiprometheus.New("pronouns.cc")
	s.Router.Use(prom.Handler)
	prom.MustRegisterDefault()

	// enable authentication for all routes (but don't require it)
	s.Router.Use(s.maybeAuth)

	// rate limit handling
	// - base is 120 req/minute (2/s)
	// - keyed by Authorization header if valid token is provided, otherwise by IP
	// - returns rate limit reset info in error
	rateLimiter := rate.NewLimiter(120, time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			_, ok := ClaimsFromContext(r.Context())
			if token := r.Header.Get("Authorization"); ok && token != "" {
				return token, nil
			}

			ip, err := httprate.KeyByIP(r)
			return ip, err
		}),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			reset, _ := strconv.Atoi(w.Header().Get("X-RateLimit-Reset"))

			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, APIError{
				Code:           ErrTooManyRequests,
				Message:        errCodeMessages[ErrTooManyRequests],
				RatelimitReset: &reset,
			})
		}),
	)

	// set scopes
	// users
	rateLimiter.Scope("GET", "/users/*", 60)
	rateLimiter.Scope("PATCH", "/users/@me", 10)

	// members
	rateLimiter.Scope("GET", "/users/*/members", 60)
	rateLimiter.Scope("GET", "/users/*/members/*", 60)

	rateLimiter.Scope("POST", "/members", 10)
	rateLimiter.Scope("GET", "/members/*", 60)
	rateLimiter.Scope("PATCH", "/members/*", 20)
	rateLimiter.Scope("DELETE", "/members/*", 5)

	// auth
	rateLimiter.Scope("*", "/auth/*", 20)
	rateLimiter.Scope("*", "/auth/tokens", 10)
	rateLimiter.Scope("*", "/auth/invites", 10)
	rateLimiter.Scope("POST", "/auth/discord/*", 10)

	s.Router.Use(rateLimiter.Handler())

	// increment the total requests counter whenever a request is made
	s.Router.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			s.DB.TotalRequests.Inc()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	// return an API error for not found + method not allowed
	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, errCodeStatuses[ErrNotFound])
		render.JSON(w, r, APIError{
			Code:    ErrNotFound,
			Message: errCodeMessages[ErrNotFound],
		})
	})

	s.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, errCodeStatuses[ErrMethodNotAllowed])
		render.JSON(w, r, APIError{
			Code:    ErrMethodNotAllowed,
			Message: errCodeMessages[ErrMethodNotAllowed],
		})
	})

	return s, nil
}

type ctxKey int

const (
	ctxKeyClaims ctxKey = 1
)
