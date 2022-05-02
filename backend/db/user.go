package db

import (
	"context"
	"regexp"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/rs/xid"
)

type User struct {
	ID          xid.ID
	Username    string
	DisplayName *string
	Bio         *string

	AvatarSource *string
	AvatarURL    *string
	Links        []string

	Discord *string
}

// usernames must match this regex
var usernameRegex = regexp.MustCompile(`[\w-.]{2,40}`)

const (
	ErrUsernameTaken = errors.Sentinel("username is already taken")

	ErrInvalidUsername  = errors.Sentinel("username contains invalid characters")
	ErrUsernameTooShort = errors.Sentinel("username is too short")
	ErrUsernameTooLong  = errors.Sentinel("username is too long")
)

// CreateUser creates a user with the given username.
func (db *DB) CreateUser(ctx context.Context, username string) (u User, err error) {
	// check if the username is valid
	// if not, return an error depending on what failed
	if !usernameRegex.MatchString(username) {
		if len(username) < 2 {
			return u, ErrUsernameTooShort
		} else if len(username) > 40 {
			return u, ErrUsernameTooLong
		}

		return u, ErrInvalidUsername
	}

	sql, args, err := sq.Insert("users").Columns("id", "username").Values(xid.New(), username).Suffix("RETURNING *").ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if v, ok := errors.Cause(err).(*pgconn.PgError); ok {
			if v.Code == "23505" { // unique constraint violation
				return u, ErrUsernameTaken
			}
		}

		return u, errors.Cause(err)
	}

	return u, nil
}
