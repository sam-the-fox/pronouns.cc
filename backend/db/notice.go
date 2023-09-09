package db

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Notice struct {
	ID        int
	Notice    string
	StartTime time.Time
	EndTime   time.Time
}

func (db *DB) Notices(ctx context.Context) (ns []Notice, err error) {
	sql, args, err := sq.Select("*").From("notices").OrderBy("id DESC").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &ns, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	return NotNull(ns), nil
}

func (db *DB) CreateNotice(ctx context.Context, notice string, start, end time.Time) (n Notice, err error) {
	sql, args, err := sq.Insert("notices").SetMap(map[string]any{
		"notice":     notice,
		"start_time": start,
		"end_time":   end,
	}).Suffix("RETURNING *").ToSql()
	if err != nil {
		return n, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &n, sql, args...)
	if err != nil {
		return n, errors.Wrap(err, "executing query")
	}
	return n, nil
}

const ErrNoNotice = errors.Sentinel("no current notice")

func (db *DB) CurrentNotice(ctx context.Context) (n Notice, err error) {
	sql, args, err := sq.Select("*").From("notices").Where("end_time > ?", time.Now()).OrderBy("id DESC").Limit(1).ToSql()
	if err != nil {
		return n, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &n, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return n, ErrNoNotice
		}

		return n, errors.Wrap(err, "executing query")
	}
	return n, nil
}
