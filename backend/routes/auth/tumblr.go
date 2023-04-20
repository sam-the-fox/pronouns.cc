package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

var tumblrOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("TUMBLR_CLIENT_ID"),
	ClientSecret: os.Getenv("TUMBLR_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://www.tumblr.com/oauth2/authorize",
		TokenURL:  "https://api.tumblr.com/v2/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
	Scopes: []string{"basic"},
}

type partialTumblrResponse struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"msg"`
	} `json:"meta"`
	Response struct {
		User struct {
			Blogs []struct {
				Name    string `json:"name"`
				Primary bool   `json:"primary"`
				UUID    string `json:"uuid"`
			} `json:"blogs"`
		} `json:"user"`
	} `json:"response"`
}

type tumblrUserInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type tumblrCallbackResponse struct {
	HasAccount bool `json:"has_account"` // if true, Token and User will be set. if false, Ticket and Tumblr will be set

	Token string        `json:"token,omitempty"`
	User  *userResponse `json:"user,omitempty"`

	Tumblr        string `json:"tumblr,omitempty"` // username, for UI purposes
	Ticket        string `json:"ticket,omitempty"`
	RequireInvite bool   `json:"require_invite"` // require an invite for signing up

	IsDeleted    bool       `json:"is_deleted"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	SelfDelete   *bool      `json:"self_delete,omitempty"`
	DeleteReason *string    `json:"delete_reason,omitempty"`
}

func (s *Server) tumblrCallback(w http.ResponseWriter, r *http.Request) error {
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

	cfg := tumblrOAuthConfig
	cfg.RedirectURL = decoded.CallbackDomain + "/auth/login/tumblr"
	token, err := cfg.Exchange(r.Context(), decoded.Code)
	if err != nil {
		log.Errorf("exchanging oauth code: %v", err)

		return server.APIError{Code: server.ErrInvalidOAuthCode}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.tumblr.com/v2/user/info", nil)
	if err != nil {
		return errors.Wrap(err, "creating user/info request")
	}

	req.Header.Set("Content-Type", "application/json")
	token.SetAuthHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending user/info request")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return errors.New("response had status code < 200 or >= 400")
	}

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading user/info response")
	}

	var tr partialTumblrResponse
	err = json.Unmarshal(jb, &tr)
	if err != nil {
		return errors.Wrap(err, "unmarshaling user/info response")
	}

	var tumblrName, tumblrID string
	for _, blog := range tr.Response.User.Blogs {
		if blog.Primary {
			tumblrName = blog.Name
			tumblrID = blog.UUID
			break
		}
	}

	if tumblrID == "" {
		return server.APIError{Code: server.ErrInternalServerError, Details: "Your Tumblr account doesn't seem to have a primary blog"}
	}

	u, err := s.DB.TumblrUser(ctx, tumblrID)
	if err == nil {
		if u.DeletedAt != nil {
			// store cancel delete token
			token := undeleteToken()
			err = s.saveUndeleteToken(ctx, u.ID, token)
			if err != nil {
				log.Errorf("saving undelete token: %v", err)
				return err
			}

			render.JSON(w, r, tumblrCallbackResponse{
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

		err = u.UpdateFromTumblr(ctx, s.DB, tumblrID, tumblrName)
		if err != nil {
			log.Errorf("updating user %v with Tumblr info: %v", u.ID, err)
		}

		// TODO: implement user + token permissions
		tokenID := xid.New()
		token, err := s.Auth.CreateToken(u.ID, tokenID, u.IsAdmin, false, true)
		if err != nil {
			return err
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

		render.JSON(w, r, tumblrCallbackResponse{
			HasAccount: true,
			Token:      token,
			User:       dbUserToUserResponse(u, fields),
		})

		return nil

	} else if err != db.ErrUserNotFound { // internal error
		return err
	}

	// no user found, so save a ticket + save their Tumblr info in Redis
	ticket := RandBase64(32)
	err = s.DB.SetJSON(ctx, "tumblr:"+ticket, tumblrUserInfo{ID: tumblrID, Name: tumblrName}, "EX", "600")
	if err != nil {
		log.Errorf("setting Tumblr user for ticket %q: %v", ticket, err)
		return err
	}

	render.JSON(w, r, tumblrCallbackResponse{
		HasAccount:    false,
		Tumblr:        tumblrName,
		Ticket:        ticket,
		RequireInvite: s.RequireInvite,
	})

	return nil
}

func (s *Server) tumblrLink(w http.ResponseWriter, r *http.Request) error {
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

	if u.Tumblr != nil {
		return server.APIError{Code: server.ErrAlreadyLinked}
	}

	tui := new(tumblrUserInfo)
	err = s.DB.GetJSON(ctx, "tumblr:"+req.Ticket, &tui)
	if err != nil {
		log.Errorf("getting tumblr user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	if tui.ID == "" {
		log.Errorf("linking user with id %v: user ID was empty", claims.UserID)
		return server.APIError{Code: server.ErrInternalServerError, Details: "Tumblr user ID is empty"}
	}

	err = u.UpdateFromTumblr(ctx, s.DB, tui.ID, tui.Name)
	if err != nil {
		return errors.Wrap(err, "updating user from tumblr")
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "getting user fields")
	}

	render.JSON(w, r, dbUserToUserResponse(u, fields))
	return nil
}

func (s *Server) tumblrUnlink(w http.ResponseWriter, r *http.Request) error {
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

	if u.Tumblr == nil {
		return server.APIError{Code: server.ErrNotLinked}
	}

	// cannot unlink last auth provider
	if u.NumProviders() <= 1 {
		return server.APIError{Code: server.ErrLastProvider}
	}

	err = u.UnlinkTumblr(ctx, s.DB)
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

func (s *Server) tumblrSignup(w http.ResponseWriter, r *http.Request) error {
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

	tui := new(tumblrUserInfo)
	err = s.DB.GetJSON(ctx, "tumblr:"+req.Ticket, &tui)
	if err != nil {
		log.Errorf("getting tumblr user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	u, err := s.DB.CreateUser(ctx, tx, req.Username)
	if err != nil {
		if errors.Cause(err) == db.ErrUsernameTaken {
			return server.APIError{Code: server.ErrUsernameTaken}
		}

		return errors.Wrap(err, "creating user")
	}

	if tui.ID == "" {
		log.Errorf("creating user with name %q: user ID was empty", req.Username)
		return server.APIError{Code: server.ErrInternalServerError, Details: "Tumblr user ID is empty"}
	}

	err = u.UpdateFromTumblr(ctx, tx, tui.ID, tui.Name)
	if err != nil {
		return errors.Wrap(err, "updating user from tumblr")
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
	err = s.DB.Redis.Do(ctx, radix.Cmd(nil, "DEL", "tumblr:"+req.Ticket))
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
