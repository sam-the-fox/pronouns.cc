package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"codeberg.org/u1f320/pronouns.cc/frontend"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := ":" + os.Getenv("PORT")

	s, err := server.New()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	// mount api routes
	mountRoutes(s)

	r := chi.NewMux()
	setupFrontend(r, s)

	e := make(chan error)

	// run server in another goroutine (for gracefully shutting down, see below)
	go func() {
		e <- http.ListenAndServe(port, r)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Infof("API server running at %v!", port)

	select {
	case <-ctx.Done():
		log.Info("Interrupt signal received, shutting down...")
		s.DB.Close()
		return
	case err := <-e:
		log.Fatalf("Error running server: %v", err)
	}
}

func setupFrontend(r chi.Router, s *server.Server) {
	r.Use(middleware.Recoverer)

	r.Get("/@{user}", a)
	r.Get("/@{user}/{member}", a)

	r.Mount("/api", s.Router)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		var (
			data []byte
			err  error
		)

		if strings.HasSuffix(r.URL.Path, ".js") {
			data, err = frontend.Data.ReadFile("dist" + r.URL.Path)
			w.Header().Add("content-type", "application/javascript")
		} else if strings.HasSuffix(r.URL.Path, ".css") {
			data, err = frontend.Data.ReadFile("dist" + r.URL.Path)
			w.Header().Add("content-type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".map") {
			data, err = frontend.Data.ReadFile("dist" + r.URL.Path)
		} else {
			data, err = frontend.Data.ReadFile("dist/index.html")
			w.Header().Add("content-type", "text/html")
		}
		if err != nil {
			panic(err)
		}

		w.Write(data)
	})
}

func a(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	member := chi.URLParam(r, "member")

	fmt.Fprintf(w, "user: %v, member: %v", user, member)
}
