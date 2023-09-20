package backend

import (
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/auth"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/member"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/meta"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/mod"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/user"
	user2 "codeberg.org/pronounscc/pronouns.cc/backend/routes/v2/user"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"

	_ "embed"
)

// mountRoutes mounts all API routes on the server's router.
// they are all mounted under /v1/
func mountRoutes(s *server.Server) {
	s.Router.Route("/v1", func(r chi.Router) {
		auth.Mount(s, r)
		user.Mount(s, r)
		member.Mount(s, r)
		meta.Mount(s, r)
		mod.Mount(s, r)
	})

	s.Router.Route("/v2", func(r chi.Router) {
		user2.Mount(s, r)
	})
}
