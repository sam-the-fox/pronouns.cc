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
	ID        xid.ID   `json:"id"`
	Name      string   `json:"name"`
	Bio       *string  `json:"bio"`
	AvatarURL *string  `json:"avatar_url"`
	Links     []string `json:"links"`

	Names    []db.Name    `json:"names"`
	Pronouns []db.Pronoun `json:"pronouns"`
	Fields   []db.Field   `json:"fields"`

	User PartialUser `json:"user"`
}

func dbMemberToMember(u db.User, m db.Member, names []db.Name, pronouns []db.Pronoun, fields []db.Field) GetMemberResponse {
	return GetMemberResponse{
		ID:        m.ID,
		Name:      m.Name,
		Bio:       m.Bio,
		AvatarURL: m.AvatarURL,
		Links:     m.Links,

		Names:    names,
		Pronouns: pronouns,
		Fields:   fields,

		User: PartialUser{
			ID:          u.ID,
			Username:    u.Username,
			DisplayName: u.DisplayName,
			AvatarURLs:  u.AvatarURLs,
		},
	}
}

type PartialUser struct {
	ID          xid.ID   `json:"id"`
	Username    string   `json:"username"`
	DisplayName *string  `json:"display_name"`
	AvatarURLs  []string `json:"avatar_urls"`
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

	names, err := s.DB.MemberNames(ctx, m.ID)
	if err != nil {
		return err
	}

	pronouns, err := s.DB.MemberPronouns(ctx, m.ID)
	if err != nil {
		return err
	}

	fields, err := s.DB.MemberFields(ctx, m.ID)
	if err != nil {
		return err
	}

	render.JSON(w, r, dbMemberToMember(u, m, names, pronouns, fields))
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

	m, err := s.DB.UserMember(ctx, u.ID, chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{
			Code: server.ErrMemberNotFound,
		}
	}

	names, err := s.DB.MemberNames(ctx, m.ID)
	if err != nil {
		return err
	}

	pronouns, err := s.DB.MemberPronouns(ctx, m.ID)
	if err != nil {
		return err
	}

	fields, err := s.DB.MemberFields(ctx, m.ID)
	if err != nil {
		return err
	}

	render.JSON(w, r, dbMemberToMember(u, m, names, pronouns, fields))
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
