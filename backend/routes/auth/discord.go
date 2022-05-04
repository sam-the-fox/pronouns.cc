package auth

import (
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/render"
	"gitlab.com/1f320/pronouns/backend/db"
	"gitlab.com/1f320/pronouns/backend/log"
	"gitlab.com/1f320/pronouns/backend/server"
	"golang.org/x/oauth2"
)

var discordOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://discord.com/api/oauth2/authorize",
		TokenURL:  "https://discord.com/api/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
	Scopes: []string{"identify"},
}

type oauthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type discordCallbackResponse struct {
	HasAccount bool `json:"has_account"` // if true, Token and User will be set. if false, Ticket and Discord will be set

	Token string   `json:"token,omitempty"`
	User  *db.User `json:"user,omitempty"`

	Discord string `json:"discord,omitempty"` // username, for UI purposes
	Ticket  string `json:"ticket,omitempty"`
}

func (s *Server) discordCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	decoded, err := Decode[oauthCallbackRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// if the state can't be validated, return
	if valid, err := s.validateCSRFState(ctx, decoded.State); !valid {
		if err != nil {
			return err
		}

		return server.APIError{Code: server.ErrInvalidState}
	}

	token, err := discordOAuthConfig.Exchange(r.Context(), decoded.Code)
	if err != nil {
		log.Errorf("exchanging oauth code: %v", err)

		return server.APIError{Code: server.ErrInvalidOAuthCode}
	}

	dg, _ := discordgo.New(token.Type() + " " + token.AccessToken)
	du, err := dg.User("@me")
	if err != nil {
		return err
	}

	u, err := s.DB.DiscordUser(ctx, du.ID)
	if err == nil {
		err = u.UpdateFromDiscord(ctx, s.DB, du)
		if err != nil {
			log.Errorf("updating user %v with Discord info: %v", u.ID, err)
		}

		token, err := s.Auth.CreateToken(u.ID)
		if err != nil {
			return err
		}

		render.JSON(w, r, discordCallbackResponse{
			HasAccount: true,
			Token:      token,
			User:       &u,
		})
	} else if err != db.ErrUserNotFound { // internal error
		return err
	}

	// no user found, so save a ticket

	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	decoded := *new(T)

	return decoded, render.Decode(r, &decoded)
}
