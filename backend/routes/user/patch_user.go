package user

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

type PatchUserRequest struct {
	Username          *string               `json:"name"`
	DisplayName       *string               `json:"display_name"`
	Bio               *string               `json:"bio"`
	MemberTitle       *string               `json:"member_title"`
	Links             *[]string             `json:"links"`
	Names             *[]db.FieldEntry      `json:"names"`
	Pronouns          *[]db.PronounEntry    `json:"pronouns"`
	Fields            *[]db.Field           `json:"fields"`
	Avatar            *string               `json:"avatar"`
	ListPrivate       *bool                 `json:"list_private"`
	CustomPreferences *db.CustomPreferences `json:"custom_preferences"`
	Flags             *[]xid.ID             `json:"flags"`
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
		req.MemberTitle == nil &&
		req.ListPrivate == nil &&
		req.Links == nil &&
		req.Fields == nil &&
		req.Names == nil &&
		req.Pronouns == nil &&
		req.Avatar == nil &&
		req.CustomPreferences == nil &&
		req.Flags == nil {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Data must not be empty",
		}
	}

	// validate display name/bio
	if common.StringLength(req.Username) > db.MaxUsernameLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Name name too long (max %d, current %d)", db.MaxUsernameLength, common.StringLength(req.Username)),
		}
	}
	if common.StringLength(req.DisplayName) > db.MaxDisplayNameLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Display name too long (max %d, current %d)", db.MaxDisplayNameLength, common.StringLength(req.DisplayName)),
		}
	}
	if common.StringLength(req.Bio) > db.MaxUserBioLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Bio too long (max %d, current %d)", db.MaxUserBioLength, common.StringLength(req.Bio)),
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

	// validate flag length
	if req.Flags != nil {
		if len(*req.Flags) > db.MaxPrideFlags {
			return server.APIError{
				Code:    server.ErrBadRequest,
				Details: fmt.Sprintf("Too many flags (max %d, current %d)", len(*req.Flags), db.MaxPrideFlags),
			}
		}
	}

	// validate custom preferences
	if req.CustomPreferences != nil {
		if count := len(*req.CustomPreferences); count > db.MaxFields {
			return server.APIError{Code: server.ErrBadRequest, Details: fmt.Sprintf("Too many custom preferences (max %d, current %d)", db.MaxFields, count)}
		}

		for k, v := range *req.CustomPreferences {
			_, err := uuid.Parse(k)
			if err != nil {
				return server.APIError{Code: server.ErrBadRequest, Details: "One or more custom preference IDs is not a UUID."}
			}
			if s := v.Validate(); s != "" {
				return server.APIError{Code: server.ErrBadRequest, Details: s}
			}
		}
	}
	customPreferences := u.CustomPreferences
	if req.CustomPreferences != nil {
		customPreferences = *req.CustomPreferences
	}

	if err := validateSlicePtr("name", req.Names, customPreferences); err != nil {
		return *err
	}

	if err := validateSlicePtr("pronoun", req.Pronouns, customPreferences); err != nil {
		return *err
	}

	if err := validateSlicePtr("field", req.Fields, customPreferences); err != nil {
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

	u, err = s.DB.UpdateUser(ctx, tx, claims.UserID, req.DisplayName, req.Bio, req.MemberTitle, req.ListPrivate, req.Links, avatarHash, req.CustomPreferences)
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

	// update flags
	if req.Flags != nil {
		err = s.DB.SetUserFlags(ctx, tx, claims.UserID, *req.Flags)
		if err != nil {
			if err == db.ErrInvalidFlagID {
				return server.APIError{Code: server.ErrBadRequest, Details: "One or more flag IDs are unknown"}
			}

			log.Errorf("updating flags for user %v: %v", claims.UserID, err)
			return err
		}
	}

	// update last active time
	err = s.DB.UpdateActiveTime(ctx, tx, claims.UserID)
	if err != nil {
		log.Errorf("updating last active time for user %v: %v", claims.UserID, err)
		return err
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

	// get flags to return (we need to return full flag objects, not the array of IDs in the request body)
	flags, err := s.DB.UserFlags(ctx, u.ID)
	if err != nil {
		log.Errorf("getting user flags: %v", err)
		return err
	}

	// echo the updated user back on success
	render.JSON(w, r, GetMeResponse{
		GetUserResponse:   dbUserToResponse(u, fields, nil, flags),
		MaxInvites:        u.MaxInvites,
		IsAdmin:           u.IsAdmin,
		ListPrivate:       u.ListPrivate,
		LastSIDReroll:     u.LastSIDReroll,
		Discord:           u.Discord,
		DiscordUsername:   u.DiscordUsername,
		Tumblr:            u.Tumblr,
		TumblrUsername:    u.TumblrUsername,
		Google:            u.Google,
		GoogleUsername:    u.GoogleUsername,
		Fediverse:         u.Fediverse,
		FediverseUsername: u.FediverseUsername,
		FediverseInstance: fediInstance,
	})
	return nil
}

type validator interface {
	Validate(custom db.CustomPreferences) string
}

// validateSlicePtr validates a slice of validators.
// If the slice is nil, a nil error is returned (assuming that the field is not required)
func validateSlicePtr[T validator](typ string, slice *[]T, custom db.CustomPreferences) *server.APIError {
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
		if s := pronouns.Validate(custom); s != "" {
			return &server.APIError{
				Code:    server.ErrBadRequest,
				Details: fmt.Sprintf("%s %d: %s", typ, i+1, s),
			}
		}
	}

	return nil
}

func (s *Server) rerollUserSID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting existing user")
	}

	if time.Now().Add(-time.Hour).Before(u.LastSIDReroll) {
		return server.APIError{Code: server.ErrRerollingTooQuickly}
	}

	newID, err := s.DB.RerollUserSID(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "updating user SID")
	}

	u.SID = newID
	render.JSON(w, r, dbUserToResponse(u, nil, nil, nil))
	return nil
}
