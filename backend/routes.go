package main

import (
	"github.com/go-chi/chi/v5"
	"gitlab.com/1f320/pronouns/backend/routes/auth"
	"gitlab.com/1f320/pronouns/backend/server"
)

// mountRoutes mounts all API routes on the server's router.
// they are all mounted under /v1/
func mountRoutes(s *server.Server) {
	s.Router.Route("/v1", func(r chi.Router) {
		auth.Mount(s, r)
	})
}
