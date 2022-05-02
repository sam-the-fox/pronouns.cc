package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/1f320/pronouns/backend/server"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{srv}
	_ = s

	r.Route("/auth/discord", func(r chi.Router) {
		r.Get("/authorize", nil) // generate csrf token, returns URL
		r.Get("/callback", nil)  // takes code + state, validates it, returns token OR discord signup ticket
		r.Get("/signup", nil)    // takes discord signup ticket to register account

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world!"))
		})
	})
}
