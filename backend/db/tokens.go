package db

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type Token struct {
	UserID      xid.ID
	TokenID     xid.ID
	Invalidated bool
	Created     time.Time
	Expires     time.Time
}

func (db *DB) TokenValid(ctx context.Context, userID, tokenID xid.ID) (valid bool, err error) {
	sql, args, err := sq.Select("*").From("tokens").
		Where("user_id = ?", userID).
		Where("token_id = ?", tokenID).
		ToSql()
	if err != nil {
		return false, errors.Wrap(err, "building sql")
	}

	var t Token
	err = pgxscan.Get(ctx, db, &t, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return false, nil
		}

		return false, errors.Wrap(err, "getting from database")
	}

	now := time.Now()
	return !t.Invalidated && t.Created.Before(now) && t.Expires.After(now), nil
}

func (db *DB) Tokens(ctx context.Context, userID xid.ID) (ts []Token, err error) {
	sql, args, err := sq.Select("*").From("tokens").
		Where("user_id = ?", userID).
		Where("expires > ?", time.Now()).
		OrderBy("created").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &ts, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "getting from database")
	}
	return ts, nil
}

// 3 months, might be customizable later
const ExpiryTime = 3 * 30 * 24 * time.Hour

// SaveToken saves a token to the database.
func (db *DB) SaveToken(ctx context.Context, userID xid.ID, tokenID xid.ID) (t Token, err error) {
	sql, args, err := sq.Insert("tokens").
		Columns("user_id", "token_id", "expires").
		Values(userID, tokenID, time.Now().Add(ExpiryTime)).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return t, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &t, sql, args...)
	if err != nil {
		return t, errors.Wrap(err, "inserting token")
	}
	return t, nil
}

func (db *DB) InvalidateToken(ctx context.Context, userID xid.ID, tokenID xid.ID) (t Token, err error) {
	sql, args, err := sq.Update("tokens").
		Where("user_id = ?").
		Where("token_id = ?").
		Set("invalidated", true).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return t, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &t, sql, args...)
	if err != nil {
		return t, errors.Wrap(err, "invalidating token")
	}
	return t, nil
}

func (db *DB) InvalidateAllTokens(ctx context.Context, q querier, userID xid.ID) error {
	sql, args, err := sq.Update("tokens").Where("user_id = ?", userID).Set("invalidated", true).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = q.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}
