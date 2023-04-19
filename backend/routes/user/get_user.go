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
	ID                xid.ID               `json:"id"`
	Username          string               `json:"name"`
	DisplayName       *string              `json:"display_name"`
	Bio               *string              `json:"bio"`
	MemberTitle       *string              `json:"member_title"`
	Avatar            *string              `json:"avatar"`
	Links             []string             `json:"links"`
	Names             []db.FieldEntry      `json:"names"`
	Pronouns          []db.PronounEntry    `json:"pronouns"`
	Members           []PartialMember      `json:"members"`
	Fields            []db.Field           `json:"fields"`
	CustomPreferences db.CustomPreferences `json:"custom_preferences"`
}

type GetMeResponse struct {
	GetUserResponse

	MaxInvites  int  `json:"max_invites"`
	IsAdmin     bool `json:"is_admin"`
	ListPrivate bool `json:"list_private"`

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`

	Tumblr         *string `json:"tumblr"`
	TumblrUsername *string `json:"tumblr_username"`

	Google         *string `json:"google"`
	GoogleUsername *string `json:"google_username"`

	Fediverse         *string `json:"fediverse"`
	FediverseUsername *string `json:"fediverse_username"`
	FediverseInstance *string `json:"fediverse_instance"`
}

type PartialMember struct {
	ID          xid.ID            `json:"id"`
	Name        string            `json:"name"`
	DisplayName *string           `json:"display_name"`
	Bio         *string           `json:"bio"`
	Avatar      *string           `json:"avatar"`
	Links       []string          `json:"links"`
	Names       []db.FieldEntry   `json:"names"`
	Pronouns    []db.PronounEntry `json:"pronouns"`
}

func dbUserToResponse(u db.User, fields []db.Field, members []db.Member) GetUserResponse {
	resp := GetUserResponse{
		ID:                u.ID,
		Username:          u.Username,
		DisplayName:       u.DisplayName,
		Bio:               u.Bio,
		MemberTitle:       u.MemberTitle,
		Avatar:            u.Avatar,
		Links:             db.NotNull(u.Links),
		Names:             db.NotNull(u.Names),
		Pronouns:          db.NotNull(u.Pronouns),
		Fields:            db.NotNull(fields),
		CustomPreferences: u.CustomPreferences,
	}

	resp.Members = make([]PartialMember, len(members))
	for i := range members {
		resp.Members[i] = PartialMember{
			ID:          members[i].ID,
			Name:        members[i].Name,
			DisplayName: members[i].DisplayName,
			Bio:         members[i].Bio,
			Avatar:      members[i].Avatar,
			Links:       db.NotNull(members[i].Links),
			Names:       db.NotNull(members[i].Names),
			Pronouns:    db.NotNull(members[i].Pronouns),
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
			if u.DeletedAt != nil {
				return server.APIError{Code: server.ErrUserNotFound}
			}

			isSelf := false
			if claims, ok := server.ClaimsFromContext(ctx); ok && claims.UserID == u.ID {
				isSelf = true
			}

			fields, err := s.DB.UserFields(ctx, u.ID)
			if err != nil {
				log.Errorf("Error getting user fields: %v", err)
				return err
			}

			var members []db.Member
			if !u.ListPrivate || isSelf {
				members, err = s.DB.UserMembers(ctx, u.ID, isSelf)
				if err != nil {
					log.Errorf("Error getting user members: %v", err)
					return err
				}
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

	isSelf := false
	if claims, ok := server.ClaimsFromContext(ctx); ok && claims.UserID == u.ID {
		isSelf = true
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user fields: %v", err)
		return err
	}

	var members []db.Member
	if !u.ListPrivate || isSelf {
		members, err = s.DB.UserMembers(ctx, u.ID, isSelf)
		if err != nil {
			log.Errorf("Error getting user members: %v", err)
			return err
		}
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

	members, err := s.DB.UserMembers(ctx, u.ID, true)
	if err != nil {
		log.Errorf("Error getting user members: %v", err)
		return err
	}

	render.JSON(w, r, GetMeResponse{
		GetUserResponse:   dbUserToResponse(u, fields, members),
		MaxInvites:        u.MaxInvites,
		IsAdmin:           u.IsAdmin,
		ListPrivate:       u.ListPrivate,
		Discord:           u.Discord,
		DiscordUsername:   u.DiscordUsername,
		Tumblr:            u.Tumblr,
		TumblrUsername:    u.TumblrUsername,
		Google:            u.Google,
		GoogleUsername:    u.GoogleUsername,
		Fediverse:         u.Fediverse,
		FediverseUsername: u.FediverseUsername,
		FediverseInstance: u.FediverseInstance,
	})
	return nil
}
