package server

import (
	"net/http"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Revision is the git commit, filled at build time
var Revision = "[unknown]"

type Server struct {
	Router *chi.Mux

	DB   *db.DB
	Auth *auth.Verifier
}

func New() (*Server, error) {
	db, err := db.New(os.Getenv("DATABASE_URL"))
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
	// enable authentication for all routes (but don't require it)
	s.Router.Use(s.maybeAuth)

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
