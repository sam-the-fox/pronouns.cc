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
	DisplayName *string     `json:"display_name"`
	Bio         *string     `json:"bio"`
	Links       *[]string   `json:"links"`
	Fields      *[]db.Field `json:"fields"`
}

// patchUser parses a PatchUserRequest and updates the user with the given ID.
func (s *Server) patchUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, _ := server.ClaimsFromContext(ctx)

	var req PatchUserRequest
	err := render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// validate that *something* is set
	if req.DisplayName == nil && req.Bio == nil && req.Links == nil && req.Fields == nil {
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

	if (req.Fields) != nil {
		// max 25 fields
		if len(*req.Fields) > db.MaxFields {
			return server.APIError{
				Code:    server.ErrBadRequest,
				Details: fmt.Sprintf("Too many fields (max %d, current %d)", db.MaxFields, len(*req.Fields)),
			}
		}

		// validate all fields
		for i, field := range *req.Fields {
			if s := field.Validate(); s != "" {
				return server.APIError{
					Code:    server.ErrBadRequest,
					Details: fmt.Sprintf("field %d: %s", i, s),
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

	u, err := s.DB.UpdateUser(ctx, tx, claims.UserID, req.DisplayName, req.Bio, req.Links)
	if err != nil && errors.Cause(err) != db.ErrNothingToUpdate {
		log.Errorf("updating user: %v", err)
		return err
	}

	var fields []db.Field
	if req.Fields != nil {
		err = s.DB.SetUserFields(ctx, tx, claims.UserID, *req.Fields)
		if err != nil {
			log.Errorf("setting fields for user %v: %v", claims.UserID, err)
			return err
		}
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

	// echo the updated user back on success
	render.JSON(w, r, dbUserToResponse(u, fields))
	return nil
}
