package auth

import (
	"net/http"
	"os"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

var googleOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://accounts.google.com/o/oauth2/auth",
		TokenURL:  "https://oauth2.googleapis.com/token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
	Scopes: []string{"openid", "https://www.googleapis.com/auth/userinfo.email"},
}

type googleCallbackResponse struct {
	HasAccount bool `json:"has_account"` // if true, Token and User will be set. if false, Ticket and Google will be set

	Token string        `json:"token,omitempty"`
	User  *userResponse `json:"user,omitempty"`

	Google         string `json:"google,omitempty"` // username, for UI purposes
	Ticket         string `json:"ticket,omitempty"`
	RequireInvite  bool   `json:"require_invite"` // require an invite for signing up
	RequireCaptcha bool   `json:"require_captcha"`

	IsDeleted    bool       `json:"is_deleted"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	SelfDelete   *bool      `json:"self_delete,omitempty"`
	DeleteReason *string    `json:"delete_reason,omitempty"`
}

type partialGoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (s *Server) googleCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	decoded, err := Decode[oauthCallbackRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// if the state can't be validated, return
	if valid, err := s.validateCSRFState(ctx, decoded.State); !valid {
		if err != nil {
			return errors.Wrap(err, "validating state")
		}

		return server.APIError{Code: server.ErrInvalidState}
	}

	cfg := googleOAuthConfig
	cfg.RedirectURL = decoded.CallbackDomain + "/auth/login/google"
	token, err := cfg.Exchange(r.Context(), decoded.Code)
	if err != nil {
		log.Errorf("exchanging oauth code: %v", err)
		return server.APIError{Code: server.ErrInvalidOAuthCode}
	}
	rawToken := token.Extra("id_token")
	if rawToken == nil {
		log.Debug("id_token was nil")
		return server.APIError{Code: server.ErrInternalServerError}
	}

	idToken, ok := rawToken.(string)
	if !ok {
		log.Debug("id_token was not a string")
		return server.APIError{Code: server.ErrInternalServerError}
	}
	payload, err := idtoken.Validate(ctx, idToken, "")
	if err != nil {
		log.Errorf("getting id token payload: %v", err)
		return server.APIError{Code: server.ErrInternalServerError}
	}

	googleID, ok := payload.Claims["sub"].(string)
	if !ok {
		log.Debug("id_token.claims.sub was not a string")
		return server.APIError{Code: server.ErrInternalServerError}
	}
	googleUsername, ok := payload.Claims["email"].(string)
	if !ok {
		log.Debug("id_token.claims.email was not a string")
		return server.APIError{Code: server.ErrInternalServerError}
	}

	u, err := s.DB.GoogleUser(ctx, googleID)
	if err == nil {
		if u.DeletedAt != nil {
			// store cancel delete token
			token := undeleteToken()
			err = s.saveUndeleteToken(ctx, u.ID, token)
			if err != nil {
				log.Errorf("saving undelete token: %v", err)
				return errors.Wrap(err, "saving undelete token")
			}

			render.JSON(w, r, googleCallbackResponse{
				HasAccount:   true,
				Token:        token,
				User:         dbUserToUserResponse(u, []db.Field{}),
				IsDeleted:    true,
				DeletedAt:    u.DeletedAt,
				SelfDelete:   u.SelfDelete,
				DeleteReason: u.DeleteReason,
			})
			return nil
		}

		err = u.UpdateFromGoogle(ctx, s.DB, googleID, googleUsername)
		if err != nil {
			log.Errorf("updating user %v with Google info: %v", u.ID, err)
		}

		// TODO: implement user + token permissions
		tokenID := xid.New()
		token, err := s.Auth.CreateToken(u.ID, tokenID, u.IsAdmin, false, true)
		if err != nil {
			return errors.Wrap(err, "creating token")
		}

		// save token to database
		_, err = s.DB.SaveToken(ctx, u.ID, tokenID, false, false)
		if err != nil {
			return errors.Wrap(err, "saving token to database")
		}

		fields, err := s.DB.UserFields(ctx, u.ID)
		if err != nil {
			return errors.Wrap(err, "querying fields")
		}

		render.JSON(w, r, googleCallbackResponse{
			HasAccount: true,
			Token:      token,
			User:       dbUserToUserResponse(u, fields),
		})

		return nil

	} else if err != db.ErrUserNotFound { // internal error
		return errors.Wrap(err, "getting user")
	}

	// no user found, so save a ticket + save their Google info in Redis
	ticket := RandBase64(32)
	err = s.DB.SetJSON(ctx, "google:"+ticket, partialGoogleUser{ID: googleID, Email: googleUsername}, "EX", "600")
	if err != nil {
		log.Errorf("setting Google user for ticket %q: %v", ticket, err)
		return errors.Wrap(err, "setting user for ticket")
	}

	render.JSON(w, r, googleCallbackResponse{
		HasAccount:     false,
		Google:         googleUsername,
		Ticket:         ticket,
		RequireInvite:  s.RequireInvite,
		RequireCaptcha: s.hcaptchaSecret != "",
	})

	return nil
}

func (s *Server) googleLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	// only site tokens can be used for this endpoint
	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	req, err := Decode[linkRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if u.Google != nil {
		return server.APIError{Code: server.ErrAlreadyLinked}
	}

	gu := new(partialGoogleUser)
	err = s.DB.GetJSON(ctx, "google:"+req.Ticket, &gu)
	if err != nil {
		log.Errorf("getting google user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	if gu.ID == "" {
		log.Errorf("linking user with id %v: user ID was empty", claims.UserID)
		return server.APIError{Code: server.ErrInternalServerError, Details: "Google user ID is empty"}
	}

	err = u.UpdateFromGoogle(ctx, s.DB, gu.ID, gu.Email)
	if err != nil {
		return errors.Wrap(err, "updating user from google")
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "getting user fields")
	}

	render.JSON(w, r, dbUserToUserResponse(u, fields))
	return nil
}

func (s *Server) googleUnlink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	// only site tokens can be used for this endpoint
	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if u.Google == nil {
		return server.APIError{Code: server.ErrNotLinked}
	}

	// cannot unlink last auth provider
	if u.NumProviders() <= 1 {
		return server.APIError{Code: server.ErrLastProvider}
	}

	err = u.UnlinkGoogle(ctx, s.DB)
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

func (s *Server) googleSignup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	req, err := Decode[signupRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if s.RequireInvite && req.InviteCode == "" {
		return server.APIError{Code: server.ErrInviteRequired}
	}

	valid, taken, err := s.DB.UsernameTaken(ctx, req.Username)
	if err != nil {
		return errors.Wrap(err, "checking if username is taken")
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

	gu := new(partialGoogleUser)
	err = s.DB.GetJSON(ctx, "google:"+req.Ticket, &gu)
	if err != nil {
		log.Errorf("getting google user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	// check captcha
	if s.hcaptchaSecret != "" {
		ok, err := s.verifyCaptcha(ctx, req.CaptchaResponse)
		if err != nil {
			log.Errorf("verifying captcha: %v", err)
			return server.APIError{Code: server.ErrInternalServerError}
		}

		if !ok {
			return server.APIError{Code: server.ErrInvalidCaptcha}
		}
	}

	u, err := s.DB.CreateUser(ctx, tx, req.Username)
	if err != nil {
		if errors.Cause(err) == db.ErrUsernameTaken {
			return server.APIError{Code: server.ErrUsernameTaken}
		}

		return errors.Wrap(err, "creating user")
	}

	if gu.ID == "" {
		log.Errorf("creating user with name %q: user ID was empty", req.Username)
		return server.APIError{Code: server.ErrInternalServerError, Details: "Google user ID is empty"}
	}

	err = u.UpdateFromGoogle(ctx, tx, gu.ID, gu.Email)
	if err != nil {
		return errors.Wrap(err, "updating user from google")
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
	err = s.DB.Redis.Do(ctx, radix.Cmd(nil, "DEL", "google:"+req.Ticket))
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
	_, err = s.DB.SaveToken(ctx, u.ID, tokenID, false, false)
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
