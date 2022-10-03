package member

import (
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type memberListResponse struct {
	ID         xid.ID   `json:"id"`
	Name       string   `json:"name"`
	Bio        *string  `json:"bio"`
	AvatarURLs []string `json:"avatar_urls"`
	Links      []string `json:"links"`
}

func membersToMemberList(ms []db.Member) []memberListResponse {
	resps := make([]memberListResponse, len(ms))
	for i := range ms {
		resps[i] = memberListResponse{
			ID:         ms[i].ID,
			Name:       ms[i].Name,
			Bio:        ms[i].Bio,
			AvatarURLs: ms[i].AvatarURLs,
			Links:      ms[i].Links,
		}
	}

	return resps
}

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

	render.JSON(w, r, membersToMemberList(ms))
	return nil
}

func (s *Server) getMeMembers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	ms, err := s.DB.UserMembers(ctx, claims.UserID)
	if err != nil {
		return err
	}

	render.JSON(w, r, membersToMemberList(ms))
	return nil
}
