package exporter

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"codeberg.org/u1f320/pronouns.cc/backend/db"
	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/xid"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "exporter",
	Usage:  "Data exporter service",
	Action: run,
}

type server struct {
	Router chi.Router
	DB     *db.DB

	exporting   map[xid.ID]struct{}
	exportingMu sync.Mutex
}

func run(c *cli.Context) error {
	port := ":" + os.Getenv("EXPORTER_PORT")

	db, err := db.New()
	if err != nil {
		log.Fatalf("creating database: %v", err)
		return err
	}

	s := &server{
		Router:    chi.NewRouter(),
		DB:        db,
		exporting: make(map[xid.ID]struct{}),
	}

	// set up middleware + the single route
	s.Router.Use(middleware.Recoverer)
	s.Router.Get("/start/{id}", s.startExport)

	e := make(chan error)

	// run server in another goroutine (for gracefully shutting down, see below)
	go func() {
		e <- http.ListenAndServe(port, s.Router)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Infof("API server running at %v!", port)

	select {
	case <-ctx.Done():
		log.Info("Interrupt signal received, shutting down...")
		s.DB.Close()
		return nil
	case err := <-e:
		log.Fatalf("Error running server: %v", err)
	}

	return nil
}

func (s *server) startExport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := xid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := s.DB.User(ctx, id)
	if err != nil {
		log.Errorf("getting user %v: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go s.doExport(u)

	w.WriteHeader(http.StatusAccepted)
}

func (s *server) doExport(u db.User) {
	s.exportingMu.Lock()
	if _, ok := s.exporting[u.ID]; ok {
		s.exportingMu.Unlock()
		log.Debugf("user %v is already being exported, aborting", u.ID)
		return
	}
	s.exporting[u.ID] = struct{}{}
	s.exportingMu.Unlock()

	defer func() {
		s.exportingMu.Lock()
		delete(s.exporting, u.ID)
		s.exportingMu.Unlock()
	}()

	ctx := context.Background()

	log.Debugf("[%v] starting export of user", u.ID)

	jsonBuffer := new(bytes.Buffer)
	encoder := json.NewEncoder(jsonBuffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	outBuffer := new(bytes.Buffer)
	zw := zip.NewWriter(outBuffer)
	defer zw.Close()

	w, err := zw.Create("user.json")
	if err != nil {
		log.Errorf("[%v] creating file in zip archive: %v", u.ID, err)
		return
	}

	log.Debugf("[%v] getting user fields", u.ID)

	fields, err := s.DB.UserFields(ctx, u.ID)
	if err != nil {
		log.Errorf("[%v] getting user fields: %v", u.ID, err)
		return
	}

	log.Debugf("[%v] getting user warnings", u.ID)

	warnings, err := s.DB.Warnings(ctx, u.ID, false)
	if err != nil {
		log.Errorf("[%v] getting warnings: %v", u.ID, err)
		return
	}

	log.Debugf("[%v] writing user json", u.ID)

	err = encoder.Encode(dbUserToExport(u, fields, warnings))
	if err != nil {
		log.Errorf("[%v] marshaling user: %v", u.ID, err)
		return
	}

	_, err = io.Copy(w, jsonBuffer)
	if err != nil {
		log.Errorf("[%v] writing user: %v", u.ID, err)
		return
	}
	jsonBuffer.Reset()

	if u.Avatar != nil {
		log.Debugf("[%v] getting user avatar", u.ID)

		w, err := zw.Create("user_avatar.webp")
		if err != nil {
			log.Errorf("[%v] creating file in zip archive: %v", u.ID, err)
			return
		}

		r, err := s.DB.UserAvatar(ctx, u.ID, *u.Avatar)
		if err != nil {
			log.Errorf("[%v] getting user avatar: %v", u.ID, err)
			return
		}
		defer r.Close()

		_, err = io.Copy(w, r)
		if err != nil {
			log.Errorf("[%v] writing user avatar: %v", u.ID, err)
			return
		}

		log.Debugf("[%v] exported user avatar", u.ID)
	}

	members, err := s.DB.UserMembers(ctx, u.ID, true)
	if err != nil {
		log.Errorf("[%v] getting user members: %v", u.ID, err)
		return
	}

	for _, m := range members {
		log.Debugf("[%v] starting export for member %v", u.ID, m.ID)

		fields, err := s.DB.MemberFields(ctx, m.ID)
		if err != nil {
			log.Errorf("[%v] getting fields for member %v: %v", u.ID, m.ID, err)
			return
		}

		w, err := zw.Create("members/" + m.Name + "-" + m.ID.String() + ".json")
		if err != nil {
			log.Errorf("[%v] creating file in zip archive: %v", u.ID, err)
			return
		}

		err = encoder.Encode(dbMemberToExport(m, fields))
		if err != nil {
			log.Errorf("[%v] marshaling member %v: %v", u.ID, m.ID, err)
			return
		}

		_, err = io.Copy(w, jsonBuffer)
		if err != nil {
			log.Errorf("[%v] writing member %v json: %v", u.ID, m.ID, err)
			return
		}
		jsonBuffer.Reset()

		if m.Avatar != nil {
			log.Debugf("[%v] getting member %v avatar", u.ID, m.ID)

			w, err := zw.Create("members/" + m.Name + "-" + m.ID.String() + "-avatar.webp")
			if err != nil {
				log.Errorf("[%v] creating file in zip archive: %v", u.ID, err)
				return
			}

			r, err := s.DB.MemberAvatar(ctx, m.ID, *m.Avatar)
			if err != nil {
				log.Errorf("[%v] getting member %v avatar: %v", u.ID, m.ID, err)
				return
			}
			defer r.Close()

			_, err = io.Copy(w, r)
			if err != nil {
				log.Errorf("[%v] writing member %v avatar: %v", u.ID, m.ID, err)
				return
			}

			log.Debugf("[%v] exported member %v avatar", u.ID, m.ID)
		}

		log.Debugf("[%v] finished export for member %v", u.ID, m.ID)
	}

	log.Debugf("[%v] finished export, uploading to object storage and saving in database", u.ID)

	err = zw.Close()
	if err != nil {
		log.Errorf("[%v] closing zip file: %v", u.ID, err)
		return
	}

	de, err := s.DB.CreateExport(ctx, u.ID, randomFilename(), outBuffer)
	if err != nil {
		log.Errorf("[%v] writing export: %v", u.ID, err)
		return
	}

	log.Debugf("[%v] finished writing export. path: %q", u.ID, de.Path())
}

func randomFilename() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(b)
}
