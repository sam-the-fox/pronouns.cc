package db

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type Invite struct {
	UserID  xid.ID
	Code    string
	Created time.Time
	Used    bool
}

func (db *DB) UserInvites(ctx context.Context, userID xid.ID) (is []Invite, err error) {
	sql, args, err := sq.Select("*").From("invites").Where("user_id = ?", userID).OrderBy("created").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &is, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying database")
	}
	if len(is) == 0 {
		is = []Invite{}
	}

	return is, nil
}

const ErrTooManyInvites = errors.Sentinel("user invite limit reached")

func (db *DB) CreateInvite(ctx context.Context, userID xid.ID) (i Invite, err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return i, errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback(ctx)

	var maxInvites, inviteCount int
	err = tx.QueryRow(ctx, "SELECT max_invites FROM users WHERE id = $1", userID).Scan(&maxInvites)
	if err != nil {
		return i, errors.Wrap(err, "querying invite limit")
	}
	err = tx.QueryRow(ctx, "SELECT count(*) FROM invites WHERE user_id = $1", userID).Scan(&inviteCount)
	if err != nil {
		return i, errors.Wrap(err, "querying current invite count")
	}

	if inviteCount >= maxInvites {
		return i, ErrTooManyInvites
	}

	b := make([]byte, 32)

	_, err = rand.Read(b)
	if err != nil {
		panic(err)
	}

	code := base64.RawURLEncoding.EncodeToString(b)

	sql, args, err := sq.Insert("invites").Columns("user_id", "code").Values(userID, code).Suffix("RETURNING *").ToSql()
	if err != nil {
		return i, errors.Wrap(err, "building insert invite sql")
	}

	err = pgxscan.Get(ctx, db, &i, sql, args...)
	if err != nil {
		return i, errors.Wrap(err, "inserting invite")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return i, errors.Wrap(err, "committing transaction")
	}
	return i, nil
}

func (db *DB) InvalidateInvite(ctx context.Context, tx pgx.Tx, code string) (valid, alreadyUsed bool, err error) {
	err = tx.QueryRow(ctx, "SELECT used FROM invites WHERE code = $1", code).Scan(&alreadyUsed)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return false, false, nil
		}

		return false, false, errors.Wrap(err, "checking if invite exists and is used")
	}

	// valid: true, already used: true
	if alreadyUsed {
		return true, true, nil
	}

	// invite is valid, not already used
	_, err = tx.Exec(ctx, "UPDATE invites SET used = true WHERE code = $1", code)
	if err != nil {
		return false, false, errors.Wrap(err, "updating invite usage")
	}

	// valid: true, already used: false
	return true, false, nil
}
