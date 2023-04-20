package db

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
)

const (
	MaxFields           = 25
	FieldNameMaxLength  = 100
	FieldEntriesLimit   = 100
	FieldEntryMaxLength = 100
)

type Field struct {
	ID      int64        `json:"-"`
	Name    string       `json:"name"`
	Entries []FieldEntry `json:"entries"`
}

// Validate validates this field. If it is invalid, a non-empty string is returned as error message.
func (f Field) Validate(custom CustomPreferences) string {
	if f.Name == "" {
		return "name cannot be empty"
	}

	if length := len([]rune(f.Name)); length > FieldNameMaxLength {
		return fmt.Sprintf("name max length is %d characters, length is %d", FieldNameMaxLength, length)
	}

	if length := len(f.Entries); length > FieldEntriesLimit {
		return fmt.Sprintf("max number of entries is %d, current number is %d", FieldEntriesLimit, length)
	}

	for i, entry := range f.Entries {
		if length := len([]rune(entry.Value)); length > FieldEntryMaxLength {
			return fmt.Sprintf("entries.%d: max length is %d characters, length is %d", i, FieldEntryMaxLength, length)
		}

		if !entry.Status.Valid(custom) {
			return fmt.Sprintf("entries.%d: status is invalid", i)
		}
	}

	return ""
}

// UserFields returns the fields associated with the given user ID.
func (db *DB) UserFields(ctx context.Context, id xid.ID) (fs []Field, err error) {
	sql, args, err := sq.Select("id", "name", "entries").From("user_fields").Where("user_id = ?", id).OrderBy("id").ToSql()
	if err != nil {
		return fs, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &fs, sql, args...)
	if err != nil {
		return fs, errors.Wrap(err, "executing query")
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

	for _, field := range fields {
		_, err := tx.Exec(ctx, "INSERT INTO user_fields (user_id, name, entries) VALUES ($1, $2, $3)", userID, field.Name, field.Entries)
		if err != nil {
			return errors.Wrap(err, "inserting new fields")
		}
	}
	return nil
}

// MemberFields returns the fields associated with the given member ID.
func (db *DB) MemberFields(ctx context.Context, id xid.ID) (fs []Field, err error) {
	sql, args, err := sq.Select("id", "name", "entries").From("member_fields").Where("member_id = ?", id).OrderBy("id").ToSql()
	if err != nil {
		return fs, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &fs, sql, args...)
	if err != nil {
		return fs, errors.Wrap(err, "executing query")
	}

	return fs, nil
}

// SetMemberFields updates the fields for the given member.
func (db *DB) SetMemberFields(ctx context.Context, tx pgx.Tx, memberID xid.ID, fields []Field) (err error) {
	sql, args, err := sq.Delete("member_fields").Where("member_id = ?", memberID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting existing fields")
	}

	for _, field := range fields {
		_, err := tx.Exec(ctx, "INSERT INTO member_fields (member_id, name, entries) VALUES ($1, $2, $3)", memberID, field.Name, field.Entries)
		if err != nil {
			return errors.Wrap(err, "inserting new fields")
		}
	}
	return nil
}
