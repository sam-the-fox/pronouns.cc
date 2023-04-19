package member

import (
	"context"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type GetMemberResponse struct {
	ID          xid.ID   `json:"id"`
	Name        string   `json:"name"`
	DisplayName *string  `json:"display_name"`
	Bio         *string  `json:"bio"`
	Avatar      *string  `json:"avatar"`
	Links       []string `json:"links"`

	Names    []db.FieldEntry   `json:"names"`
	Pronouns []db.PronounEntry `json:"pronouns"`
	Fields   []db.Field        `json:"fields"`

	User PartialUser `json:"user"`

	Unlisted *bool `json:"unlisted,omitempty"`
}

func dbMemberToMember(u db.User, m db.Member, fields []db.Field, isOwnMember bool) GetMemberResponse {
	r := GetMemberResponse{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Bio:         m.Bio,
		Avatar:      m.Avatar,
		Links:       db.NotNull(m.Links),

		Names:    db.NotNull(m.Names),
		Pronouns: db.NotNull(m.Pronouns),
		Fields:   db.NotNull(fields),

		User: PartialUser{
			ID:                u.ID,
			Username:          u.Username,
			DisplayName:       u.DisplayName,
			Avatar:            u.Avatar,
			CustomPreferences: u.CustomPreferences,
		},
	}

	if isOwnMember {
		r.Unlisted = &m.Unlisted
	}

	return r
}

type PartialUser struct {
	ID                xid.ID               `json:"id"`
	Username          string               `json:"name"`
	DisplayName       *string              `json:"display_name"`
	Avatar            *string              `json:"avatar"`
	CustomPreferences db.CustomPreferences `json:"custom_preferences"`
}

func (s *Server) getMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id, err := xid.FromString(chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{
			Code: server.ErrMemberNotFound,
		}
	}

	m, err := s.DB.Member(ctx, id)
	if err != nil {
		return server.APIError{
			Code: server.ErrMemberNotFound,
		}
	}

	u, err := s.DB.User(ctx, m.UserID)
	if err != nil {
		return err
	}

	if u.DeletedAt != nil {
		return server.APIError{Code: server.ErrMemberNotFound}
	}

	isOwnMember := false
	if claims, ok := server.ClaimsFromContext(ctx); ok && claims.UserID == u.ID {
		isOwnMember = true
	}

	fields, err := s.DB.MemberFields(ctx, m.ID)
	if err != nil {
		return err
	}

	render.JSON(w, r, dbMemberToMember(u, m, fields, isOwnMember))
	return nil
}

func (s *Server) getUserMember(w http.ResponseWriter, r *http.Request) error {
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

	isOwnMember := false
	if claims, ok := server.ClaimsFromContext(ctx); ok && claims.UserID == u.ID {
		isOwnMember = true
	}

	m, err := s.DB.UserMember(ctx, u.ID, chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{
			Code: server.ErrMemberNotFound,
		}
	}

	fields, err := s.DB.MemberFields(ctx, m.ID)
	if err != nil {
		return err
	}

	render.JSON(w, r, dbMemberToMember(u, m, fields, isOwnMember))
	return nil
}

func (s *Server) parseUser(ctx context.Context, userRef string) (u db.User, err error) {
	if id, err := xid.FromString(userRef); err != nil {
		u, err := s.DB.User(ctx, id)
		if err == nil {
			return u, nil
		}
	}

	return s.DB.Username(ctx, userRef)
}
