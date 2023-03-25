package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/xid"
)

type partialMisskeyAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (s *Server) misskeyCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	decoded, err := Decode[fediOauthCallbackRequest](r)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	app, err := s.DB.FediverseApp(ctx, decoded.Instance)
	if err != nil {
		log.Errorf("getting app for instance %q: %v", decoded.Instance, err)

		if err == db.ErrNoInstanceApp {
			// can we get here?
			return server.APIError{Code: server.ErrNotFound}
		}
	}

	userkeyReq := struct {
		AppSecret string `json:"appSecret"`
		Token     string `json:"token"`
	}{AppSecret: app.ClientSecret, Token: decoded.Code}

	b, err := json.Marshal(userkeyReq)
	if err != nil {
		return errors.Wrap(err, "marshaling json")
	}

	// make me user request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+decoded.Instance+"/api/auth/session/userkey", bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "creating userkey request")
	}
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending i request")
	}
	defer resp.Body.Close()

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading i response")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		log.Errorf("POST userkey for instance %q (type %v): %v", app.Instance, app.InstanceType, string(jb))
		return errors.Wrap(err, "error on misskey's end")
	}

	var mu struct {
		User partialMisskeyAccount `json:"user"`
	}
	err = json.Unmarshal(jb, &mu)
	if err != nil {
		return errors.Wrap(err, "unmarshaling userkey response")
	}

	u, err := s.DB.FediverseUser(ctx, mu.User.ID, app.ID)
	if err == nil {
		if u.DeletedAt != nil {
			// store cancel delete token
			token := undeleteToken()
			err = s.saveUndeleteToken(ctx, u.ID, token)
			if err != nil {
				log.Errorf("saving undelete token: %v", err)
				return err
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

		err = u.UpdateFromFedi(ctx, s.DB, mu.User.ID, mu.User.Username, app.ID)
		if err != nil {
			log.Errorf("updating user %v with misskey info: %v", u.ID, err)
		}

		// TODO: implement user + token permissions
		tokenID := xid.New()
		token, err := s.Auth.CreateToken(u.ID, tokenID, u.IsAdmin, false, true)
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

		render.JSON(w, r, fediCallbackResponse{
			HasAccount: true,
			Token:      token,
			User:       dbUserToUserResponse(u, fields),
		})

		return nil

	} else if err != db.ErrUserNotFound { // internal error
		return err
	}

	// no user found, so save a ticket + save their Misskey info in Redis
	ticket := RandBase64(32)
	err = s.DB.SetJSON(ctx, "misskey:"+ticket, mu.User, "EX", "600")
	if err != nil {
		log.Errorf("setting misskey user for ticket %q: %v", ticket, err)
		return err
	}

	render.JSON(w, r, fediCallbackResponse{
		HasAccount:    false,
		Fediverse:     mu.User.Username,
		Ticket:        ticket,
		RequireInvite: s.RequireInvite,
	})

	return nil
}

func (s *Server) misskeyLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	// only site tokens can be used for this endpoint
	if claims.APIToken || !claims.TokenWrite {
		return server.APIError{Code: server.ErrInvalidToken}
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

	mu := new(partialMisskeyAccount)
	err = s.DB.GetJSON(ctx, "misskey:"+req.Ticket, &mu)
	if err != nil {
		log.Errorf("getting misskey user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	err = u.UpdateFromFedi(ctx, s.DB, mu.ID, mu.Username, app.ID)
	if err != nil {
		return errors.Wrap(err, "updating user from misskey")
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "getting user fields")
	}

	render.JSON(w, r, dbUserToUserResponse(u, fields))
	return nil
}

func (s *Server) misskeySignup(w http.ResponseWriter, r *http.Request) error {
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

	mu := new(partialMisskeyAccount)
	err = s.DB.GetJSON(ctx, "misskey:"+req.Ticket, &mu)
	if err != nil {
		log.Errorf("getting misskey user for ticket: %v", err)

		return server.APIError{Code: server.ErrInvalidTicket}
	}

	u, err := s.DB.CreateUser(ctx, tx, req.Username)
	if err != nil {
		if errors.Cause(err) == db.ErrUsernameTaken {
			return server.APIError{Code: server.ErrUsernameTaken}
		}

		return errors.Wrap(err, "creating user")
	}

	err = u.UpdateFromFedi(ctx, tx, mu.ID, mu.Username, app.ID)
	if err != nil {
		return errors.Wrap(err, "updating user from misskey")
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
	err = s.DB.Redis.Do(ctx, radix.Cmd(nil, "DEL", "misskey:"+req.Ticket))
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

func (s *Server) noAppMisskeyURL(ctx context.Context, w http.ResponseWriter, r *http.Request, softwareName, instance string) error {
	log.Debugf("creating application on misskey-compatible instance %q", instance)

	b, err := json.Marshal(misskeyAppRequest{
		Name:        "pronouns.cc (+" + s.BaseURL + ")",
		Description: "pronouns.cc on " + s.BaseURL,
		CallbackURL: s.BaseURL + "/auth/login/misskey/" + instance,
		Permission:  []string{"read:account"},
	})
	if err != nil {
		log.Errorf("marshaling app json: %v", err)
		return errors.Wrap(err, "marshaling json")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+instance+"/api/app/create", bytes.NewReader(b))
	if err != nil {
		log.Errorf("creating POST apps request for %q: %v", instance, err)
		return errors.Wrap(err, "creating POST apps request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("sending POST apps request for %q: %v", instance, err)
		return errors.Wrap(err, "sending POST apps request")
	}
	defer resp.Body.Close()

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("reading response for request: %v", err)
		return errors.Wrap(err, "reading response")
	}

	var ma misskeyApp
	err = json.Unmarshal(jb, &ma)
	if err != nil {
		return errors.Wrap(err, "unmarshaling misskey app")
	}

	app, err := s.DB.CreateFediverseApp(ctx, instance, softwareName, ma.ID, ma.Secret)
	if err != nil {
		log.Errorf("saving app for %q: %v", instance, err)
		return errors.Wrap(err, "creating app")
	}

	_, url, err := s.misskeyURL(ctx, app)
	if err != nil {
		log.Errorf("generating URL for misskey %q: %v", instance, err)
		return errors.Wrap(err, "generating URL")
	}

	render.JSON(w, r, FediResponse{
		URL: url,
	})
	return nil
}

type misskeyAppRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permission  []string `json:"permission"`
	CallbackURL string   `json:"callbackUrl"`
}

type misskeyApp struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

func (s *Server) misskeyURL(ctx context.Context, app db.FediverseApp) (token, url string, err error) {
	genSession := struct {
		AppSecret string `json:"appSecret"`
	}{AppSecret: app.ClientSecret}

	b, err := json.Marshal(genSession)
	if err != nil {
		return token, url, errors.Wrap(err, "marshaling json")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+app.Instance+"/api/auth/session/generate", bytes.NewReader(b))
	if err != nil {
		log.Errorf("creating POST session request for %q: %v", app.Instance, err)
		return token, url, errors.Wrap(err, "creating POST apps request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("sending POST session request for %q: %v", app.Instance, err)
		return token, url, errors.Wrap(err, "sending POST apps request")
	}
	defer resp.Body.Close()

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("reading response for request: %v", err)
		return token, url, errors.Wrap(err, "reading response")
	}

	var genSessionResp struct {
		Token string `json:"token"`
		URL   string `json:"url"`
	}

	err = json.Unmarshal(jb, &genSessionResp)
	if err != nil {
		return token, url, errors.Wrap(err, "unmarshaling misskey response")
	}

	return genSessionResp.Token, genSessionResp.URL, nil
}
