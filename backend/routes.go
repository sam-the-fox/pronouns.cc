package main

import (
	"codeberg.org/u1f320/pronouns.cc/backend/routes/auth"
	"codeberg.org/u1f320/pronouns.cc/backend/routes/bot"
	"codeberg.org/u1f320/pronouns.cc/backend/routes/member"
	"codeberg.org/u1f320/pronouns.cc/backend/routes/meta"
	"codeberg.org/u1f320/pronouns.cc/backend/routes/user"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
)

// mountRoutes mounts all API routes on the server's router.
// they are all mounted under /v1/
func mountRoutes(s *server.Server) {
	// future-proofing for API versions
	s.Router.Route("/v1", func(r chi.Router) {
		auth.Mount(s, r)
		user.Mount(s, r)
		member.Mount(s, r)
		bot.Mount(s, r)
		meta.Mount(s, r)
	})
}
