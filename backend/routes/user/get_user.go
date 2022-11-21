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
	ID          xid.ID          `json:"id"`
	Username    string          `json:"name"`
	DisplayName *string         `json:"display_name"`
	Bio         *string         `json:"bio"`
	AvatarURLs  []string        `json:"avatar_urls"`
	Links       []string        `json:"links"`
	Names       []db.Name       `json:"names"`
	Pronouns    []db.Pronoun    `json:"pronouns"`
	Members     []PartialMember `json:"members"`
	Fields      []db.Field      `json:"fields"`
}

type GetMeResponse struct {
	GetUserResponse

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`
}

type PartialMember struct {
	ID        xid.ID  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}

func dbUserToResponse(u db.User, fields []db.Field, names []db.Name, pronouns []db.Pronoun) GetUserResponse {
	return GetUserResponse{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURLs:  u.AvatarURLs,
		Links:       u.Links,
		Names:       names,
		Pronouns:    pronouns,
		Fields:      fields,
	}
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

			names, err := s.DB.UserNames(ctx, u.ID)
			if err != nil {
				log.Errorf("getting user names: %v", err)
				return err
			}

			pronouns, err := s.DB.UserPronouns(ctx, u.ID)
			if err != nil {
				log.Errorf("getting user pronouns: %v", err)
				return err
			}

			render.JSON(w, r, dbUserToResponse(u, fields, names, pronouns))
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

	names, err := s.DB.UserNames(ctx, u.ID)
	if err != nil {
		log.Errorf("getting user names: %v", err)
		return err
	}

	pronouns, err := s.DB.UserPronouns(ctx, u.ID)
	if err != nil {
		log.Errorf("getting user pronouns: %v", err)
		return err
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user fields: %v", err)
		return err
	}

	render.JSON(w, r, dbUserToResponse(u, fields, names, pronouns))
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

	names, err := s.DB.UserNames(ctx, u.ID)
	if err != nil {
		log.Errorf("getting user names: %v", err)
		return err
	}

	pronouns, err := s.DB.UserPronouns(ctx, u.ID)
	if err != nil {
		log.Errorf("getting user pronouns: %v", err)
		return err
	}

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		log.Errorf("Error getting user fields: %v", err)
		return err
	}

	render.JSON(w, r, GetMeResponse{
		GetUserResponse: dbUserToResponse(u, fields, names, pronouns),
		Discord:         u.Discord,
		DiscordUsername: u.DiscordUsername,
	})
	return nil
}
