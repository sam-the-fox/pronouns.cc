package user

import (
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type GetUserResponse struct {
	ID          xid.ID            `json:"id"`
	Username    string            `json:"name"`
	DisplayName *string           `json:"display_name"`
	Bio         *string           `json:"bio"`
	AvatarURLs  []string          `json:"avatar_urls"`
	Links       []string          `json:"links"`
	Names       []db.FieldEntry   `json:"names"`
	Pronouns    []db.PronounEntry `json:"pronouns"`
	Members     []PartialMember   `json:"members"`
	Fields      []db.Field        `json:"fields"`
}

type GetMeResponse struct {
	GetUserResponse

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`
}

type PartialMember struct {
	ID          xid.ID   `json:"id"`
	Name        string   `json:"name"`
	DisplayName *string  `json:"display_name"`
	AvatarURLs  []string `json:"avatar_urls"`
}

func dbUserToResponse(u db.User, fields []db.Field, members []db.Member) GetUserResponse {
	resp := GetUserResponse{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURLs:  u.AvatarURLs,
		Links:       u.Links,
		Names:       u.Names,
		Pronouns:    u.Pronouns,
		Fields:      fields,
	}

	resp.Members = make([]PartialMember, len(members))
	for i := range members {
		resp.Members[i] = PartialMember{
			ID:          members[i].ID,
			Name:        members[i].Name,
			DisplayName: members[i].DisplayName,
			AvatarURLs:  members[i].AvatarURLs,
		}
	}

	return resp
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userRef := chi.URLParamFromCtx(ctx, "userRef")

	if id, err := xid.FromString(userRef); err == nil {
		u, err := s.DB.User(ctx, id)
		if err == nil {
			fields, err := s.DB.UserFields(ctx, u.ID)
			if err != nil {
				log.Errorf("Error getting user fields: %v", err)
				return err
			}

			members, err := s.DB.UserMembers(ctx, u.ID)
			if err != nil {
				log.Errorf("Error getting user members: %v", err)
				return err
			}

			render.JSON(w, r, dbUserToResponse(u, fields, members))
			return nil
		} else if err != db.ErrUserNotFound {
			log.Errorf("Error getting user by ID: %v", err)
			return err
		}
		// otherwise, we fall back to checking usernames
	}

	u, err := s.DB.Username(ctx, userRef)
	if err == db.ErrUserNotFound {
		return server.APIError{
			Code: server.ErrUserNotFound,
		}

	} else if err != nil {
		log.Errorf("Error getting user by username: %v", err)
		return err
	}

	if u.DeletedAt != nil {
		return server.APIError{Code: server.ErrUserNotFound}
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user fields: %v", err)
		return err
	}

	members, err := s.DB.UserMembers(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user members: %v", err)
		return err
	}

	render.JSON(w, r, dbUserToResponse(u, fields, members))
	return nil
}

func (s *Server) getMeUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		log.Errorf("Error getting user: %v", err)
		return err
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user fields: %v", err)
		return err
	}

	members, err := s.DB.UserMembers(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user members: %v", err)
		return err
	}

	render.JSON(w, r, GetMeResponse{
		GetUserResponse: dbUserToResponse(u, fields, members),
		Discord:         u.Discord,
		DiscordUsername: u.DiscordUsername,
	})
	return nil
}
