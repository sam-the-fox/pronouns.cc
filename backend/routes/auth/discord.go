package auth

import (
	"net/http"
	"os"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/render"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/xid"
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

type discordOauthCallbackRequest struct {
	CallbackDomain string `json:"callback_domain"`
	Code           string `json:"code"`
	State          string `json:"state"`
}

type discordCallbackResponse struct {
	HasAccount bool `json:"has_account"` // if true, Token and User will be set. if false, Ticket and Discord will be set

	Token string        `json:"token,omitempty"`
	User  *userResponse `json:"user,omitempty"`

	Discord       string `json:"discord,omitempty"` // username, for UI purposes
	Ticket        string `json:"ticket,omitempty"`
	RequireInvite bool   `json:"require_invite"` // require an invite for signing up

	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (s *Server) discordCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	decoded, err := Decode[discordOauthCallbackRequest](r)
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
	cfg.RedirectURL = decoded.CallbackDomain + "/auth/login/discord"
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
		if u.DeletedAt != nil && *u.SelfDelete {
			// store cancel delete token
			token := undeleteToken()
			err = s.saveUndeleteToken(ctx, u.ID, token)
			if err != nil {
				log.Errorf("saving undelete token: %v", err)
				return err
			}

			render.JSON(w, r, discordCallbackResponse{
				HasAccount: true,
				Token:      token,
				User:       dbUserToUserResponse(u, []db.Field{}),
				IsDeleted:  true,
				DeletedAt:  u.DeletedAt,
			})
			return nil
		}

		err = u.UpdateFromDiscord(ctx, s.DB, du)
		if err != nil {
			log.Errorf("updating user %v with Discord info: %v", u.ID, err)
		}

		// TODO: implement user + token permissions
		tokenID := xid.New()
		token, err := s.Auth.CreateToken(u.ID, tokenID, false, false, true)
		if err != nil {
			return err
		}

		// save token to database
		_, err = s.DB.SaveToken(ctx, u.ID, tokenID)
		if err != nil {
			return errors.Wrap(err, "saving token to database")
		}

		fields, err := s.DB.UserFields(ctx, u.ID)
		if err != nil {
			return errors.Wrap(err, "querying fields")
		}

		render.JSON(w, r, discordCallbackResponse{
			HasAccount: true,
			Token:      token,
			User:       dbUserToUserResponse(u, fields),
		})

		return nil

	} else if err != db.ErrUserNotFound { // internal error
		return err
	}

	// no user found, so save a ticket + save their Discord info in Redis
	ticket := RandBase64(32)
	err = s.DB.SetJSON(ctx, "discord:"+ticket, du, "EX", "600")
	if err != nil {
		log.Errorf("setting Discord user for ticket %q: %v", ticket, err)
		return err
	}

	render.JSON(w, r, discordCallbackResponse{
		HasAccount:    false,
		Discord:       du.String(),
		Ticket:        ticket,
		RequireInvite: s.RequireInvite,
	})

	return nil
}

type linkRequest struct {
	Ticket string `json:"ticket"`
}

func (s *Server) discordLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	// only site tokens can be used for this endpoint
	if claims.APIToken || !claims.TokenWrite {
		return server.APIError{Code: server.ErrInvalidToken}
	}

	req, err := Decode[linkRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if u.Discord != nil {
		return server.APIError{Code: server.ErrAlreadyLinked}
	}

	du := new(discordgo.User)
	err = s.DB.GetJSON(ctx, "discord:"+req.Ticket, &du)
	if err != nil {
		log.Errorf("getting discord user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	err = u.UpdateFromDiscord(ctx, s.DB, du)
	if err != nil {
		return errors.Wrap(err, "updating user from discord")
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "getting user fields")
	}

	render.JSON(w, r, dbUserToUserResponse(u, fields))
	return nil
}

func (s *Server) discordUnlink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	// only site tokens can be used for this endpoint
	if claims.APIToken || !claims.TokenWrite {
		return server.APIError{Code: server.ErrInvalidToken}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if u.Discord == nil {
		return server.APIError{Code: server.ErrNotLinked}
	}

	err = u.UnlinkDiscord(ctx, s.DB)
	if err != nil {
		return errors.Wrap(err, "updating user in db")
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "getting user fields")
	}

	render.JSON(w, r, dbUserToUserResponse(u, fields))
	return nil
}

type discordSignupRequest struct {
	Ticket     string `json:"ticket"`
	Username   string `json:"username"`
	InviteCode string `json:"invite_code"`
}

type signupResponse struct {
	User  userResponse `json:"user"`
	Token string       `json:"token"`
}

func (s *Server) discordSignup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	req, err := Decode[discordSignupRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if s.RequireInvite && req.InviteCode == "" {
		return server.APIError{Code: server.ErrInviteRequired}
	}

	valid, taken, err := s.DB.UsernameTaken(ctx, req.Username)
	if err != nil {
		return err
	}
	if !valid {
		return server.APIError{Code: server.ErrInvalidUsername}
	}
	if taken {
		return server.APIError{Code: server.ErrUsernameTaken}
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback(ctx)

	du := new(discordgo.User)
	err = s.DB.GetJSON(ctx, "discord:"+req.Ticket, &du)
	if err != nil {
		log.Errorf("getting discord user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	u, err := s.DB.CreateUser(ctx, tx, req.Username)
	if err != nil {
		if errors.Cause(err) == db.ErrUsernameTaken {
			return server.APIError{Code: server.ErrUsernameTaken}
		}

		return errors.Wrap(err, "creating user")
	}

	err = u.UpdateFromDiscord(ctx, tx, du)
	if err != nil {
		return errors.Wrap(err, "updating user from discord")
	}

	if s.RequireInvite {
		valid, used, err := s.DB.InvalidateInvite(ctx, tx, req.InviteCode)
		if err != nil {
			return errors.Wrap(err, "checking and invalidating invite")
		}

		if !valid {
			return server.APIError{Code: server.ErrInviteRequired}
		}

		if used {
			return server.APIError{Code: server.ErrInviteAlreadyUsed}
		}
	}

	// delete sign up ticket
	err = s.DB.Redis.Do(ctx, radix.Cmd(nil, "DEL", "discord:"+req.Ticket))
	if err != nil {
		return errors.Wrap(err, "deleting signup ticket")
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	// create token
	// TODO: implement user + token permissions
	tokenID := xid.New()
	token, err := s.Auth.CreateToken(u.ID, tokenID, false, false, true)
	if err != nil {
		return errors.Wrap(err, "creating token")
	}

	// save token to database
	_, err = s.DB.SaveToken(ctx, u.ID, tokenID)
	if err != nil {
		return errors.Wrap(err, "saving token to database")
	}

	// return user
	render.JSON(w, r, signupResponse{
		User:  *dbUserToUserResponse(u, nil),
		Token: token,
	})
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	decoded := *new(T)

	return decoded, render.Decode(r, &decoded)
}
