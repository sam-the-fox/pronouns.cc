package user

import (
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "creating transaction")
	}
	defer tx.Rollback(ctx)

	err = s.DB.DeleteUser(ctx, tx, claims.UserID, true, "")
	if err != nil {
		return errors.Wrap(err, "setting user as deleted")
	}

	err = s.DB.InvalidateAllTokens(ctx, tx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "invalidating tokens")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	render.NoContent(w, r)
	return nil
}
