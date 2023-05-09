package user

import (
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	*server.Server

	ExporterPath string
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{
		Server:       srv,
		ExporterPath: "http://127.0.0.1:" + os.Getenv("EXPORTER_PORT"),
	}

	r.Route("/users", func(r chi.Router) {
		r.Get("/{userRef}", server.WrapHandler(s.getUser))

		r.With(server.MustAuth).Group(func(r chi.Router) {
			r.Get("/@me", server.WrapHandler(s.getMeUser))
			r.Patch("/@me", server.WrapHandler(s.patchUser))
			r.Delete("/@me", server.WrapHandler(s.deleteUser))

			r.Get("/@me/export/start", server.WrapHandler(s.startExport))
			r.Get("/@me/export", server.WrapHandler(s.getExport))

			r.Get("/@me/flags", server.WrapHandler(s.getUserFlags))
			r.Post("/@me/flags", server.WrapHandler(s.postUserFlag))
			r.Patch("/@me/flags", server.WrapHandler(s.patchUserFlag))
			r.Delete("/@me/flags", server.WrapHandler(s.deleteUserFlag))
		})
	})
}
