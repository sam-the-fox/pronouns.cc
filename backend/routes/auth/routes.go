package auth

import (
	"net/http"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type Server struct {
	*server.Server

	RequireInvite bool
}

type userResponse struct {
	ID          xid.ID   `json:"id"`
	Username    string   `json:"name"`
	DisplayName *string  `json:"display_name"`
	Bio         *string  `json:"bio"`
	AvatarURLs  []string `json:"avatar_urls"`
	Links       []string `json:"links"`

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`
}

func dbUserToUserResponse(u db.User) *userResponse {
	return &userResponse{
		ID:              u.ID,
		Username:        u.Username,
		DisplayName:     u.DisplayName,
		Bio:             u.Bio,
		AvatarURLs:      u.AvatarURLs,
		Links:           u.Links,
		Discord:         u.Discord,
		DiscordUsername: u.DiscordUsername,
	}
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{
		Server:        srv,
		RequireInvite: os.Getenv("REQUIRE_INVITE") == "true",
	}

	r.Route("/auth", func(r chi.Router) {
		// check if username is taken
		r.Get("/username", server.WrapHandler(s.usernameTaken))

		// generate csrf token, returns all supported OAuth provider URLs
		r.Post("/urls", server.WrapHandler(s.oauthURLs))

		r.Route("/discord", func(r chi.Router) {
			// takes code + state, validates it, returns token OR discord signup ticket
			r.Post("/callback", server.WrapHandler(s.discordCallback))
			// takes discord signup ticket to register account
			r.Post("/signup", server.WrapHandler(s.discordSignup))
		})

		// invite routes
		r.With(server.MustAuth).Get("/invites", server.WrapHandler(s.getInvites))
		r.With(server.MustAuth).Post("/invites", server.WrapHandler(s.createInvite))

		// tokens
		r.With(server.MustAuth).Get("/tokens", server.WrapHandler(s.getTokens))
		r.With(server.MustAuth).Post("/tokens", server.WrapHandler(s.createToken))
		r.With(server.MustAuth).Delete("/tokens", server.WrapHandler(s.deleteToken))
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

func (s *Server) usernameTaken(w http.ResponseWriter, r *http.Request) error {
	type Response struct {
		Valid bool `json:"valid"`
		Taken bool `json:"taken"`
	}

	name := r.FormValue("username")
	if name == "" {
		render.JSON(w, r, Response{
			Valid: false,
		})
		return nil
	}

	valid, taken, err := s.DB.UsernameTaken(r.Context(), name)
	if err != nil {
		return err
	}

	render.JSON(w, r, Response{
		Valid: valid,
		Taken: taken,
	})
	return nil
}
