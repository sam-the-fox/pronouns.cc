package auth

import (
	"net/http"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type Server struct {
	*server.Server

	RequireInvite bool
	BaseURL       string

	hcaptchaSitekey string
	hcaptchaSecret  string
}

type userResponse struct {
	ID          xid.ID            `json:"id"`
	SnowflakeID common.UserID     `json:"id_new"`
	Username    string            `json:"name"`
	DisplayName *string           `json:"display_name"`
	Bio         *string           `json:"bio"`
	Avatar      *string           `json:"avatar"`
	Links       []string          `json:"links"`
	Names       []db.FieldEntry   `json:"names"`
	Pronouns    []db.PronounEntry `json:"pronouns"`
	Fields      []db.Field        `json:"fields"`

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`

	Tumblr         *string `json:"tumblr"`
	TumblrUsername *string `json:"tumblr_username"`

	Google         *string `json:"google"`
	GoogleUsername *string `json:"google_username"`

	Fediverse         *string `json:"fediverse"`
	FediverseUsername *string `json:"fediverse_username"`
	FediverseInstance *string `json:"fediverse_instance"`
}

func dbUserToUserResponse(u db.User, fields []db.Field) *userResponse {
	return &userResponse{
		ID:                u.ID,
		SnowflakeID:       u.SnowflakeID,
		Username:          u.Username,
		DisplayName:       u.DisplayName,
		Bio:               u.Bio,
		Avatar:            u.Avatar,
		Links:             db.NotNull(u.Links),
		Names:             db.NotNull(u.Names),
		Pronouns:          db.NotNull(u.Pronouns),
		Fields:            db.NotNull(fields),
		Discord:           u.Discord,
		DiscordUsername:   u.DiscordUsername,
		Tumblr:            u.Tumblr,
		TumblrUsername:    u.TumblrUsername,
		Google:            u.Google,
		GoogleUsername:    u.GoogleUsername,
		Fediverse:         u.Fediverse,
		FediverseUsername: u.FediverseUsername,
		FediverseInstance: u.FediverseInstance,
	}
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{
		Server:          srv,
		RequireInvite:   os.Getenv("REQUIRE_INVITE") == "true",
		BaseURL:         os.Getenv("BASE_URL"),
		hcaptchaSitekey: os.Getenv("HCAPTCHA_SITEKEY"),
		hcaptchaSecret:  os.Getenv("HCAPTCHA_SECRET"),
	}

	r.Route("/auth", func(r chi.Router) {
		// check if username is taken
		r.Get("/username", server.WrapHandler(s.usernameTaken))

		// generate csrf token, returns all supported OAuth provider URLs
		r.Post("/urls", server.WrapHandler(s.oauthURLs))
		r.Get("/urls/fediverse", server.WrapHandler(s.getFediverseURL))

		r.Route("/discord", func(r chi.Router) {
			// takes code + state, validates it, returns token OR discord signup ticket
			r.Post("/callback", server.WrapHandler(s.discordCallback))
			// takes discord signup ticket to register account
			r.Post("/signup", server.WrapHandler(s.discordSignup))
			// takes discord signup ticket to link to existing account
			r.With(server.MustAuth).Post("/add-provider", server.WrapHandler(s.discordLink))
			// removes discord link from existing account
			r.With(server.MustAuth).Post("/remove-provider", server.WrapHandler(s.discordUnlink))
		})

		r.Route("/tumblr", func(r chi.Router) {
			r.Post("/callback", server.WrapHandler(s.tumblrCallback))
			r.Post("/signup", server.WrapHandler(s.tumblrSignup))
			r.With(server.MustAuth).Post("/add-provider", server.WrapHandler(s.tumblrLink))
			r.With(server.MustAuth).Post("/remove-provider", server.WrapHandler(s.tumblrUnlink))
		})

		r.Route("/google", func(r chi.Router) {
			r.Post("/callback", server.WrapHandler(s.googleCallback))
			r.Post("/signup", server.WrapHandler(s.googleSignup))
			r.With(server.MustAuth).Post("/add-provider", server.WrapHandler(s.googleLink))
			r.With(server.MustAuth).Post("/remove-provider", server.WrapHandler(s.googleUnlink))
		})

		r.Route("/mastodon", func(r chi.Router) {
			r.Post("/callback", server.WrapHandler(s.mastodonCallback))
			r.Post("/signup", server.WrapHandler(s.mastodonSignup))
			r.With(server.MustAuth).Post("/add-provider", server.WrapHandler(s.mastodonLink))
			r.With(server.MustAuth).Post("/remove-provider", server.WrapHandler(s.mastodonUnlink))
		})

		r.Route("/misskey", func(r chi.Router) {
			r.Post("/callback", server.WrapHandler(s.misskeyCallback))
			r.Post("/signup", server.WrapHandler(s.misskeySignup))
			r.With(server.MustAuth).Post("/add-provider", server.WrapHandler(s.misskeyLink))
		})

		// invite routes
		r.With(server.MustAuth).Get("/invites", server.WrapHandler(s.getInvites))
		r.With(server.MustAuth).Post("/invites", server.WrapHandler(s.createInvite))

		// tokens
		r.With(server.MustAuth).Get("/tokens", server.WrapHandler(s.getTokens))
		r.With(server.MustAuth).Post("/tokens", server.WrapHandler(s.createToken))
		r.With(server.MustAuth).Delete("/tokens", server.WrapHandler(s.deleteToken))

		// cancel user delete
		// uses a special token, so handled in the function itself
		r.Get("/cancel-delete", server.WrapHandler(s.cancelDelete))
		// force user delete
		// uses a special token (same as above)
		r.Get("/force-delete", server.WrapHandler(s.forceDelete))
	})
}

type oauthURLsRequest struct {
	CallbackDomain string `json:"callback_domain"`
}

type oauthURLsResponse struct {
	Discord string `json:"discord,omitempty"`
	Tumblr  string `json:"tumblr,omitempty"`
	Google  string `json:"google,omitempty"`
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
	var resp oauthURLsResponse

	if discordOAuthConfig.ClientID != "" {
		discordCfg := discordOAuthConfig
		discordCfg.RedirectURL = req.CallbackDomain + "/auth/login/discord"
		resp.Discord = discordCfg.AuthCodeURL(state) + "&prompt=none"
	}
	if tumblrOAuthConfig.ClientID != "" {
		tumblrCfg := tumblrOAuthConfig
		tumblrCfg.RedirectURL = req.CallbackDomain + "/auth/login/tumblr"
		resp.Tumblr = tumblrCfg.AuthCodeURL(state)
	}
	if googleOAuthConfig.ClientID != "" {
		googleCfg := googleOAuthConfig
		googleCfg.RedirectURL = req.CallbackDomain + "/auth/login/google"
		resp.Google = googleCfg.AuthCodeURL(state)
	}

	render.JSON(w, r, resp)
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
