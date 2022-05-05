package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
	"gitlab.com/1f320/pronouns/backend/db"
	"gitlab.com/1f320/pronouns/backend/log"
	"gitlab.com/1f320/pronouns/backend/server"
)

type GetUserResponse struct {
	ID          xid.ID   `json:"id"`
	Username    string   `json:"username"`
	DisplayName *string  `json:"display_name"`
	Bio         *string  `json:"bio"`
	AvatarURL   *string  `json:"avatar_url"`
	Links       []string `json:"links"`
}

type GetMeResponse struct {
	ID           xid.ID   `json:"id"`
	Username     string   `json:"username"`
	DisplayName  *string  `json:"display_name"`
	Bio          *string  `json:"bio"`
	AvatarSource *string  `json:"avatar_source"`
	AvatarURL    *string  `json:"avatar_url"`
	Links        []string `json:"links"`

	Discord         *string `json:"discord"`
	DiscordUsername *string `json:"discord_username"`
}

type PartialMember struct {
	ID        xid.ID  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}

func dbUserToResponse(u db.User) GetUserResponse {
	return GetUserResponse{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		Links:       u.Links,
	}
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userRef := chi.URLParamFromCtx(ctx, "userRef")

	if id, err := xid.FromString(userRef); err == nil {
		u, err := s.DB.User(ctx, id)
		if err == nil {
			render.JSON(w, r, dbUserToResponse(u))
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

	render.JSON(w, r, dbUserToResponse(u))
	return nil
}
