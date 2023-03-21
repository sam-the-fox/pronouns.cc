package mod

import (
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

const MaxReasonLength = 2000

type CreateReportRequest struct {
	Reason string `json:"reason"`
}

func (s *Server) createUserReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	userID, err := xid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest, Details: "Invalid user ID"}
	}

	u, err := s.DB.User(ctx, userID)
	if err != nil {
		if err == db.ErrUserNotFound {
			return server.APIError{Code: server.ErrUserNotFound}
		}

		log.Errorf("getting user %v: %v", userID, err)
		return errors.Wrap(err, "getting user")
	}

	var req CreateReportRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if len(req.Reason) > MaxReasonLength {
		return server.APIError{Code: server.ErrBadRequest, Details: "Reason cannot exceed 2000 characters"}
	}

	report, err := s.DB.CreateReport(ctx, claims.UserID, u.ID, nil, req.Reason)
	if err != nil {
		log.Errorf("creating report for %v: %v", u.ID, err)
		return errors.Wrap(err, "creating report")
	}

	render.JSON(w, r, map[string]any{"created": true, "created_at": report.CreatedAt})
	return nil
}
