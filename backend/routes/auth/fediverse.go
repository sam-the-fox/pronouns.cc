package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
	"github.com/go-chi/render"
)

type FediResponse struct {
	URL string `json:"url"`
}

func (s *Server) getFediverseURL(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	instance := r.FormValue("instance")
	if instance == "" {
		return server.APIError{Code: server.ErrBadRequest, Details: "Instance URL is empty"}
	}

	// Too many people tried using @username@fediverse.example despite the warning
	if strings.Contains(instance, "@") {
		return server.APIError{Code: server.ErrBadRequest, Details: "Instance URL should only be the base URL, without username"}
	}

	app, err := s.DB.FediverseApp(ctx, instance)
	if err != nil {
		return s.noAppFediverseURL(ctx, w, r, instance)
	}

	if app.Misskey() {
		_, url, err := s.misskeyURL(ctx, app)
		if err != nil {
			return errors.Wrap(err, "generating misskey URL")
		}

		render.JSON(w, r, FediResponse{
			URL: url,
		})
		return nil
	}

	state, err := s.setCSRFState(r.Context())
	if err != nil {
		return errors.Wrap(err, "setting CSRF state")
	}

	render.JSON(w, r, FediResponse{
		URL: app.ClientConfig().AuthCodeURL(state),
	})
	return nil
}

func (s *Server) noAppFediverseURL(ctx context.Context, w http.ResponseWriter, r *http.Request, instance string) error {
	softwareName, err := nodeinfo(ctx, instance)
	if err != nil {
		log.Errorf("querying instance %q nodeinfo: %v", instance, err)
	}

	switch softwareName {
	case "misskey", "foundkey", "calckey", "firefish":
		return s.noAppMisskeyURL(ctx, w, r, softwareName, instance)
	case "mastodon", "pleroma", "akkoma", "pixelfed", "gotosocial":
	case "glitchcafe", "hometown":
		// plural.cafe (potentially other instances too?) runs Mastodon but changes the software name
		// Hometown is a lightweight fork of Mastodon so we can just treat it the same
		// changing it back to mastodon here for consistency
		softwareName = "mastodon"
	default:
		return server.APIError{Code: server.ErrUnsupportedInstance}
	}

	log.Debugf("creating application on mastodon-compatible instance %q", instance)

	formData := url.Values{
		"client_name":   {"pronouns.cc (+" + s.BaseURL + ")"},
		"redirect_uris": {s.BaseURL + "/auth/login/mastodon/" + instance},
		"scopes":        {"read:accounts"},
		"website":       {s.BaseURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://"+instance+"/api/v1/apps", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Errorf("creating POST apps request for %q: %v", instance, err)
		return errors.Wrap(err, "creating POST apps request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("sending POST apps request for %q: %v", instance, err)
		return errors.Wrap(err, "sending POST apps request")
	}
	defer resp.Body.Close()

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("reading response for request: %v", err)
		return errors.Wrap(err, "reading response")
	}

	var ma mastodonApplication
	err = json.Unmarshal(jb, &ma)
	if err != nil {
		return errors.Wrap(err, "unmarshaling mastodon app")
	}

	app, err := s.DB.CreateFediverseApp(ctx, instance, softwareName, ma.ClientID, ma.ClientSecret)
	if err != nil {
		log.Errorf("saving app for %q: %v", instance, err)
		return errors.Wrap(err, "creating app")
	}

	state, err := s.setCSRFState(r.Context())
	if err != nil {
		return errors.Wrap(err, "setting CSRF state")
	}

	render.JSON(w, r, FediResponse{
		URL: app.ClientConfig().AuthCodeURL(state),
	})
	return nil
}

type mastodonApplication struct {
	Name         string `json:"name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
