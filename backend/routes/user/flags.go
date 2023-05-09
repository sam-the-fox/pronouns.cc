package user

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/u1f320/pronouns.cc/backend/common"
	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

func (s *Server) getUserFlags(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	flags, err := s.DB.AccountFlags(ctx, claims.UserID)
	if err != nil {
		return errors.Wrapf(err, "getting flags for account %v", claims.UserID)
	}

	render.JSON(w, r, flags)
	return nil
}

type postUserFlagRequest struct {
	Flag        string `json:"flag"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) postUserFlag(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	flags, err := s.DB.AccountFlags(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "getting current user flags")
	}
	if len(flags) >= db.MaxPrideFlags {
		return server.APIError{
			Code: server.ErrFlagLimitReached,
		}
	}

	var req postUserFlagRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// remove whitespace from all fields
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if s := common.StringLength(&req.Name); s > db.MaxPrideFlagTitleLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("name too long, must be %v characters or less, is %v", db.MaxPrideFlagTitleLength, s),
		}
	}
	if s := common.StringLength(&req.Description); s > db.MaxPrideFlagDescLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("description too long, must be %v characters or less, is %v", db.MaxPrideFlagDescLength, s),
		}
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "starting transaction")
	}
	defer tx.Rollback(ctx)

	flag, err := s.DB.CreateFlag(ctx, tx, claims.UserID, req.Name, req.Description)
	if err != nil {
		log.Errorf("creating flag: %v", err)
		return errors.Wrap(err, "creating flag")
	}

	webp, err := s.DB.ConvertFlag(req.Flag)
	if err != nil {
		if err == db.ErrInvalidDataURI {
			return server.APIError{Code: server.ErrBadRequest, Message: "invalid data URI"}
		}
		return errors.Wrap(err, "converting flag")
	}

	hash, err := s.DB.WriteFlag(ctx, flag.ID, webp)
	if err != nil {
		return errors.Wrap(err, "writing flag")
	}

	flag, err = s.DB.EditFlag(ctx, tx, flag.ID, nil, nil, &hash)
	if err != nil {
		return errors.Wrap(err, "setting hash for flag")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	render.JSON(w, r, flag)
	return nil
}

func (s *Server) patchUserFlag(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) deleteUserFlag(w http.ResponseWriter, r *http.Request) error {
	return nil
}
