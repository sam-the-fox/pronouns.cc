package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/common"
	"codeberg.org/u1f320/pronouns.cc/backend/icons"
	"emperror.dev/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"
)

type User struct {
	ID          xid.ID
	Username    string
	DisplayName *string
	Bio         *string
	MemberTitle *string
	LastActive  time.Time

	Avatar *string
	Links  []string

	Names    []FieldEntry
	Pronouns []PronounEntry

	Discord         *string
	DiscordUsername *string

	Fediverse         *string
	FediverseUsername *string
	FediverseAppID    *int64
	FediverseInstance *string

	Tumblr         *string
	TumblrUsername *string

	Google         *string
	GoogleUsername *string

	MaxInvites  int
	IsAdmin     bool
	ListPrivate bool

	DeletedAt    *time.Time
	SelfDelete   *bool
	DeleteReason *string

	CustomPreferences CustomPreferences
}

type CustomPreferences = map[string]CustomPreference

type CustomPreference struct {
	Icon      string         `json:"icon"`
	Tooltip   string         `json:"tooltip"`
	Size      PreferenceSize `json:"size"`
	Muted     bool           `json:"muted"`
	Favourite bool           `json:"favourite"`
}

func (c CustomPreference) Validate() string {
	if !icons.IsValid(c.Icon) {
		return fmt.Sprintf("custom preference icon %q is invalid", c.Icon)
	}

	if c.Tooltip == "" {
		return "custom preference tooltip is empty"
	}
	if common.StringLength(&c.Tooltip) > FieldEntryMaxLength {
		return fmt.Sprintf("custom preference tooltip is too long, max %d characters, is %d characters", FieldEntryMaxLength, common.StringLength(&c.Tooltip))
	}

	if c.Size != PreferenceSizeLarge && c.Size != PreferenceSizeNormal && c.Size != PreferenceSizeSmall {
		return fmt.Sprintf("custom preference size %q is invalid", string(c.Size))
	}

	return ""
}

type PreferenceSize string

const (
	PreferenceSizeLarge  PreferenceSize = "large"
	PreferenceSizeNormal PreferenceSize = "normal"
	PreferenceSizeSmall  PreferenceSize = "small"
)

func (u User) NumProviders() (numProviders int) {
	if u.Discord != nil {
		numProviders++
	}
	if u.Fediverse != nil {
		numProviders++
	}
	if u.Tumblr != nil {
		numProviders++
	}
	if u.Google != nil {
		numProviders++
	}
	return numProviders
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

const (
	SelfDeleteAfter = 30 * 24 * time.Hour
	ModDeleteAfter  = 180 * 24 * time.Hour
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

	sql, args, err := sq.Insert("users").Columns("id", "username").Values(xid.New(), username).Suffix("RETURNING *").ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, tx, &u, sql, args...)
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
	return u, nil
}

func (db *DB) FediverseUser(ctx context.Context, userID string, instanceAppID int64) (u User, err error) {
	sql, args, err := sq.Select("*", "(SELECT instance FROM fediverse_apps WHERE id = users.fediverse_app_id) AS fediverse_instance").
		From("users").
		Where("fediverse = ?", userID).Where("fediverse_app_id = ?", instanceAppID).
		ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "executing query")
	}
	return u, nil
}

func (u *User) UpdateFromFedi(ctx context.Context, ex Execer, userID, username string, appID int64) error {
	sql, args, err := sq.Update("users").
		Set("fediverse", userID).
		Set("fediverse_username", username).
		Set("fediverse_app_id", appID).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Fediverse = &userID
	u.FediverseUsername = &username
	u.FediverseAppID = &appID
	return nil
}

func (u *User) UnlinkFedi(ctx context.Context, ex Execer) error {
	sql, args, err := sq.Update("users").
		Set("fediverse", nil).
		Set("fediverse_username", nil).
		Set("fediverse_app_id", nil).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Fediverse = nil
	u.FediverseUsername = nil
	u.FediverseAppID = nil
	return nil
}

