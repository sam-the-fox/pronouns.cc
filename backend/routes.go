package backend

import (
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/auth"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/bot"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/member"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/meta"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/mod"
	"codeberg.org/pronounscc/pronouns.cc/backend/routes/v1/user"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	_ "embed"
)

//go:embed openapi.html
var openapi string

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
		mod.Mount(s, r)
	})

	// API docs
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r, openapi)
	})
}
