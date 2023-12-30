package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/xid"
)

type fediOauthCallbackRequest struct {
	Instance string `json:"instance"`
	Code     string `json:"code"`
	State    string `json:"state"`
}

type fediCallbackResponse struct {
	HasAccount bool `json:"has_account"` // if true, Token and User will be set. if false, Ticket and Fediverse will be set

	Token string        `json:"token,omitempty"`
	User  *userResponse `json:"user,omitempty"`

	Fediverse      string `json:"fediverse,omitempty"` // username, for UI purposes
	Ticket         string `json:"ticket,omitempty"`
	RequireInvite  bool   `json:"require_invite"` // require an invite for signing up
	RequireCaptcha bool   `json:"require_captcha"`

	IsDeleted    bool       `json:"is_deleted"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	SelfDelete   *bool      `json:"self_delete,omitempty"`
	DeleteReason *string    `json:"delete_reason,omitempty"`
}

type partialMastodonAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (s *Server) mastodonCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	decoded, err := Decode[fediOauthCallbackRequest](r)
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

	app, err := s.DB.FediverseApp(ctx, decoded.Instance)
	if err != nil {
		log.Errorf("getting app for instance %q: %v", decoded.Instance, err)

		if err == db.ErrNoInstanceApp {
			// can we get here?
			return server.APIError{Code: server.ErrNotFound}
		}
	}

	token, err := app.ClientConfig().Exchange(ctx, decoded.Code)
	if err != nil {
		log.Errorf("exchanging oauth code: %v", err)

		return server.APIError{Code: server.ErrInvalidOAuthCode}
	}

	// make me user request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+decoded.Instance+"/api/v1/accounts/verify_credentials", nil)
	if err != nil {
		return errors.Wrap(err, "creating verify_credentials request")
	}
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token.Type()+" "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending verify_credentials request")
	}
	defer resp.Body.Close()

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading verify_credentials response")
	}

	var mu partialMastodonAccount
	err = json.Unmarshal(jb, &mu)
	if err != nil {
		return errors.Wrap(err, "unmarshaling verify_credentials response")
	}

	u, err := s.DB.FediverseUser(ctx, mu.ID, app.ID)
	if err == nil {
		if u.DeletedAt != nil {
			// store cancel delete token
			token := undeleteToken()
			err = s.saveUndeleteToken(ctx, u.ID, token)
			if err != nil {
				log.Errorf("saving undelete token: %v", err)
				return errors.Wrap(err, "saving undelete token")
			}

			render.JSON(w, r, fediCallbackResponse{
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

		err = u.UpdateFromFedi(ctx, s.DB, mu.ID, mu.Username, app.ID)
		if err != nil {
			log.Errorf("updating user %v with mastoAPI info: %v", u.ID, err)
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

		render.JSON(w, r, fediCallbackResponse{
			HasAccount: true,
			Token:      token,
			User:       dbUserToUserResponse(u, fields),
		})

		return nil

	} else if err != db.ErrUserNotFound { // internal error
		return errors.Wrap(err, "getting user")
	}

	// no user found, so save a ticket + save their Mastodon info in Redis
	ticket := RandBase64(32)
	err = s.DB.SetJSON(ctx, "mastodon:"+ticket, mu, "EX", "600")
	if err != nil {
		log.Errorf("setting mastoAPI user for ticket %q: %v", ticket, err)
		return errors.Wrap(err, "setting user for ticket")
	}

	render.JSON(w, r, fediCallbackResponse{
		HasAccount:     false,
		Fediverse:      mu.Username,
		Ticket:         ticket,
		RequireInvite:  s.RequireInvite,
		RequireCaptcha: s.hcaptchaSecret != "",
	})

	return nil
}

type fediLinkRequest struct {
	Instance string `json:"instance"`
	Ticket   string `json:"ticket"`
}

func (s *Server) mastodonLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	// only site tokens can be used for this endpoint
	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	req, err := Decode[fediLinkRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	app, err := s.DB.FediverseApp(ctx, req.Instance)
	if err != nil {
		return errors.Wrap(err, "getting instance application")
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if u.Fediverse != nil {
		return server.APIError{Code: server.ErrAlreadyLinked}
	}

	mu := new(partialMastodonAccount)
	err = s.DB.GetJSON(ctx, "mastodon:"+req.Ticket, &mu)
	if err != nil {
		log.Errorf("getting mastoAPI user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	if mu.ID == "" {
		log.Errorf("linking user with id %v: user ID was empty", claims.UserID)
		return server.APIError{Code: server.ErrInternalServerError, Details: "Mastodon user ID is empty"}
	}

	err = u.UpdateFromFedi(ctx, s.DB, mu.ID, mu.Username, app.ID)
	if err != nil {
		return errors.Wrap(err, "updating user from mastoAPI")
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "getting user fields")
	}

	render.JSON(w, r, dbUserToUserResponse(u, fields))
	return nil
}

func (s *Server) mastodonUnlink(w http.ResponseWriter, r *http.Request) error {
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

	if u.Fediverse == nil {
		return server.APIError{Code: server.ErrNotLinked}
	}

	// cannot unlink last auth provider
	if u.NumProviders() <= 1 {
		return server.APIError{Code: server.ErrLastProvider}
	}

	err = u.UnlinkFedi(ctx, s.DB)
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

type fediSignupRequest struct {
	Instance        string `json:"instance"`
	Ticket          string `json:"ticket"`
	Username        string `json:"username"`
	InviteCode      string `json:"invite_code"`
	CaptchaResponse string `json:"captcha_response"`
}

func (s *Server) mastodonSignup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	req, err := Decode[fediSignupRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if s.RequireInvite && req.InviteCode == "" {
		return server.APIError{Code: server.ErrInviteRequired}
	}

	app, err := s.DB.FediverseApp(ctx, req.Instance)
	if err != nil {
		return errors.Wrap(err, "getting instance application")
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
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("rolling back transaction:", err)
		}
	}()

	mu := new(partialMastodonAccount)
	err = s.DB.GetJSON(ctx, "mastodon:"+req.Ticket, &mu)
	if err != nil {
		log.Errorf("getting mastoAPI user for ticket: %v", err)

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

	if mu.ID == "" {
		log.Errorf("creating user with name %q: user ID was empty", req.Username)
		return server.APIError{Code: server.ErrInternalServerError, Details: "Mastodon user ID is empty"}
	}

	err = u.UpdateFromFedi(ctx, tx, mu.ID, mu.Username, app.ID)
	if err != nil {
		return errors.Wrap(err, "updating user from mastoAPI")
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
	err = s.DB.Redis.Do(ctx, radix.Cmd(nil, "DEL", "mastodon:"+req.Ticket))
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
