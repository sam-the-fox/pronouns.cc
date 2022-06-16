package db

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

const (
	MaxFields           = 25
	FieldNameMaxLength  = 100
	FieldEntriesLimit   = 100
	FieldEntryMaxLength = 50
)

type Field struct {
	ID          int64    `json:"-"`
	Name        string   `json:"name"`
	Favourite   []string `json:"favourite"`
	Okay        []string `json:"okay"`
	Jokingly    []string `json:"jokingly"`
	FriendsOnly []string `json:"friends_only"`
	Avoid       []string `json:"avoid"`
}

// Validate validates this field. If it is invalid, a non-empty string is returned as error message.
func (f Field) Validate() string {
	if f.Name == "" {
		return "name cannot be empty"
	}

	if length := len([]rune(f.Name)); length > FieldNameMaxLength {
		return fmt.Sprintf("name max length is %d characters, length is %d", FieldNameMaxLength, length)
	}

	if length := len(f.Favourite) + len(f.Okay) + len(f.Jokingly) + len(f.FriendsOnly) + len(f.Avoid); length > FieldEntriesLimit {
		return fmt.Sprintf("max number of entries is %d, current number is %d", FieldEntriesLimit, length)
	}

	for i, entry := range f.Favourite {
		if length := len([]rune(entry)); length > FieldEntryMaxLength {
			return fmt.Sprintf("favourite.%d: name max length is %d characters, length is %d", i, FieldEntryMaxLength, length)
		}
	}

	for i, entry := range f.Okay {
		if length := len([]rune(entry)); length > FieldEntryMaxLength {
			return fmt.Sprintf("okay.%d: name max length is %d characters, length is %d", i, FieldEntryMaxLength, length)
		}
	}

	for i, entry := range f.Jokingly {
		if length := len([]rune(entry)); length > FieldEntryMaxLength {
			return fmt.Sprintf("jokingly.%d: name max length is %d characters, length is %d", i, FieldEntryMaxLength, length)
		}
	}

	for i, entry := range f.FriendsOnly {
		if length := len([]rune(entry)); length > FieldEntryMaxLength {
			return fmt.Sprintf("friends_only.%d: name max length is %d characters, length is %d", i, FieldEntryMaxLength, length)
		}
	}

	for i, entry := range f.Avoid {
		if length := len([]rune(entry)); length > FieldEntryMaxLength {
			return fmt.Sprintf("avoid.%d: name max length is %d characters, length is %d", i, FieldEntryMaxLength, length)
		}
	}

	return ""
}

// UserFields returns the fields associated with the given user ID.
func (db *DB) UserFields(ctx context.Context, id xid.ID) (fs []Field, err error) {
	sql, args, err := sq.
		Select("id", "name", "favourite", "okay", "jokingly", "friends_only", "avoid").
		From("user_fields").Where("user_id = ?", id).OrderBy("id ASC").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &fs, sql, args...)
	if err != nil {
		return nil, errors.Cause(err)
	}
	return fs, nil
}

// SetUserFields updates the fields for the given user.
func (db *DB) SetUserFields(ctx context.Context, tx pgx.Tx, userID xid.ID, fields []Field) (err error) {
	sql, args, err := sq.Delete("user_fields").Where("user_id = ?", userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting existing fields")
	}

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"user_fields"},
		[]string{"name", "favourite", "okay", "jokingly", "friends_only", "avoid"},
		pgx.CopyFromSlice(len(fields), func(i int) ([]any, error) {
			return []any{
				fields[i].Name,
				fields[i].Favourite,
				fields[i].Okay,
				fields[i].Jokingly,
				fields[i].FriendsOnly,
				fields[i].Avoid,
			}, nil
		}))
	if err != nil {
		return errors.Wrap(err, "inserting new fields")
	}
	return nil
}
