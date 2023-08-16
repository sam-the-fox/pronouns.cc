package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/pronounscc/pronouns.cc/backend/server"
	"emperror.dev/errors"
)

const hcaptchaURL = "https://hcaptcha.com/siteverify"

type hcaptchaResponse struct {
	Success bool `json:"success"`
}

// verifyCaptcha verifies a captcha response.
func (s *Server) verifyCaptcha(ctx context.Context, response string) (ok bool, err error) {
	vals := url.Values{
		"response": []string{response},
		"secret":   []string{s.hcaptchaSecret},
		"sitekey":  []string{s.hcaptchaSitekey},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hcaptchaURL, strings.NewReader(vals.Encode()))
	if err != nil {
		return false, errors.Wrap(err, "creating request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "pronouns.cc/"+server.Tag)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "sending request")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, errors.Sentinel("error status code")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "reading body")
	}

	var hr hcaptchaResponse
	err = json.Unmarshal(b, &hr)
	if err != nil {
		return false, errors.Wrap(err, "unmarshaling json")
	}

	return hr.Success, nil
}
