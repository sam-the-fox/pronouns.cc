package member

import (
	"github.com/go-chi/chi/v5"

	"codeberg.org/u1f320/pronouns.cc/backend/server"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{Server: srv}

	// member list
	r.Get("/users/{userRef}/members", server.WrapHandler(s.getUserMembers))
	r.With(server.MustAuth).Get("/users/@me/members", server.WrapHandler(s.getMeMembers))

	// user-scoped member lookup (including custom urls)
	r.Get("/users/{userRef}/members/{memberRef}", server.WrapHandler(s.getUserMember))
	r.With(server.MustAuth).Get("/users/@me/members/{memberRef}", server.WrapHandler(s.getMeMember))

	r.Route("/members", func(r chi.Router) {
		// any member by ID
		r.Get("/{memberRef}", server.WrapHandler(s.getMember))

		// create, edit, and delete members
		r.With(server.MustAuth).Post("/", server.WrapHandler(s.createMember))
		r.With(server.MustAuth).Patch("/{memberRef}", server.WrapHandler(s.patchMember))
		r.With(server.MustAuth).Delete("/{memberRef}", server.WrapHandler(s.deleteMember))
	})
}
