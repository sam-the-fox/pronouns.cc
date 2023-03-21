package mod

import (
	"net/http"
	"strconv"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

func (s *Server) getReports(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	showClosed := r.FormValue("closed") == "true"
	var before int
	if s := r.FormValue("before"); s != "" {
		before, err = strconv.Atoi(s)
		if err != nil {
			return server.APIError{Code: server.ErrBadRequest, Details: "\"before\": invalid ID"}
		}
	}

	reports, err := s.DB.Reports(ctx, showClosed, before)
	if err != nil {
		log.Errorf("getting reports: %v", err)
		return errors.Wrap(err, "getting reports from database")
	}

	render.JSON(w, r, reports)
	return nil
}

func (s *Server) getReportsByUser(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	var before int
	if s := r.FormValue("before"); s != "" {
		before, err = strconv.Atoi(s)
		if err != nil {
			return server.APIError{Code: server.ErrBadRequest, Details: "\"before\": invalid ID"}
		}
	}

	userID, err := xid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest, Details: "Invalid user ID"}
	}

	reports, err := s.DB.ReportsByUser(ctx, userID, before)
	if err != nil {
		log.Errorf("getting reports: %v", err)
		return errors.Wrap(err, "getting reports from database")
	}

	render.JSON(w, r, reports)
	return nil
}

func (s *Server) getReportsByReporter(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	var before int
	if s := r.FormValue("before"); s != "" {
		before, err = strconv.Atoi(s)
		if err != nil {
			return server.APIError{Code: server.ErrBadRequest, Details: "\"before\": invalid ID"}
		}
	}

	userID, err := xid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest, Details: "Invalid user ID"}
	}

	reports, err := s.DB.ReportsByReporter(ctx, userID, before)
	if err != nil {
		log.Errorf("getting reports: %v", err)
		return errors.Wrap(err, "getting reports from database")
	}

	render.JSON(w, r, reports)
	return nil
}
