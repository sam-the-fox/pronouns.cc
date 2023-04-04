package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
)

func (db *DB) SetUserNamesPronouns(ctx context.Context, tx pgx.Tx, userID xid.ID, names []FieldEntry, pronouns []PronounEntry) (err error) {
	sql, args, err := sq.Update("users").Set("names", names).Set("pronouns", pronouns).Where("id = ?", userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

func (db *DB) SetMemberNamesPronouns(ctx context.Context, tx pgx.Tx, memberID xid.ID, names []FieldEntry, pronouns []PronounEntry) (err error) {
	sql, args, err := sq.Update("members").Set("names", names).Set("pronouns", pronouns).Where("id = ?", memberID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}
