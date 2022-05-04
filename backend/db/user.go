package db

import (
	"context"
	"regexp"

	"emperror.dev/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
	ErrUserNotFound = errors.Sentinel("user not found")

	ErrUsernameTaken    = errors.Sentinel("username is already taken")
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

// DiscordUser fetches a user by Discord user ID.
func (db *DB) DiscordUser(ctx context.Context, discordID string) (u User, err error) {
	sql, args, err := sq.Select("*").From("users").Where("discord = ?", discordID).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}
		return u, errors.Cause(err)
	}

	return u, nil
}

func (u *User) UpdateFromDiscord(ctx context.Context, db pgxscan.Querier, du *discordgo.User) error {
	builder := sq.Update("users").
		Set("discord_username", du.String()).
		Where("id = ?", u.ID).
		Suffix("RETURNING *")

	if u.AvatarSource == nil || *u.AvatarSource == "discord" {
		builder = builder.
			Set("avatar_source", "discord").
			Set("avatar_url", du.AvatarURL("1024"))
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	return pgxscan.Get(ctx, db, u, sql, args...)
}

// User gets a user by ID.
func (db *DB) User(ctx context.Context, id xid.ID) (u User, err error) {
	err = pgxscan.Get(ctx, db, &u, "select * from users where id = $1", id)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Cause(err)
	}

	return u, nil
}

// Username gets a user by username.
func (db *DB) Username(ctx context.Context, name string) (u User, err error) {
	err = pgxscan.Get(ctx, db, &u, "select * from users where username = $1", name)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Cause(err)
	}

	return u, nil
}
