package mod

import (
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

const MaxReasonLength = 2000

type CreateReportRequest struct {
	Reason string `json:"reason"`
}

func (s *Server) createUserReport(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	var u db.User
	if id, err := xid.FromString(chi.URLParam(r, "id")); err == nil {
		u, err = s.DB.User(ctx, id)
		if err != nil {
			if err == db.ErrUserNotFound {
				return server.APIError{Code: server.ErrUserNotFound}
			}

			log.Errorf("getting user %v: %v", id, err)
			return errors.Wrap(err, "getting user")
		}
	} else {
		id, err := common.ParseSnowflake(chi.URLParam(r, "id"))
		if err != nil {
			return server.APIError{Code: server.ErrUserNotFound}
		}

		u, err = s.DB.UserBySnowflake(ctx, common.UserID(id))
		if err != nil {
			if err == db.ErrUserNotFound {
				return server.APIError{Code: server.ErrUserNotFound}
			}

			log.Errorf("getting user %v: %v", id, err)
			return errors.Wrap(err, "getting user")
		}
	}

	if u.DeletedAt != nil {
		return server.APIError{Code: server.ErrUserNotFound}
	}

	var req CreateReportRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if len(req.Reason) > MaxReasonLength {
		return server.APIError{Code: server.ErrBadRequest, Details: "Reason cannot exceed 2000 characters"}
	}

	_, err = s.DB.CreateReport(ctx, claims.UserID, u.ID, nil, req.Reason)
	if err != nil {
		log.Errorf("creating report for %v: %v", u.ID, err)
		return errors.Wrap(err, "creating report")
	}

	render.NoContent(w, r)
	return nil
}

func (s *Server) createMemberReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	var m db.Member
	if id, err := xid.FromString(chi.URLParam(r, "id")); err == nil {
		m, err = s.DB.Member(ctx, id)
		if err != nil {
			if err == db.ErrMemberNotFound {
				return server.APIError{Code: server.ErrMemberNotFound}
			}

			log.Errorf("getting user %v: %v", id, err)
			return errors.Wrap(err, "getting user")
		}
	} else {
		id, err := common.ParseSnowflake(chi.URLParam(r, "id"))
		if err != nil {
			return server.APIError{Code: server.ErrUserNotFound}
		}

		m, err = s.DB.MemberBySnowflake(ctx, common.MemberID(id))
		if err != nil {
			if err == db.ErrMemberNotFound {
				return server.APIError{Code: server.ErrMemberNotFound}
			}

			log.Errorf("getting user %v: %v", id, err)
			return errors.Wrap(err, "getting user")
		}
	}

	u, err := s.DB.User(ctx, m.UserID)
	if err != nil {
		log.Errorf("getting user %v: %v", m.UserID, err)
		return errors.Wrap(err, "getting user")
	}

	if u.DeletedAt != nil {
		return server.APIError{Code: server.ErrMemberNotFound}
	}

	var req CreateReportRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if len(req.Reason) > MaxReasonLength {
		return server.APIError{Code: server.ErrBadRequest, Details: "Reason cannot exceed 2000 characters"}
	}

	_, err = s.DB.CreateReport(ctx, claims.UserID, u.ID, &m.ID, req.Reason)
	if err != nil {
		log.Errorf("creating report for %v: %v", m.ID, err)
		return errors.Wrap(err, "creating report")
	}

	render.NoContent(w, r)
	return nil
}
