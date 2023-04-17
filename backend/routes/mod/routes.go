package mod

import (
	"net/http"

	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	*server.Server
}

func Mount(srv *server.Server, r chi.Router) {
	s := &Server{Server: srv}

	r.With(MustAdmin).Route("/admin", func(r chi.Router) {
		r.Get("/reports", server.WrapHandler(s.getReports))
		r.Get("/reports/by-user/{id}", server.WrapHandler(s.getReportsByUser))
		r.Get("/reports/by-reporter/{id}", server.WrapHandler(s.getReportsByReporter))

		r.Patch("/reports/{id}", server.WrapHandler(s.resolveReport))
	})

	r.With(MustAdmin).Handle("/metrics", promhttp.Handler())

	r.With(server.MustAuth).Post("/users/{id}/reports", server.WrapHandler(s.createUserReport))
	r.With(server.MustAuth).Post("/members/{id}/reports", server.WrapHandler(s.createMemberReport))

	r.With(server.MustAuth).Get("/auth/warnings", server.WrapHandler(s.getWarnings))
	r.With(server.MustAuth).Post("/auth/warnings/{id}/ack", server.WrapHandler(s.ackWarning))
}

func MustAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, ok := server.ClaimsFromContext(r.Context())
		if !ok {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, server.APIError{
				Code:    server.ErrForbidden,
				Message: "Forbidden",
			})
			return
		}

		if !claims.UserIsAdmin {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, server.APIError{
				Code:    server.ErrForbidden,
				Message: "Forbidden",
			})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
