package auth

import (
	"net/http"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type getTokenResponse struct {
	TokenID  xid.ID    `json:"id"`
	APIOnly  bool      `json:"api_only"`
	ReadOnly bool      `json:"read_only"`
	Created  time.Time `json:"created"`
	Expires  time.Time `json:"expires"`
}

func dbTokenToGetResponse(t db.Token) getTokenResponse {
	return getTokenResponse{
		TokenID:  t.TokenID,
		APIOnly:  t.APIOnly,
		ReadOnly: t.ReadOnly,
		Created:  t.Created,
		Expires:  t.Expires,
	}
}

func (s *Server) getTokens(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	tokens, err := s.DB.Tokens(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting tokens")
	}

	resps := make([]getTokenResponse, len(tokens))
	for i := range tokens {
		resps[i] = dbTokenToGetResponse(tokens[i])
	}

	render.JSON(w, r, resps)
	return nil
}

func (s *Server) deleteToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback(ctx)

	err = s.DB.InvalidateAllTokens(ctx, tx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "invalidating tokens")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	render.NoContent(w, r)
	return nil
}

type createTokenResponse struct {
	Token    string    `json:"token"`
	TokenID  xid.ID    `json:"id"`
	APIOnly  bool      `json:"api_only"`
	ReadOnly bool      `json:"read_only"`
	Created  time.Time `json:"created"`
	Expires  time.Time `json:"expires"`
}

func (s *Server) createToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting me user")
	}

	readOnly := r.FormValue("read_only") == "true"
	tokenID := xid.New()
	tokenStr, err := s.Auth.CreateToken(claims.UserID, tokenID, u.IsAdmin, true, !readOnly)
	if err != nil {
		return errors.Wrap(err, "creating token")
	}

	t, err := s.DB.SaveToken(ctx, claims.UserID, tokenID, true, readOnly)
	if err != nil {
		return errors.Wrap(err, "saving token")
	}

	render.JSON(w, r, createTokenResponse{
		Token:    tokenStr,
		TokenID:  t.TokenID,
		APIOnly:  t.APIOnly,
		ReadOnly: t.ReadOnly,
		Created:  t.Created,
		Expires:  t.Expires,
	})
	return nil
}
