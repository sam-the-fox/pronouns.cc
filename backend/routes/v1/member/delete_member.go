package member

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
)

func (s *Server) deleteMember(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "this token is read-only"}
	}

	var m db.Member
	if id, err := xid.FromString(chi.URLParam(r, "memberRef")); err == nil {
		m, err = s.DB.Member(ctx, id)
		if err != nil {
			if err == db.ErrMemberNotFound {
				return server.APIError{Code: server.ErrMemberNotFound}
			}

			return errors.Wrap(err, "getting member")
		}
	} else if id, err := common.ParseSnowflake(chi.URLParam(r, "memberRef")); err == nil {
		m, err = s.DB.MemberBySnowflake(ctx, common.MemberID(id))
		if err != nil {
			if err == db.ErrMemberNotFound {
				return server.APIError{Code: server.ErrMemberNotFound}
			}

			return errors.Wrap(err, "getting member")
		}
	} else {
		return server.APIError{Code: server.ErrMemberNotFound}
	}

	if m.UserID != claims.UserID {
		return server.APIError{Code: server.ErrNotOwnMember}
	}

	err = s.DB.DeleteMember(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "deleting member")
	}

	if m.Avatar != nil {
		err = s.DB.DeleteMemberAvatar(ctx, m.ID, *m.Avatar)
		if err != nil {
			return errors.Wrap(err, "deleting member avatar")
		}
	}

	// update last active time
	err = s.DB.UpdateActiveTime(ctx, s.DB, claims.UserID)
	if err != nil {
		log.Errorf("updating last active time for user %v: %v", claims.UserID, err)
		return errors.Wrap(err, "updating last active time")
	}

	render.NoContent(w, r)
	return nil
}
