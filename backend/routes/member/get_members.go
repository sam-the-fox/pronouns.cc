package member

import (
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (s *Server) getUserMembers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	u, err := s.parseUser(ctx, chi.URLParam(r, "userRef"))
	if err != nil {
		return server.APIError{
			Code: server.ErrUserNotFound,
		}
	}

	ms, err := s.DB.UserMembers(ctx, u.ID)
	if err != nil {
		return err
	}

	render.JSON(w, r, ms)
	return nil
}

func (s *Server) getMeMembers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	ms, err := s.DB.UserMembers(ctx, claims.UserID)
	if err != nil {
		return err
	}

	render.JSON(w, r, ms)
	return nil
}
