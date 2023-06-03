package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
)

const errNoNodeinfoURL = errors.Sentinel("no valid nodeinfo rel found")

// nodeinfo queries an instance's nodeinfo and returns the software name.
func nodeinfo(ctx context.Context, instance string) (softwareName string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+instance+"/.well-known/nodeinfo", nil)
	if err != nil {
		return "", errors.Wrap(err, "creating .well-known/nodeinfo request")
	}
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "sending .well-known/nodeinfo request")
	}
	defer resp.Body.Close()

	jb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading .well-known/nodeinfo response")
	}

	var wkr wellKnownResponse
	err = json.Unmarshal(jb, &wkr)
	if err != nil {
		return "", errors.Wrap(err, "unmarshaling .well-known/nodeinfo response")
	}

	var nodeinfoURL string
	for _, link := range wkr.Links {
		if link.Rel == "http://nodeinfo.diaspora.software/ns/schema/2.0" {
			nodeinfoURL = link.Href
			break
		}
	}

	if nodeinfoURL == "" {
		return "", errNoNodeinfoURL
	}

	req, err = http.NewRequestWithContext(ctx, "GET", nodeinfoURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating nodeinfo request")
	}
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "sending nodeinfo request")
	}
	defer resp.Body.Close()

	jb, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading nodeinfo response")
	}

	var ni partialNodeinfo
	err = json.Unmarshal(jb, &ni)
	if err != nil {
		return "", errors.Wrap(err, "unmarshaling nodeinfo response")
	}

	return ni.Software.Name, nil
}

type wellKnownResponse struct {
	Links []wellKnownLink `json:"links"`
}

type wellKnownLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type partialNodeinfo struct {
	Software nodeinfoSoftware `json:"software"`
}

type nodeinfoSoftware struct {
	Name string `json:"name"`
}
