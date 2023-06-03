package meta

import (
	"net/http"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{Server: srv}

	r.Get("/meta", server.WrapHandler(s.meta))
}

type MetaResponse struct {
	GitRepository string    `json:"git_repository"`
	GitCommit     string    `json:"git_commit"`
	Users         MetaUsers `json:"users"`
	Members       int64     `json:"members"`
	RequireInvite bool      `json:"require_invite"`
}

type MetaUsers struct {
	Total       int64 `json:"total"`
	ActiveMonth int64 `json:"active_month"`
	ActiveWeek  int64 `json:"active_week"`
	ActiveDay   int64 `json:"active_day"`
}

func (s *Server) meta(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	numUsers, err := s.DB.TotalUserCount(ctx)
	if err != nil {
		return errors.Wrap(err, "querying user count")
	}

	activeMonth, err := s.DB.ActiveUsers(ctx, db.ActiveMonth)
	if err != nil {
		return errors.Wrap(err, "querying user count")
	}

	activeWeek, err := s.DB.ActiveUsers(ctx, db.ActiveWeek)
	if err != nil {
		return errors.Wrap(err, "querying user count")
	}

	activeDay, err := s.DB.ActiveUsers(ctx, db.ActiveDay)
	if err != nil {
		return errors.Wrap(err, "querying user count")
	}

	numMembers, err := s.DB.TotalMemberCount(ctx)
	if err != nil {
		return errors.Wrap(err, "querying user count")
	}

	render.JSON(w, r, MetaResponse{
		GitRepository: server.Repository,
		GitCommit:     server.Revision,
		Users: MetaUsers{
			Total:       numUsers,
			ActiveMonth: activeMonth,
			ActiveWeek:  activeWeek,
			ActiveDay:   activeDay,
		},
		Members:       numMembers,
		RequireInvite: os.Getenv("REQUIRE_INVITE") == "true",
	})
	return nil
}
