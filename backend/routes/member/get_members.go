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
	ID          xid.ID            `json:"id"`
	Name        string            `json:"name"`
	DisplayName *string           `json:"display_name"`
	Bio         *string           `json:"bio"`
	Avatar      *string           `json:"avatar"`
	Links       []string          `json:"links"`
	Names       []db.FieldEntry   `json:"names"`
	Pronouns    []db.PronounEntry `json:"pronouns"`
	Unlisted    bool              `json:"unlisted"`
}

func membersToMemberList(ms []db.Member, isSelf bool) []memberListResponse {
	resps := make([]memberListResponse, len(ms))
	for i := range ms {
		resps[i] = memberListResponse{
			ID:       ms[i].ID,
			Name:     ms[i].Name,
			Bio:      ms[i].Bio,
			Avatar:   ms[i].Avatar,
			Links:    db.NotNull(ms[i].Links),
			Names:    db.NotNull(ms[i].Names),
			Pronouns: db.NotNull(ms[i].Pronouns),
		}

		if isSelf {
			resps[i].Unlisted = ms[i].Unlisted
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

	if u.DeletedAt != nil {
		return server.APIError{Code: server.ErrUserNotFound}
	}

	isSelf := false
	if claims, ok := server.ClaimsFromContext(ctx); ok && claims.UserID == u.ID {
		isSelf = true
	}

	if u.ListPrivate && !isSelf {
		return server.APIError{Code: server.ErrMemberListPrivate}
	}

	ms, err := s.DB.UserMembers(ctx, u.ID, isSelf)
	if err != nil {
		return err
	}

	render.JSON(w, r, membersToMemberList(ms, isSelf))
	return nil
}

func (s *Server) getMeMembers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	ms, err := s.DB.UserMembers(ctx, claims.UserID, true)
	if err != nil {
		return err
	}

	render.JSON(w, r, membersToMemberList(ms, true))
	return nil
}
