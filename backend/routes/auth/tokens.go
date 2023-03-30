package auth

import (
	"net/http"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type getTokenResponse struct {
	TokenID xid.ID    `json:"id"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

func dbTokenToGetResponse(t db.Token) getTokenResponse {
	return getTokenResponse{
		TokenID: t.TokenID,
		Created: t.Created,
		Expires: t.Expires,
	}
}

func (s *Server) getTokens(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

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

	if !claims.TokenWrite || claims.APIToken {
		return server.APIError{Code: server.ErrInvalidToken}
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

func (s *Server) createToken(w http.ResponseWriter, r *http.Request) error {
	// unimplemented right now
	return server.APIError{Code: server.ErrForbidden}
}
