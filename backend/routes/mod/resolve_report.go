package mod

import (
	"net/http"
	"strconv"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type resolveReportRequest struct {
	Warn   bool   `json:"warn"`
	Ban    bool   `json:"ban"`
	Delete bool   `json:"delete"`
	Reason string `json:"reason"`
}

func (s *Server) resolveReport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	var req resolveReportRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if req.Reason == "" {
		return server.APIError{Code: server.ErrBadRequest, Details: "Reason cannot be empty"}
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		log.Errorf("creating transaction: %v", err)
		return errors.Wrap(err, "creating transaction")
	}
	defer tx.Rollback(ctx)

	report, err := s.DB.Report(ctx, tx, id)
	if err != nil {
		if err == db.ErrReportNotFound {
			return server.APIError{Code: server.ErrNotFound}
		}
		log.Errorf("getting report: %v", err)
		return errors.Wrap(err, "getting report")
	}

	if report.ResolvedAt != nil {
		return server.APIError{Code: server.ErrReportAlreadyHandled}
	}

	err = s.DB.ResolveReport(ctx, tx, report.ID, claims.UserID, req.Reason)
	if err != nil {
		log.Errorf("resolving report: %v", err)
	}

	if req.Warn || req.Ban {
		_, err = s.DB.CreateWarning(ctx, tx, report.UserID, req.Reason)
		if err != nil {
			log.Errorf("creating warning: %v", err)
		}
	}

	if req.Ban {
		err = s.DB.DeleteUser(ctx, tx, report.UserID, false, req.Reason)
		if err != nil {
			log.Errorf("banning user: %v", err)
		}

		if req.Delete {
			err = s.DB.InvalidateAllTokens(ctx, tx, report.UserID)
			if err != nil {
				return errors.Wrap(err, "invalidating tokens")
			}

			err = s.DB.CleanUser(ctx, report.UserID)
			if err != nil {
				log.Errorf("cleaning user data: %v", err)
				return errors.Wrap(err, "cleaning user")
			}

			err = s.DB.DeleteUserMembers(ctx, tx, report.UserID)
			if err != nil {
				log.Errorf("deleting members: %v", err)
				return errors.Wrap(err, "deleting members")
			}

			err = s.DB.ResetUser(ctx, tx, report.UserID)
			if err != nil {
				log.Errorf("resetting user data: %v", err)
				return errors.Wrap(err, "resetting user")
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("committing transaction: %v", err)
		return errors.Wrap(err, "committing transaction")
	}

	render.NoContent(w, r)
	return nil
}
