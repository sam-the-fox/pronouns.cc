package member

import (
	"fmt"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
)

type PatchMemberRequest struct {
	Name        *string       `json:"name"`
	Bio         *string       `json:"bio"`
	DisplayName *string       `json:"display_name"`
	Links       *[]string     `json:"links"`
	Names       *[]db.Name    `json:"names"`
	Pronouns    *[]db.Pronoun `json:"pronouns"`
	Fields      *[]db.Field   `json:"fields"`
	Avatar      *string       `json:"avatar"`
}

func (s *Server) patchMember(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	id, err := xid.FromString(chi.URLParam(r, "memberRef"))
	if err != nil {
		return server.APIError{Code: server.ErrMemberNotFound}
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
	if req.Name != nil && len(*req.Name) > db.MaxMemberNameLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Name name too long (max %d, current %d)", db.MaxMemberNameLength, len(*req.Name)),
		}
	}
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
		return err
	}

	if err := validateSlicePtr("pronoun", req.Pronouns); err != nil {
		return err
	}

	if err := validateSlicePtr("field", req.Fields); err != nil {
		return err
	}

	// update avatar
	var avatarURLs []string = nil
	if req.Avatar != nil {
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

		webpURL, jpgURL, err := s.DB.WriteMemberAvatar(ctx, id, webp, jpg)
		if err != nil {
			log.Errorf("uploading member avatar: %v", err)
			return err
		}
		avatarURLs = []string{webpURL, jpgURL}
	}

	// start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		log.Errorf("creating transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	m, err = s.DB.UpdateMember(ctx, tx, id, req.Name, req.DisplayName, req.Bio, req.Links, avatarURLs)
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

	var (
		names    []db.Name
		pronouns []db.Pronoun
		fields   []db.Field
	)

	if req.Names != nil {
		err = s.DB.SetMemberNames(ctx, tx, id, *req.Names)
		if err != nil {
			log.Errorf("setting names for member %v: %v", id, err)
			return err
		}
		names = *req.Names
	} else {
		names, err = s.DB.MemberNames(ctx, id)
		if err != nil {
			log.Errorf("getting names for member %v: %v", id, err)
			return err
		}
	}

	if req.Pronouns != nil {
		err = s.DB.SetMemberPronouns(ctx, tx, id, *req.Pronouns)
		if err != nil {
			log.Errorf("setting pronouns for member %v: %v", id, err)
			return err
		}
		pronouns = *req.Pronouns
	} else {
		pronouns, err = s.DB.MemberPronouns(ctx, id)
		if err != nil {
			log.Errorf("getting fields for member %v: %v", id, err)
			return err
		}
	}

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

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("committing transaction: %v", err)
		return err
	}

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	// echo the updated member back on success
	render.JSON(w, r, dbMemberToMember(u, m, names, pronouns, fields))
	return nil
}