// DiscordUser fetches a user by Discord user ID.
func (db *DB) DiscordUser(ctx context.Context, discordID string) (u User, err error) {
	sql, args, err := sq.Select("*", "(SELECT instance FROM fediverse_apps WHERE id = users.fediverse_app_id) AS fediverse_instance").
		From("users").Where("discord = ?", discordID).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "executing query")
	}
	return u, nil
}

func (u *User) UpdateFromDiscord(ctx context.Context, ex Execer, du *discordgo.User) error {
	sql, args, err := sq.Update("users").
		Set("discord", du.ID).
		Set("discord_username", du.String()).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Discord = &du.ID
	username := du.String()
	u.DiscordUsername = &username

	return nil
}

func (u *User) UnlinkDiscord(ctx context.Context, ex Execer) error {
	sql, args, err := sq.Update("users").
		Set("discord", nil).
		Set("discord_username", nil).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Discord = nil
	u.DiscordUsername = nil

	return nil
}

// TumblrUser fetches a user by Tumblr user ID.
func (db *DB) TumblrUser(ctx context.Context, tumblrID string) (u User, err error) {
	sql, args, err := sq.Select("*", "(SELECT instance FROM fediverse_apps WHERE id = users.fediverse_app_id) AS fediverse_instance").
		From("users").Where("tumblr = ?", tumblrID).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "executing query")
	}
	return u, nil
}

func (u *User) UpdateFromTumblr(ctx context.Context, ex Execer, tumblrID, tumblrUsername string) error {
	sql, args, err := sq.Update("users").
		Set("tumblr", tumblrID).
		Set("tumblr_username", tumblrUsername).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Tumblr = &tumblrID
	u.TumblrUsername = &tumblrUsername

	return nil
}

func (u *User) UnlinkTumblr(ctx context.Context, ex Execer) error {
	sql, args, err := sq.Update("users").
		Set("tumblr", nil).
		Set("tumblr_username", nil).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Tumblr = nil
	u.TumblrUsername = nil

	return nil
}

// GoogleUser fetches a user by Google user ID.
func (db *DB) GoogleUser(ctx context.Context, googleID string) (u User, err error) {
	sql, args, err := sq.Select("*", "(SELECT instance FROM fediverse_apps WHERE id = users.fediverse_app_id) AS fediverse_instance").
		From("users").Where("google = ?", googleID).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "executing query")
	}
	return u, nil
}

func (u *User) UpdateFromGoogle(ctx context.Context, ex Execer, googleID, googleUsername string) error {
	sql, args, err := sq.Update("users").
		Set("google", googleID).
		Set("google_username", googleUsername).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Google = &googleID
	u.GoogleUsername = &googleUsername

	return nil
}

func (u *User) UnlinkGoogle(ctx context.Context, ex Execer) error {
	sql, args, err := sq.Update("users").
		Set("google", nil).
		Set("google_username", nil).
		Where("id = ?", u.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	u.Google = nil
	u.GoogleUsername = nil

	return nil
}

// User gets a user by ID.
func (db *DB) User(ctx context.Context, id xid.ID) (u User, err error) {
	sql, args, err := sq.Select("*", "(SELECT instance FROM fediverse_apps WHERE id = users.fediverse_app_id) AS fediverse_instance").
		From("users").Where("id = ?", id).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "getting user from db")
	}

	return u, nil
}

