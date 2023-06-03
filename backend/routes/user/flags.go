package user

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/xid"
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
		} else if err == db.ErrFileTooLarge {
			return server.APIError{Code: server.ErrBadRequest, Message: "data URI exceeds 512 KB"}
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

type patchUserFlagRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (s *Server) patchUserFlag(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	flagID, err := xid.FromString(chi.URLParam(r, "flagID"))
	if err != nil {
		return server.APIError{Code: server.ErrNotFound, Details: "Invalid flag ID"}
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
	var found bool
	for _, flag := range flags {
		if flag.ID == flagID {
			found = true
			break
		}
	}
	if !found {
		return server.APIError{Code: server.ErrNotFound, Details: "No flag with that ID found"}
	}

	var req patchUserFlagRequest
	err = render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		*req.Description = strings.TrimSpace(*req.Description)
	}

	if req.Name == nil && req.Description == nil {
		return server.APIError{Code: server.ErrBadRequest, Details: "Request cannot be empty"}
	}

	if s := common.StringLength(req.Name); s > db.MaxPrideFlagTitleLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("name too long, must be %v characters or less, is %v", db.MaxPrideFlagTitleLength, s),
		}
	}
	if s := common.StringLength(req.Description); s > db.MaxPrideFlagDescLength {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("description too long, must be %v characters or less, is %v", db.MaxPrideFlagDescLength, s),
		}
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback(ctx)

	flag, err := s.DB.EditFlag(ctx, tx, flagID, req.Name, req.Description, nil)
	if err != nil {
		return errors.Wrap(err, "updating flag")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	render.JSON(w, r, flag)
	return nil
}

func (s *Server) deleteUserFlag(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	if !claims.TokenWrite {
		return server.APIError{Code: server.ErrMissingPermissions, Details: "This token is read-only"}
	}

	flagID, err := xid.FromString(chi.URLParam(r, "flagID"))
	if err != nil {
		return server.APIError{Code: server.ErrNotFound, Details: "Invalid flag ID"}
	}

	flag, err := s.DB.UserFlag(ctx, flagID)
	if err != nil {
		if err == db.ErrFlagNotFound {
			return server.APIError{Code: server.ErrNotFound, Details: "Flag not found"}
		}

		return errors.Wrap(err, "getting flag object")
	}
	if flag.UserID != claims.UserID {
		return server.APIError{Code: server.ErrNotFound, Details: "Flag not found"}
	}

	err = s.DB.DeleteFlag(ctx, flag.ID, flag.Hash)
	if err != nil {
		return errors.Wrap(err, "deleting flag")
	}

	render.NoContent(w, r)
	return nil
}
