package user

import (
	"fmt"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

type PatchUserRequest struct {
	Username    *string            `json:"username"`
	DisplayName *string            `json:"display_name"`
	Bio         *string            `json:"bio"`
	Links       *[]string          `json:"links"`
	Names       *[]db.FieldEntry   `json:"names"`
	Pronouns    *[]db.PronounEntry `json:"pronouns"`
	Fields      *[]db.Field        `json:"fields"`
	Avatar      *string            `json:"avatar"`
}

// patchUser parses a PatchUserRequest and updates the user with the given ID.
func (s *Server) patchUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	var req PatchUserRequest
	err := render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// get existing user, for comparison later
	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting existing user")
	}

	// validate that *something* is set
	if req.Username == nil &&
		req.DisplayName == nil &&
		req.Bio == nil &&
		req.Links == nil &&
		req.Fields == nil &&
		req.Names == nil &&
		req.Pronouns == nil &&
		req.Avatar == nil {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Data must not be empty",
		}
	}

	// validate display name/bio
	if req.DisplayName != nil && len(*req.DisplayName) > db.MaxDisplayNameLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Display name too long (max %d, current %d)", db.MaxDisplayNameLength, len(*req.DisplayName)),
		}
	}
	if req.Bio != nil && len(*req.Bio) > db.MaxUserBioLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Bio too long (max %d, current %d)", db.MaxUserBioLength, len(*req.Bio)),
		}
	}

	// validate links
	if req.Links != nil {
		if len(*req.Links) > db.MaxUserLinksLength {
			return server.APIError{
				Code:    server.ErrBadRequest,
				Details: fmt.Sprintf("Too many links (max %d, current %d)", db.MaxUserLinksLength, len(*req.Links)),
			}
		}

		for i, link := range *req.Links {
			if len(link) > db.MaxLinkLength {
				return server.APIError{
					Code:    server.ErrBadRequest,
					Details: fmt.Sprintf("Link %d too long (max %d, current %d)", i, db.MaxLinkLength, len(link)),
				}
			}
		}
	}

	if err := validateSlicePtr("name", req.Names); err != nil {
		return *err
	}

	if err := validateSlicePtr("pronoun", req.Pronouns); err != nil {
		return *err
	}

	if err := validateSlicePtr("field", req.Fields); err != nil {
		return *err
	}

	// update avatar
	var avatarHash *string = nil
	if req.Avatar != nil {
		if *req.Avatar == "" {
			if u.Avatar != nil {
				err = s.DB.DeleteUserAvatar(ctx, u.ID, *u.Avatar)
				if err != nil {
					log.Errorf("deleting user avatar: %v", err)
					return errors.Wrap(err, "deleting avatar")
				}
			}
			avatarHash = req.Avatar
		} else {
			webp, jpg, err := s.DB.ConvertAvatar(*req.Avatar)
			if err != nil {
				if err == db.ErrInvalidDataURI {
					return server.APIError{
						Code:    server.ErrBadRequest,
						Details: "invalid avatar data URI",
					}
				} else if err == db.ErrInvalidContentType {
					return server.APIError{
						Code:    server.ErrBadRequest,
						Details: "invalid avatar content type",
					}
				}

				log.Errorf("converting user avatar: %v", err)
				return err
			}

			hash, err := s.DB.WriteUserAvatar(ctx, claims.UserID, webp, jpg)
			if err != nil {
				log.Errorf("uploading user avatar: %v", err)
				return err
			}
			avatarHash = &hash

			// delete current avatar if user has one
			if u.Avatar != nil {
				err = s.DB.DeleteUserAvatar(ctx, claims.UserID, *u.Avatar)
				if err != nil {
					log.Errorf("deleting existing avatar: %v", err)
				}
			}
		}
	}

	// start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		log.Errorf("creating transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// update username
	if req.Username != nil && *req.Username != u.Username {
		err = s.DB.UpdateUsername(ctx, tx, claims.UserID, *req.Username)
		if err != nil {
			switch err {
			case db.ErrUsernameTaken:
				return server.APIError{Code: server.ErrUsernameTaken}
			case db.ErrInvalidUsername:
				return server.APIError{Code: server.ErrInvalidUsername}
			default:
				return errors.Wrap(err, "updating username")
			}
		}
	}

	u, err = s.DB.UpdateUser(ctx, tx, claims.UserID, req.DisplayName, req.Bio, req.Links, avatarHash)
	if err != nil && errors.Cause(err) != db.ErrNothingToUpdate {
		log.Errorf("updating user: %v", err)
		return err
	}

	if req.Names != nil || req.Pronouns != nil {
		names := u.Names
		pronouns := u.Pronouns

		if req.Names != nil {
			names = *req.Names
		}
		if req.Pronouns != nil {
			pronouns = *req.Pronouns
		}

		err = s.DB.SetUserNamesPronouns(ctx, tx, claims.UserID, names, pronouns)
		if err != nil {
			log.Errorf("setting names for member %v: %v", claims.UserID, err)
			return err
		}
		u.Names = names
		u.Pronouns = pronouns
	}

	var fields []db.Field
	if req.Fields != nil {
		err = s.DB.SetUserFields(ctx, tx, claims.UserID, *req.Fields)
		if err != nil {
			log.Errorf("setting fields for user %v: %v", claims.UserID, err)
			return err
		}
		fields = *req.Fields
	} else {
		fields, err = s.DB.UserFields(ctx, claims.UserID)
		if err != nil {
			log.Errorf("getting fields for user %v: %v", claims.UserID, err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("committing transaction: %v", err)
		return err
	}

	// get fedi instance name if the user has a linked fedi account
	var fediInstance *string
	if u.FediverseAppID != nil {
		app, err := s.DB.FediverseAppByID(ctx, *u.FediverseAppID)
		if err == nil {
			fediInstance = &app.Instance
		}
	}

	// echo the updated user back on success
	render.JSON(w, r, GetMeResponse{
		GetUserResponse:   dbUserToResponse(u, fields, nil),
		MaxInvites:        u.MaxInvites,
		IsAdmin:           u.IsAdmin,
		Discord:           u.Discord,
		DiscordUsername:   u.DiscordUsername,
		Fediverse:         u.Fediverse,
		FediverseUsername: u.FediverseUsername,
		FediverseInstance: fediInstance,
	})
	return nil
}

type validator interface {
	Validate() string
}

// validateSlicePtr validates a slice of validators.
// If the slice is nil, a nil error is returned (assuming that the field is not required)
func validateSlicePtr[T validator](typ string, slice *[]T) *server.APIError {
	if slice == nil {
		return nil
	}

	max := db.MaxFields
	if typ != "field" {
		max = db.FieldEntriesLimit
	}

	// max 25 fields
	if len(*slice) > max {
		return &server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Too many %ss (max %d, current %d)", typ, max, len(*slice)),
		}
	}

	// validate all fields
	for i, pronouns := range *slice {
		if s := pronouns.Validate(); s != "" {
			return &server.APIError{
				Code:    server.ErrBadRequest,
				Details: fmt.Sprintf("%s %d: %s", typ, i+1, s),
			}
		}
	}

	return nil
}
