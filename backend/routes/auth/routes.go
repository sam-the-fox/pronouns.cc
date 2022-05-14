package auth

import (
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{srv}

	r.Route("/auth", func(r chi.Router) {
		// generate csrf token, returns all supported OAuth provider URLs
		r.Post("/urls", server.WrapHandler(s.oauthURLs))

		r.Route("/discord", func(r chi.Router) {
			// takes code + state, validates it, returns token OR discord signup ticket
			r.Post("/callback", server.WrapHandler(s.discordCallback))
			// takes discord signup ticket to register account
			r.Post("/signup", nil)
		})
	})
}

type oauthURLsRequest struct {
	CallbackDomain string `json:"callback_domain"`
}

type oauthURLsResponse struct {
	Discord string `json:"discord"`
}

func (s *Server) oauthURLs(w http.ResponseWriter, r *http.Request) error {
	req, err := Decode[oauthURLsRequest](r)
	if err != nil {
		log.Error(err)

		return server.APIError{Code: server.ErrBadRequest}
	}

	// generate CSRF state
	state, err := s.setCSRFState(r.Context())
	if err != nil {
		return errors.Wrap(err, "setting CSRF state")
	}

	// copy Discord config and set redirect url
	discordCfg := discordOAuthConfig
	discordCfg.RedirectURL = req.CallbackDomain + "/login/discord"

	render.JSON(w, r, oauthURLsResponse{
		Discord: discordCfg.AuthCodeURL(state),
	})
	return nil
}
