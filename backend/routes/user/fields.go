package user

import (
	"fmt"
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/render"
)

type PatchFieldsRequest struct {
	Fields []db.Field `json:"fields"`
}

func (s *Server) patchUserFields(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims, _ := server.ClaimsFromContext(ctx)

	var req PatchFieldsRequest
	err := render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	// max 25 fields
	if len(req.Fields) > db.MaxFields {
		return server.APIError{
			Code:    server.ErrBadRequest,
			Details: fmt.Sprintf("Too many fields (max %d, current %d)", db.MaxFields, len(req.Fields)),
		}
	}

	// validate all fields
	for i, field := range req.Fields {
		if s := field.Validate(); s != "" {
			return server.APIError{
				Code:    server.ErrBadRequest,
				Details: fmt.Sprintf("field %d: %s", i, s),
			}
		}
	}

	err = s.DB.SetUserFields(ctx, claims.UserID, req.Fields)
	if err != nil {
		log.Errorf("setting fields for user %v: %v", claims.UserID, err)
		return err
	}

	// echo the fields back on success
	render.JSON(w, r, req)
	return nil
}
