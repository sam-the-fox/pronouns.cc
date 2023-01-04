package db

import (
	"context"
	"fmt"

	"codeberg.org/u1f320/pronouns.cc/backend/db/queries"
	"emperror.dev/errors"
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
	ID      int64        `json:"-"`
	Name    string       `json:"name"`
	Entries []FieldEntry `json:"entries"`
}

// Validate validates this field. If it is invalid, a non-empty string is returned as error message.
func (f Field) Validate() string {
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

		if entry.Status == StatusUnknown || entry.Status >= wordStatusMax {
			return fmt.Sprintf("entries.%d: status is invalid, must be between 1 and %d, is %d", i, wordStatusMax-1, entry.Status)
		}
	}

	return ""
}

// UserFields returns the fields associated with the given user ID.
func (db *DB) UserFields(ctx context.Context, id xid.ID) (fs []Field, err error) {
	qfields, err := db.q.GetUserFields(ctx, id.String())
	if err != nil {
		return nil, errors.Wrap(err, "querying fields")
	}

	fs = make([]Field, len(qfields))
	for i := range qfields {
		fs[i] = Field{
			ID:      int64(*qfields[i].ID),
			Name:    *qfields[i].Name,
			Entries: dbEntriesToFieldEntries(qfields[i].Entries),
		}
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

	querier := queries.NewQuerier(tx)
	for _, field := range fields {
		querier.InsertUserField(ctx, queries.InsertUserFieldParams{
			UserID:  userID.String(),
			Name:    field.Name,
			Entries: entriesToDBEntries(field.Entries),
		})
	}
	if err != nil {
		return errors.Wrap(err, "inserting new fields")
	}
	return nil
}

// MemberFields returns the fields associated with the given member ID.
func (db *DB) MemberFields(ctx context.Context, id xid.ID) (fs []Field, err error) {
	qfields, err := db.q.GetMemberFields(ctx, id.String())
	if err != nil {
		return nil, errors.Wrap(err, "querying fields")
	}

	fs = make([]Field, len(qfields))
	for i := range qfields {
		fs[i] = Field{
			ID:      int64(*qfields[i].ID),
			Name:    *qfields[i].Name,
			Entries: dbEntriesToFieldEntries(qfields[i].Entries),
		}
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

	querier := queries.NewQuerier(tx)
	for _, field := range fields {
		querier.InsertMemberField(ctx, queries.InsertMemberFieldParams{
			MemberID: memberID.String(),
			Name:     field.Name,
			Entries:  entriesToDBEntries(field.Entries),
		})
	}
	if err != nil {
		return errors.Wrap(err, "inserting new fields")
	}
	return nil
}
