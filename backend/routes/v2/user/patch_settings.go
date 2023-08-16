package user

import (
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/aarondl/opt/omitnull"
	"github.com/go-chi/render"
)

type PatchSettingsRequest struct {
	ReadChangelog      omitnull.Val[string] `json:"read_changelog"`
	ReadSettingsNotice omitnull.Val[string] `json:"read_settings_notice"`
}

func (s *Server) PatchSettings(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)
	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	var req PatchSettingsRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if !req.ReadChangelog.IsUnset() {
		u.Settings.ReadChangelog = req.ReadChangelog.GetOrZero()
	}
	if !req.ReadSettingsNotice.IsUnset() {
		u.Settings.ReadSettingsNotice = req.ReadSettingsNotice.GetOrZero()
	}

	err = s.DB.UpdateUserSettings(ctx, u.ID, u.Settings)
	if err != nil {
		return errors.Wrap(err, "updating user settings")
	}

	render.JSON(w, r, u.Settings)
	return nil
}
