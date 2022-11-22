package member

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
)

func (s *Server) deleteMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	id, err := xid.FromString(chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{Code: server.ErrMemberNotFound}
	}

	m, err := s.DB.Member(ctx, id)
	if err != nil {
		if err == db.ErrMemberNotFound {
			return server.APIError{Code: server.ErrMemberNotFound}
		}

		return errors.Wrap(err, "getting member")
	}

	if m.UserID != claims.UserID {
		return server.APIError{Code: server.ErrNotOwnMember}
	}

	err = s.DB.DeleteMember(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "deleting member")
	}

	render.NoContent(w, r)
	return nil
}
