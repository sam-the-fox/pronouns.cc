package mod

import (
	"net/http"
	"strconv"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type warning struct {
	db.Warning
	Read bool `json:"read"`
}

func dbWarningsToResponse(ws []db.Warning) []warning {
	out := make([]warning, len(ws))
	for i := range ws {
		out[i] = warning{ws[i], ws[i].ReadAt != nil}
	}
	return out
}

func (s *Server) getWarnings(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)
	showAll := r.FormValue("all") == "true"

	warnings, err := s.DB.Warnings(ctx, claims.UserID, !showAll)
	if err != nil {
		log.Errorf("getting warnings: %v", err)
		return errors.Wrap(err, "getting warnings from database")
	}

	render.JSON(w, r, dbWarningsToResponse(warnings))
	return nil
}

func (s *Server) ackWarning(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	ok, err := s.DB.AckWarning(ctx, claims.UserID, id)
	if err != nil {
		log.Errorf("acknowledging warning: %v", err)
		return errors.Wrap(err, "acknowledging warning")
	}
	if !ok {
		return server.APIError{Code: server.ErrNotFound}
	}

	render.NoContent(w, r)
	return nil
}
