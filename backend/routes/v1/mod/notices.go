package mod

import (
	"net/http"
	"time"

	"codeberg.org/pronounscc/pronouns.cc/backend/common"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/aarondl/opt/omit"
	"github.com/go-chi/render"
)

type createNoticeRequest struct {
	Notice string              `json:"notice"`
	Start  omit.Val[time.Time] `json:"start"`
	End    time.Time           `json:"end"`
}

type noticeResponse struct {
	ID        int       `json:"id"`
	Notice    string    `json:"notice"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}

func (s *Server) createNotice(w http.ResponseWriter, r *http.Request) error {
	var req createNoticeRequest
	err := render.Decode(r, &req)
	if err != nil {
		return server.APIError{Code: server.ErrBadRequest}
	}

	if common.StringLength(&req.Notice) > 2000 {
		return server.APIError{Code: server.ErrBadRequest, Details: "Notice is too long, max 2000 characters"}
	}

	start := req.Start.GetOr(time.Now())
	if req.End.IsZero() {
		return server.APIError{Code: server.ErrBadRequest, Details: "`end` is missing or invalid"}
	}

	n, err := s.DB.CreateNotice(r.Context(), req.Notice, start, req.End)
	if err != nil {
		return errors.Wrap(err, "creating notice")
	}

	render.JSON(w, r, noticeResponse{
		ID:        n.ID,
		Notice:    n.Notice,
		StartTime: n.StartTime,
		EndTime:   n.EndTime,
	})
	return nil
}
