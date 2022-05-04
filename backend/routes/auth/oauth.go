package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/mediocregopher/radix/v4"
)

// numStates is the number of CSRF states stored in Redis at any one time.
// This must be an integer.
const numStates = "1000"

// setCSRFState generates a random string to use as state, then stores that in Redis.
func (s *Server) setCSRFState(ctx context.Context) (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	state := base64.URLEncoding.EncodeToString(b)

	err = s.DB.MultiCmd(ctx,
		radix.Cmd(nil, "LPUSH", "csrf", state),
		radix.Cmd(nil, "LTRIM", "csrf", "0", numStates),
	)
	return state, err
}

// validateCSRFState checks if the given state exists in Redis.
func (s *Server) validateCSRFState(ctx context.Context, state string) (matched bool, err error) {
	var num int
	err = s.DB.Redis.Do(ctx, radix.Cmd(&num, "LREM", "csrf", "1", state))
	if err != nil {
		return
	}
	return num > 0, nil
}
