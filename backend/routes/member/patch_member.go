package member

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/u1f320/pronouns.cc/backend/common"
	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type PatchMemberRequest struct {
	Name        *string            `json:"name"`
	Bio         *string            `json:"bio"`
	DisplayName *string            `json:"display_name"`
	Links       *[]string          `json:"links"`
	Names       *[]db.FieldEntry   `json:"names"`
	Pronouns    *[]db.PronounEntry `json:"pronouns"`
	Fields      *[]db.Field        `json:"fields"`
	Avatar      *string            `json:"avatar"`
	Unlisted    *bool              `json:"unlisted"`
	Flags       *[]xid.ID          `json:"flags"`
}

func (s *Server) patchMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	id, err := xid.FromString(chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{Code: server.ErrMemberNotFound}
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	m, err := s.DB.Member(ctx, id)
	if err != nil {
		if err == db.ErrMemberNotFound {
			return server.APIError{Code: server.ErrMemberNotFound}
		}

		return errors.Wrap(err, "getting member")
	}

	if m.UserID != claims.UserID {
		return server.APIError{Code: server.ErrNotOwnMember}
	}

	var req PatchMemberRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// validate that *something* is set
	if req.DisplayName == nil &&
		req.Name == nil &&
		req.Bio == nil &&
		req.Unlisted == nil &&
		req.Links == nil &&
		req.Fields == nil &&
		req.Names == nil &&
		req.Pronouns == nil &&
		req.Avatar == nil &&
		req.Flags == nil {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Data must not be empty",
		}
	}

	// trim whitespace from strings
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}
	if req.DisplayName != nil {
		*req.DisplayName = strings.TrimSpace(*req.DisplayName)
	}
	if req.Bio != nil {
		*req.Bio = strings.TrimSpace(*req.Bio)
	}

	if req.Name != nil && *req.Name == "" {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Name must not be empty",
		}
	} else if req.Name != nil && len(*req.Name) > 100 {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Name may not be longer than 100 characters",
		}
	}

	// validate member name
	if req.Name != nil {
		if !db.MemberNameValid(*req.Name) {
			return server.APIError{
				Code:    server.ErrBadRequest,
				Details: "Member name cannot contain any of the following: @, \\, ?, !, #, /, \\, [, ], \", ', $, %, &, (, ), +, <, =, >, ^, |, ~, `, , and cannot be one or two periods.",
			}
		}
	}

	// validate display name/bio
	if common.StringLength(req.Name) > db.MaxMemberNameLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Name name too long (max %d, current %d)", db.MaxMemberNameLength, common.StringLength(req.Name)),
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
			Details: fmt.Sprintf("Bio too long (max %d, current %d)", db.MaxUserBioLength, common.StringLength(req.Name)),
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

	if err := validateSlicePtr("name", req.Names, u.CustomPreferences); err != nil {
		return *err
	}

	if err := validateSlicePtr("pronoun", req.Pronouns, u.CustomPreferences); err != nil {
		return *err
	}

	if err := validateSlicePtr("field", req.Fields, u.CustomPreferences); err != nil {
		return *err
	}

	// update avatar
	var avatarHash *string = nil
	if req.Avatar != nil {
		if *req.Avatar == "" {
			if m.Avatar != nil {
				err = s.DB.DeleteMemberAvatar(ctx, m.ID, *m.Avatar)
				if err != nil {
					log.Errorf("deleting member avatar: %v", err)
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

				log.Errorf("converting member avatar: %v", err)
				return err
			}

			hash, err := s.DB.WriteMemberAvatar(ctx, m.ID, webp, jpg)
			if err != nil {
				log.Errorf("uploading member avatar: %v", err)
				return err
			}
			avatarHash = &hash

			// delete current avatar if member has one
			if m.Avatar != nil {
				err = s.DB.DeleteMemberAvatar(ctx, m.ID, *m.Avatar)
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

	m, err = s.DB.UpdateMember(ctx, tx, id, req.Name, req.DisplayName, req.Bio, req.Unlisted, req.Links, avatarHash)
	if err != nil {
		switch errors.Cause(err) {
		case db.ErrNothingToUpdate:
		case db.ErrMemberNameInUse:
			return server.APIError{Code: server.ErrMemberNameInUse}
		default:
			log.Errorf("updating member: %v", err)
			return errors.Wrap(err, "updating member in db")
		}

	}

	if req.Names != nil || req.Pronouns != nil {
		names := m.Names
		pronouns := m.Pronouns

		if req.Names != nil {
			names = *req.Names
		}
		if req.Pronouns != nil {
			pronouns = *req.Pronouns
		}

		err = s.DB.SetMemberNamesPronouns(ctx, tx, id, names, pronouns)
		if err != nil {
			log.Errorf("setting names for member %v: %v", id, err)
			return err
		}
		m.Names = names
		m.Pronouns = pronouns
	}

	var fields []db.Field
	if req.Fields != nil {
		err = s.DB.SetMemberFields(ctx, tx, id, *req.Fields)
		if err != nil {
			log.Errorf("setting fields for member %v: %v", id, err)
			return err
		}
		fields = *req.Fields
	} else {
		fields, err = s.DB.MemberFields(ctx, id)
		if err != nil {
			log.Errorf("getting fields for member %v: %v", id, err)
			return err
		}
	}

	// update flags
	if req.Flags != nil {
		err = s.DB.SetMemberFlags(ctx, tx, m.ID, *req.Flags)
		if err != nil {
			if err == db.ErrInvalidFlagID {
				return server.APIError{Code: server.ErrBadRequest, Details: "One or more flag IDs are unknown"}
			}

			log.Errorf("updating flags for member %v: %v", m.ID, err)
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

	// get flags to return (we need to return full flag objects, not the array of IDs in the request body)
	flags, err := s.DB.MemberFlags(ctx, m.ID)
	if err != nil {
		log.Errorf("getting user flags: %v", err)
		return err
	}

	// echo the updated member back on success
	render.JSON(w, r, dbMemberToMember(u, m, fields, flags, true))
	return nil
}
