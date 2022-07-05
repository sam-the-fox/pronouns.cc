package main

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"codeberg.org/u1f320/pronouns.cc/backend/server"
	"codeberg.org/u1f320/pronouns.cc/frontend"
	"emperror.dev/errors"

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

	r.Get("/@{user}", nil)
	r.Get("/@{user}/{member}", nil)

	r.Mount("/api", s.Router)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		data, err := frontend.Data.ReadFile("dist" + r.URL.Path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				data, err = frontend.Data.ReadFile("dist/index.html")
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					log.Errorf("serving frontend file: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			log.Errorf("serving frontend file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Add("content-type", "application/javascript")
		} else if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Add("content-type", "text/css")
		} else {
			w.Header().Add("content-type", "text/html")
		}

		w.Write(data)
	})
}
