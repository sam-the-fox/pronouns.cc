package db

import (
	"context"
	"regexp"

	"codeberg.org/u1f320/pronouns.cc/backend/db/queries"
	"emperror.dev/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type User struct {
	ID          xid.ID
	Username    string
	DisplayName *string
	Bio         *string

	AvatarURLs []string `db:"avatar_urls"`
	Links      []string

	Names    []FieldEntry
	Pronouns []PronounEntry

	Discord         *string
	DiscordUsername *string

	MaxInvites int
}

// usernames must match this regex
var usernameRegex = regexp.MustCompile(`^[\w-.]{2,40}$`)

const (
	ErrUserNotFound = errors.Sentinel("user not found")

	ErrUsernameTaken    = errors.Sentinel("username is already taken")
	ErrInvalidUsername  = errors.Sentinel("username contains invalid characters")
	ErrUsernameTooShort = errors.Sentinel("username is too short")
	ErrUsernameTooLong  = errors.Sentinel("username is too long")
)

const (
	MaxUsernameLength    = 40
	MaxDisplayNameLength = 100
	MaxUserBioLength     = 1000
	MaxUserLinksLength   = 25
	MaxLinkLength        = 256
)

// CreateUser creates a user with the given username.
func (db *DB) CreateUser(ctx context.Context, tx pgx.Tx, username string) (u User, err error) {
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

	sql, args, err := sq.Insert("users").Columns("id", "username").Values(xid.New(), username).Suffix("RETURNING id").ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	var id xid.ID
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			// unique constraint violation
			if pge.Code == "23505" {
				return u, ErrUsernameTaken
			}
		}

		return u, errors.Cause(err)
	}

	return db.getUser(ctx, tx, id)
}

// DiscordUser fetches a user by Discord user ID.
func (db *DB) DiscordUser(ctx context.Context, discordID string) (u User, err error) {
	sql, args, err := sq.Select("id").From("users").Where("discord = ?", discordID).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	var id xid.ID
	err = db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return u, errors.Wrap(err, "executing id query")
	}

	return db.getUser(ctx, db, id)
}

func (u *User) UpdateFromDiscord(ctx context.Context, db querier, du *discordgo.User) error {
	sql, args, err := sq.Update("users").
		Set("discord", du.ID).
		Set("discord_username", du.String()).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Discord = &du.ID
	username := du.String()
	u.DiscordUsername = &username

	return nil
}

func (db *DB) getUser(ctx context.Context, q querier, id xid.ID) (u User, err error) {
	qu, err := queries.NewQuerier(q).GetUserByID(ctx, id.String())
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "getting user from database")
	}

	u = User{
		ID:              id,
		Username:        qu.Username,
		DisplayName:     qu.DisplayName,
		Bio:             qu.Bio,
		AvatarURLs:      qu.AvatarUrls,
		Names:           fieldEntriesFromDB(qu.Names),
		Pronouns:        pronounsFromDB(qu.Pronouns),
		Links:           qu.Links,
		Discord:         qu.Discord,
		DiscordUsername: qu.DiscordUsername,
		MaxInvites:      int(qu.MaxInvites),
	}

	return u, nil
}

// User gets a user by ID.
func (db *DB) User(ctx context.Context, id xid.ID) (u User, err error) {
	return db.getUser(ctx, db, id)
}

// Username gets a user by username.
func (db *DB) Username(ctx context.Context, name string) (u User, err error) {
	qu, err := db.q.GetUserByUsername(ctx, name)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "getting user from db")
	}

	id, err := xid.FromString(qu.ID)
	if err != nil {
		return u, errors.Wrap(err, "parsing ID")
	}

	u = User{
		ID:              id,
		Username:        qu.Username,
		DisplayName:     qu.DisplayName,
		Bio:             qu.Bio,
		AvatarURLs:      qu.AvatarUrls,
		Names:           fieldEntriesFromDB(qu.Names),
		Pronouns:        pronounsFromDB(qu.Pronouns),
		Links:           qu.Links,
		Discord:         qu.Discord,
		DiscordUsername: qu.DiscordUsername,
		MaxInvites:      int(qu.MaxInvites),
	}

	return u, nil
}

// UsernameTaken checks if the given username is already taken.
func (db *DB) UsernameTaken(ctx context.Context, username string) (valid, taken bool, err error) {
	if !usernameRegex.MatchString(username) {
		return false, false, nil
	}

	err = db.QueryRow(ctx, "select exists (select id from users where username = $1)", username).Scan(&taken)
	return true, taken, err
}

// UpdateUsername validates the given username, then updates the given user's name to it if valid.
func (db *DB) UpdateUsername(ctx context.Context, tx pgx.Tx, id xid.ID, newName string) error {
	if !usernameRegex.MatchString(newName) {
		return ErrInvalidUsername
	}

	sql, args, err := sq.Update("users").Set("username", newName).Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			// unique constraint violation
			if pge.Code == "23505" {
				return ErrUsernameTaken
			}
		}

		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) UpdateUser(
	ctx context.Context,
	tx pgx.Tx, id xid.ID,
	displayName, bio *string,
	links *[]string,
	avatarURLs []string,
) (u User, err error) {
	if displayName == nil && bio == nil && links == nil && avatarURLs == nil {
		return db.getUser(ctx, tx, id)
	}

	builder := sq.Update("users").Where("id = ?", id)
	if displayName != nil {
		if *displayName == "" {
			builder = builder.Set("display_name", nil)
		} else {
			builder = builder.Set("display_name", *displayName)
		}
	}
	if bio != nil {
		if *bio == "" {
			builder = builder.Set("bio", nil)
		} else {
			builder = builder.Set("bio", *bio)
		}
	}
	if links != nil {
		if len(*links) == 0 {
			builder = builder.Set("links", nil)
		} else {
			builder = builder.Set("links", *links)
		}
	}

	if avatarURLs != nil {
		if len(avatarURLs) == 0 {
			builder = builder.Set("avatar_urls", nil)
		} else {
			builder = builder.Set("avatar_urls", avatarURLs)
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return u, errors.Wrap(err, "executing sql")
	}

	u, err = db.getUser(ctx, tx, id)
	if err != nil {
		return u, errors.Wrap(err, "getting updated user")
	}
	return u, nil
}
