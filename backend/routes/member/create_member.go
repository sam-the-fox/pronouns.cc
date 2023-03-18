package member

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

type CreateMemberRequest struct {
	Name        string            `json:"name"`
	DisplayName *string           `json:"display_name"`
	Bio         string            `json:"bio"`
	Avatar      string            `json:"avatar"`
	Links       []string          `json:"links"`
	Names       []db.FieldEntry   `json:"names"`
	Pronouns    []db.PronounEntry `json:"pronouns"`
	Fields      []db.Field        `json:"fields"`
}

func (s *Server) createMember(w http.ResponseWriter, r *http.Request) (err error) {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	u, err := s.DB.User(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	memberCount, err := s.DB.MemberCount(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting member count")
	}
	if memberCount > db.MaxMemberCount {
		return server.APIError{
			Code: server.ErrMemberLimitReached,
		}
	}

	var cmr CreateMemberRequest
	err = render.Decode(r, &cmr)
	if err != nil {
		if _, ok := err.(server.APIError); ok {
			return err
		}

		return server.APIError{Code: server.ErrBadRequest}
	}

	// remove whitespace from all fields
	cmr.Name = strings.TrimSpace(cmr.Name)
	cmr.Bio = strings.TrimSpace(cmr.Bio)
	if cmr.DisplayName != nil {
		*cmr.DisplayName = strings.TrimSpace(*cmr.DisplayName)
	}

	// validate everything
	if cmr.Name == "" {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Name may not be empty",
		}
	} else if len(cmr.Name) > 100 {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Name may not be longer than 100 characters",
		}
	}

	if !db.MemberNameValid(cmr.Name) {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: "Member name cannot contain any of the following: @, \\, ?, !, #, /, \\, [, ], \", ', $, %, &, (, ), +, <, =, >, ^, |, ~, `, ,",
		}
	}

	if err := validateSlicePtr("name", &cmr.Names); err != nil {
		return *err
	}

	if err := validateSlicePtr("pronoun", &cmr.Pronouns); err != nil {
		return *err
	}

	if err := validateSlicePtr("field", &cmr.Fields); err != nil {
		return *err
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "starting transaction")
	}
	defer tx.Rollback(ctx)

	m, err := s.DB.CreateMember(ctx, tx, claims.UserID, cmr.Name, cmr.DisplayName, cmr.Bio, cmr.Links)
	if err != nil {
		if errors.Cause(err) == db.ErrMemberNameInUse {
			return server.APIError{Code: server.ErrMemberNameInUse}
		}

		return err
	}

	// set names, pronouns, fields
	err = s.DB.SetMemberNamesPronouns(ctx, tx, m.ID, cmr.Names, cmr.Pronouns)
	if err != nil {
		log.Errorf("setting names and pronouns for member %v: %v", m.ID, err)
		return err
	}
	m.Names = cmr.Names
	m.Pronouns = cmr.Pronouns

	err = s.DB.SetMemberFields(ctx, tx, m.ID, cmr.Fields)
	if err != nil {
		log.Errorf("setting fields for member %v: %v", m.ID, err)
		return err
	}

	if cmr.Avatar != "" {
		webp, jpg, err := s.DB.ConvertAvatar(cmr.Avatar)
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

		err = tx.QueryRow(ctx, "UPDATE members SET avatar = $1 WHERE id = $2", hash, m.ID).Scan(&m.Avatar)
		if err != nil {
			return errors.Wrap(err, "setting avatar urls in db")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	render.JSON(w, r, dbMemberToMember(u, m, cmr.Fields))
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
