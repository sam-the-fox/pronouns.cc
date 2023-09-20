package user

import (
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

func (s *Server) GetSettings(w http.ResponseWriter, r *http.Request) (err error) {
	claims, _ := server.ClaimsFromContext(r.Context())
	u, err := s.DB.User(r.Context(), claims.UserID)
	if err != nil {
		log.Errorf("getting user: %v", err)
		return errors.Wrap(err, "getting user")
	}

	render.JSON(w, r, u.Settings)
	return nil
}