// Username gets a user by username.
func (db *DB) Username(ctx context.Context, name string) (u User, err error) {
	sql, args, err := sq.Select("*").From("users").Where("username = ?", name).ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &u, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return u, ErrUserNotFound
		}

		return u, errors.Wrap(err, "getting user from db")
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
	memberTitle *string, listPrivate *bool,
	links *[]string,
	avatar *string,
	customPreferences *CustomPreferences,
) (u User, err error) {
	if displayName == nil && bio == nil && links == nil && avatar == nil && memberTitle == nil && listPrivate == nil && customPreferences == nil {
		sql, args, err := sq.Select("*").From("users").Where("id = ?", id).ToSql()
		if err != nil {
			return u, errors.Wrap(err, "building sql")
		}

		err = pgxscan.Get(ctx, db, &u, sql, args...)
		if err != nil {
			return u, errors.Wrap(err, "getting user from db")
		}

		return u, nil
	}

	builder := sq.Update("users").Where("id = ?", id).Suffix("RETURNING *")
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
	if memberTitle != nil {
		if *memberTitle == "" {
			builder = builder.Set("member_title", nil)
		} else {
			builder = builder.Set("member_title", *memberTitle)
		}
	}
	if links != nil {
		builder = builder.Set("links", *links)
	}
	if listPrivate != nil {
		builder = builder.Set("list_private", *listPrivate)
	}
	if customPreferences != nil {
		builder = builder.Set("custom_preferences", *customPreferences)
	}

	if avatar != nil {
		if *avatar == "" {
			builder = builder.Set("avatar", nil)
		} else {
			builder = builder.Set("avatar", avatar)
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return u, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, tx, &u, sql, args...)
	if err != nil {
		return u, errors.Wrap(err, "executing sql")
	}
	return u, nil
}

func (db *DB) DeleteUser(ctx context.Context, tx pgx.Tx, id xid.ID, selfDelete bool, reason string) error {
	builder := sq.Update("users").Set("deleted_at", time.Now().UTC()).Set("self_delete", selfDelete).Where("id = ?", id)
	if !selfDelete {
		builder = builder.Set("delete_reason", reason)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) UndoDeleteUser(ctx context.Context, id xid.ID) error {
	sql, args, err := sq.Update("users").
		Set("deleted_at", nil).
		Set("self_delete", nil).
		Set("delete_reason", nil).
		Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) ForceDeleteUser(ctx context.Context, id xid.ID) error {
	sql, args, err := sq.Delete("users").Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) DeleteUserMembers(ctx context.Context, tx pgx.Tx, id xid.ID) error {
	sql, args, err := sq.Delete("members").Where("user_id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) ResetUser(ctx context.Context, tx pgx.Tx, id xid.ID) error {
	err := db.SetUserFields(ctx, tx, id, []Field{})
	if err != nil {
		return errors.Wrap(err, "deleting fields")
	}

	hasher := sha256.New()
	_, err = hasher.Write(id.Bytes())
	if err != nil {
		return errors.Wrap(err, "hashing user id")
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	sql, args, err := sq.Update("users").
		Set("username", "deleted-"+hash).
		Set("display_name", nil).
		Set("bio", nil).
		Set("links", nil).
		Set("names", "[]").
		Set("pronouns", "[]").
		Set("avatar", nil).
		Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) CleanUser(ctx context.Context, id xid.ID) error {
	u, err := db.User(ctx, id)
	if err != nil {
		return errors.Wrap(err, "getting user")
	}

	if u.Avatar != nil {
		err = db.DeleteUserAvatar(ctx, u.ID, *u.Avatar)
		if err != nil {
			return errors.Wrap(err, "deleting user avatar")
		}
	}

	var exports []DataExport
	err = pgxscan.Select(ctx, db, &exports, "SELECT * FROM data_exports WHERE user_id = $1", u.ID)
	if err != nil {
		return errors.Wrap(err, "getting export iles")
	}

	for _, de := range exports {
		err = db.DeleteExport(ctx, de)
		if err != nil {
			continue
		}
	}

	members, err := db.UserMembers(ctx, u.ID, true)
	if err != nil {
		return errors.Wrap(err, "getting members")
	}

	for _, m := range members {
		if m.Avatar == nil {
			continue
		}

		err = db.DeleteMemberAvatar(ctx, m.ID, *m.Avatar)
		if err != nil {
			continue
		}
	}
	return nil
}
