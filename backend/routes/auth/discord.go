package auth

import (
	"fmt"
	"net/http"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/render"
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
	CallbackDomain string `json:"callback_domain"`
	Code           string `json:"code"`
	State          string `json:"state"`
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

	cfg := discordOAuthConfig
	cfg.RedirectURL = decoded.CallbackDomain + "/login/discord"
	token, err := cfg.Exchange(r.Context(), decoded.Code)
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

		fmt.Println(token)

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
