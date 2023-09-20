package member

import (
	"context"
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type GetMemberResponse struct {
	ID          xid.ID          `json:"id"`
	SnowflakeID common.MemberID `json:"id_new"`
	SID         string          `json:"sid"`
	Name        string          `json:"name"`
	DisplayName *string         `json:"display_name"`
	Bio         *string         `json:"bio"`
	Avatar      *string         `json:"avatar"`
	Links       []string        `json:"links"`

	Names    []db.FieldEntry   `json:"names"`
	Pronouns []db.PronounEntry `json:"pronouns"`
	Fields   []db.Field        `json:"fields"`
	Flags    []db.MemberFlag   `json:"flags"`

	User PartialUser `json:"user"`

	Unlisted *bool `json:"unlisted,omitempty"`
}

func dbMemberToMember(u db.User, m db.Member, fields []db.Field, flags []db.MemberFlag, isOwnMember bool) GetMemberResponse {
	r := GetMemberResponse{
		ID:          m.ID,
		SnowflakeID: m.SnowflakeID,
		SID:         m.SID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Bio:         m.Bio,
		Avatar:      m.Avatar,
		Links:       db.NotNull(m.Links),

		Names:    db.NotNull(m.Names),
		Pronouns: db.NotNull(m.Pronouns),
		Fields:   db.NotNull(fields),
		Flags:    flags,

		User: PartialUser{
			ID:                u.ID,
			SnowflakeID:       u.SnowflakeID,
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
	SnowflakeID       common.UserID        `json:"id_new"`
	Username          string               `json:"name"`
	DisplayName       *string              `json:"display_name"`
	Avatar            *string              `json:"avatar"`
	CustomPreferences db.CustomPreferences `json:"custom_preferences"`
}

func (s *Server) getMember(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()

	var m db.Member
	if id, err := xid.FromString(chi.URLParam(r, "memberRef")); err == nil {
		m, err = s.DB.Member(ctx, id)
		if err != nil {
			log.Errorf("getting member by xid: %v", err)
		}
	}
	// xid was not valid
	if !m.SnowflakeID.IsValid() {
		id, err := common.ParseSnowflake(chi.URLParam(r, "memberRef"))
		if err != nil {
			return server.APIError{
				Code: server.ErrMemberNotFound,
			}
		}

		m, err = s.DB.MemberBySnowflake(ctx, common.MemberID(id))
		if err != nil {
			return server.APIError{
				Code: server.ErrMemberNotFound,
			}
		}
	}

	u, err := s.DB.User(ctx, m.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
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
		return errors.Wrap(err, "getting member fields")
	}

	flags, err := s.DB.MemberFlags(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "getting member flags")
	}

	render.JSON(w, r, dbMemberToMember(u, m, fields, flags, isOwnMember))
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
		return errors.Wrap(err, "getting member fields")
	}

	flags, err := s.DB.MemberFlags(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "getting member flags")
	}

	render.JSON(w, r, dbMemberToMember(u, m, fields, flags, isOwnMember))
	return nil
}

func (s *Server) getMeMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting me user")
	}

	m, err := s.DB.UserMember(ctx, claims.UserID, chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{
			Code: server.ErrMemberNotFound,
		}
	}

	fields, err := s.DB.MemberFields(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "getting member fields")
	}

	flags, err := s.DB.MemberFlags(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "getting member flags")
	}

	render.JSON(w, r, dbMemberToMember(u, m, fields, flags, true))
	return nil
}

func (s *Server) parseUser(ctx context.Context, userRef string) (u db.User, err error) {
	// check xid first
	if id, err := xid.FromString(userRef); err == nil {
		u, err := s.DB.User(ctx, id)
		if err == nil {
			return u, nil
		}
	}

	// if not an xid, check by snowflake
	if id, err := common.ParseSnowflake(userRef); err == nil {
		u, err := s.DB.UserBySnowflake(ctx, common.UserID(id))
		if err == nil {
			return u, nil
		}
	}

	// else, use username
	return s.DB.Username(ctx, userRef)
}
