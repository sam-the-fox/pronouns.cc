package user

import (
	"net/http"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/render"
)

func (s *Server) startExport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	hasExport, err := s.DB.HasRecentExport(ctx, claims.UserID)
	if err != nil {
		log.Errorf("checking if user has recent export: %v", err)
		return server.APIError{Code: server.ErrInternalServerError}
	}
	if hasExport {
		return server.APIError{Code: server.ErrRecentExport}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", s.ExporterPath+"/start/"+claims.UserID.String(), nil)
	if err != nil {
		log.Errorf("creating start export request: %v", err)
		return server.APIError{Code: server.ErrInternalServerError}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("executing start export request: %v", err)
		return server.APIError{Code: server.ErrInternalServerError}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		log.Errorf("got non-%v code: %v", http.StatusAccepted, resp.StatusCode)
		return server.APIError{
			Code: server.ErrInternalServerError,
		}
	}

	render.NoContent(w, r)
	return nil
}

type dataExportResponse struct {
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) getExport(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	de, err := s.DB.UserExport(ctx, claims.UserID)
	if err != nil {
		if err == db.ErrNoExport {
			return server.APIError{Code: server.ErrNotFound}
		}

		log.Errorf("getting export for user %v: %v", claims.UserID, err)
		return err
	}

	render.JSON(w, r, dataExportResponse{
		Path:      de.Path(),
		CreatedAt: de.CreatedAt,
	})
	return nil
}
