package user

import (
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"github.com/go-chi/render"
)

func (s *Server) GetSettings(w http.ResponseWriter, r *http.Request) (err error) {
	claims, _ := server.ClaimsFromContext(r.Context())
	u, err := s.DB.User(r.Context(), claims.UserID)
	if err != nil {
		log.Errorf("getting user: %v", err)
		return err
	}

	render.JSON(w, r, u.Settings)
	return nil
}
