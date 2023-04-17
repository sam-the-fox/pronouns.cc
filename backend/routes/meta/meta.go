package meta

import (
	"net/http"
	"os"

	"codeberg.org/u1f320/pronouns.cc/backend/server"
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
	GitRepository string `json:"git_repository"`
	GitCommit     string `json:"git_commit"`
	Users         int64  `json:"users"`
	Members       int64  `json:"members"`
	RequireInvite bool   `json:"require_invite"`
}

func (s *Server) meta(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	numUsers, err := s.DB.TotalUserCount(ctx)
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
		Users:         numUsers,
		Members:       numMembers,
		RequireInvite: os.Getenv("REQUIRE_INVITE") == "true",
	})
	return nil
}
