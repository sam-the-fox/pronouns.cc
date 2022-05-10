package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/rs/xid"
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
