package auth

import (
	"net/http"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
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

type deleteTokenResponse struct {
	TokenID     xid.ID    `json:"id"`
	Invalidated bool      `json:"invalidated"`
	Created     time.Time `json:"time"`
}

func (s *Server) deleteToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	tokenID, err := xid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	t, err := s.DB.InvalidateToken(ctx, claims.UserID, tokenID)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return server.APIError{Code: server.ErrNotFound}
		}

		return errors.Wrap(err, "invalidating token")
	}

	render.JSON(w, r, deleteTokenResponse{
		TokenID:     t.TokenID,
		Invalidated: t.Invalidated,
		Created:     t.Created,
	})
	return nil
}

func (s *Server) createToken(w http.ResponseWriter, r *http.Request) error {
	// unimplemented right now
	return server.APIError{Code: server.ErrForbidden}
}
