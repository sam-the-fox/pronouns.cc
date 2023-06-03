package prns

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"

	dbpkg "codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "shortener",
	Usage:  "URL shortener service",
	Action: run,
}

func run(c *cli.Context) error {
	port := ":" + os.Getenv("PRNS_PORT")
	baseURL := os.Getenv("BASE_URL")

	db, err := dbpkg.New()
	if err != nil {
		log.Fatalf("creating database: %v", err)
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("recovered from panic: %v", err)
			}
		}()

		id := strings.TrimPrefix(r.URL.Path, "/")
		if len(id) == 5 {
			u, err := db.UserBySID(r.Context(), id)
			if err != nil {
				if err != dbpkg.ErrUserNotFound {
					log.Errorf("getting user: %v", err)
				}

				http.Redirect(w, r, baseURL, http.StatusTemporaryRedirect)
				return
			}
			http.Redirect(w, r, baseURL+"/@"+u.Username, http.StatusTemporaryRedirect)
			return
		}

		if len(id) == 6 {
			m, err := db.MemberBySID(r.Context(), id)
			if err != nil {
				if err != dbpkg.ErrMemberNotFound {
					log.Errorf("getting member: %v", err)
				}

				http.Redirect(w, r, baseURL, http.StatusTemporaryRedirect)
				return
			}

			u, err := db.User(r.Context(), m.UserID)
			if err != nil {
				log.Errorf("getting user for member %v: %v", m.ID, err)

				http.Redirect(w, r, baseURL, http.StatusTemporaryRedirect)
				return
			}

			http.Redirect(w, r, baseURL+"/@"+u.Username+"/"+m.Name, http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, baseURL, http.StatusTemporaryRedirect)
	})

	e := make(chan error)
	go func() {
		e <- http.ListenAndServe(port, nil)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Infof("API server running at %v!", port)

	select {
	case <-ctx.Done():
		log.Info("Interrupt signal received, shutting down...")
		db.Close()
		return nil
	case err := <-e:
		log.Fatalf("Error running server: %v", err)
	}

	return nil
}
