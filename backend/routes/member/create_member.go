package member

import (
	"context"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/render"
)

type CreateMemberRequest struct {
	Name      string       `json:"name"`
	Bio       *string      `json:"bio"`
	AvatarURL *string      `json:"avatar_url"`
	Links     []string     `json:"links"`
	Names     []db.Name    `json:"names"`
	Pronouns  []db.Pronoun `json:"pronouns"`
	Fields    []db.Field   `json:"fields"`
}

func (s *Server) createMember(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()

	var cmr CreateMemberRequest
	err = render.Decode(r, &cmr)
	if err != nil {
		if _, ok := err.(server.APIError); ok {
			return err
		}

		return server.APIError{Code: server.ErrBadRequest}
	}

	ctx = context.WithValue(ctx, render.StatusCtxKey, 204)
	render.NoContent(w, r)

	return nil
}
