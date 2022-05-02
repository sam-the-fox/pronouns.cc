package server

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/1f320/pronouns/backend/db"
	"gitlab.com/1f320/pronouns/backend/server/auth"
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

	return s, nil
}

type ctxKey int

const (
	ctxKeyClaims ctxKey = 1
)
