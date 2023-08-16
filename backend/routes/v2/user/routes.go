package user

import (
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{
		Server: srv,
	}

	r.Route("/users", func(r chi.Router) {
		r.With(server.MustAuth).Group(func(r chi.Router) {
			r.Get("/@me/settings", server.WrapHandler(s.GetSettings))
			r.Patch("/@me/settings", server.WrapHandler(s.PatchSettings))
		})
	})
}
