package auth

import (
	"net/http"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

type inviteResponse struct {
	Code    string    `json:"code"`
	Created time.Time `json:"created"`
	Used    bool      `json:"used"`
}

func dbInviteToResponse(i db.Invite) inviteResponse {
	return inviteResponse{
		Code:    i.Code,
		Created: i.Created,
		Used:    i.Used,
	}
}

func (s *Server) getInvites(w http.ResponseWriter, r *http.Request) error {
	if !s.RequireInvite {
		return server.APIError{Code: server.ErrInvitesDisabled}
	}

	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	is, err := s.DB.UserInvites(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user invites")
	}

	resps := make([]inviteResponse, len(is))
	for i := range is {
		resps[i] = dbInviteToResponse(is[i])
	}

	render.JSON(w, r, resps)
	return nil
}

func (s *Server) createInvite(w http.ResponseWriter, r *http.Request) error {
	if !s.RequireInvite {
		return server.APIError{Code: server.ErrInvitesDisabled}
	}

	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if claims.APIToken {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This endpoint cannot be used by API tokens"}
	}

	inv, err := s.DB.CreateInvite(ctx, claims.UserID)
	if err != nil {
		if err == db.ErrTooManyInvites {
			return server.APIError{Code: server.ErrInviteLimitReached}
		}

		return errors.Wrap(err, "creating invite")
	}

	render.JSON(w, r, dbInviteToResponse(inv))
	return nil
}
