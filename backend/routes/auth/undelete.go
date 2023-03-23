package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/mediocregopher/radix/v4"
	"github.com/rs/xid"
)

func (s *Server) cancelDelete(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	token := r.Header.Get("X-Delete-Token")
	if token == "" {
		return server.APIError{Code: server.ErrForbidden}
	}

	id, err := s.getUndeleteToken(ctx, token)
	if err != nil {
		log.Errorf("getting undelete token: %v", err)
		return server.APIError{Code: server.ErrNotFound} // assume invalid token
	}

	// only self deleted users can undelete themselves
	u, err := s.DB.User(ctx, id)
	if err != nil {
		log.Errorf("getting user: %v", err)
		return errors.Wrap(err, "getting user")
	}
	if !*u.SelfDelete {
		return server.APIError{Code: server.ErrForbidden}
	}

	err = s.DB.UndoDeleteUser(ctx, id)
	if err != nil {
		log.Errorf("executing undelete query: %v", err)
	}

	render.JSON(w, r, map[string]any{"success": true})
	return nil
}

func undeleteToken() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *Server) saveUndeleteToken(ctx context.Context, userID xid.ID, token string) error {
	err := s.DB.Redis.Do(ctx, radix.Cmd(nil, "SET", "undelete:"+token, userID.String(), "EX", "3600"))
	if err != nil {
		return errors.Wrap(err, "setting undelete key")
	}
	return nil
}

func (s *Server) getUndeleteToken(ctx context.Context, token string) (userID xid.ID, err error) {
	var idString string
	err = s.DB.Redis.Do(ctx, radix.Cmd(&idString, "GETDEL", "undelete:"+token))
	if err != nil {
		return userID, errors.Wrap(err, "getting undelete key")
	}

	userID, err = xid.FromString(idString)
	if err != nil {
		return userID, errors.Wrap(err, "parsing ID")
	}
	return userID, nil
}

func (s *Server) forceDelete(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	token := r.Header.Get("X-Delete-Token")
	if token == "" {
		return server.APIError{Code: server.ErrForbidden}
	}

	id, err := s.getUndeleteToken(ctx, token)
	if err != nil {
		log.Errorf("getting delete token: %v", err)
		return server.APIError{Code: server.ErrNotFound} // assume invalid token
	}

	err = s.DB.CleanUser(ctx, id)
	if err != nil {
		log.Errorf("cleaning user data: %v", err)
		return errors.Wrap(err, "cleaning user")
	}

	err = s.DB.ForceDeleteUser(ctx, id)
	if err != nil {
		log.Errorf("force deleting user: %v", err)
		return errors.Wrap(err, "deleting user")
	}

	render.JSON(w, r, map[string]any{"success": true})
	return nil
}
