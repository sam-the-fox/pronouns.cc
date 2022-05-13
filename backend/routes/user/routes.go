package user

import (
	"github.com/go-chi/chi/v5"
	"gitlab.com/1f320/pronouns/backend/server"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{srv}

	r.Route("/users", func(r chi.Router) {
		r.With(server.MustAuth).Get("/@me", server.WrapHandler(s.getMeUser))

		r.Get("/{userRef}", server.WrapHandler(s.getUser))
	})
}
