package user

import (
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{srv}

	r.Route("/users", func(r chi.Router) {
		r.Get("/{userRef}", server.WrapHandler(s.getUser))

		r.With(server.MustAuth).Group(func(r chi.Router) {
			r.Get("/@me", server.WrapHandler(s.getMeUser))
			r.Patch("/@me/fields", server.WrapHandler(s.patchUserFields))
		})
	})
}
