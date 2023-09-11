package meta

import (
	"net/http"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend/db"
	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
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
	GitRepository string      `json:"git_repository"`
	GitCommit     string      `json:"git_commit"`
	Users         MetaUsers   `json:"users"`
	Members       int64       `json:"members"`
	RequireInvite bool        `json:"require_invite"`
	Notice        *MetaNotice `json:"notice"`
}

type MetaNotice struct {
	ID     int    `json:"id"`
	Notice string `json:"notice"`
}

type MetaUsers struct {
	Total       int64 `json:"total"`
	ActiveMonth int64 `json:"active_month"`
	ActiveWeek  int64 `json:"active_week"`
	ActiveDay   int64 `json:"active_day"`
}

func (s *Server) meta(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	numUsers, numMembers, activeDay, activeWeek, activeMonth := s.DB.Counts(ctx)

	var notice *MetaNotice
	if n, err := s.DB.CurrentNotice(ctx); err != nil {
		if err != db.ErrNoNotice {
			log.Errorf("getting notice: %v", err)
		}
	} else {
		notice = &MetaNotice{
			ID:     n.ID,
			Notice: n.Notice,
		}
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
		Notice:        notice,
	})
	return nil
}
