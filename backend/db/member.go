package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type Member struct {
	ID         xid.ID
	UserID     xid.ID
	Name       string
	Bio        *string
	AvatarURLs []string `db:"avatar_urls"`
	Links      []string
}

const ErrMemberNotFound = errors.Sentinel("member not found")

func (db *DB) Member(ctx context.Context, id xid.ID) (m Member, err error) {
	sql, args, err := sq.Select("*").From("members").Where("id = ?", id).ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &m, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return m, ErrMemberNotFound
		}

		return m, errors.Wrap(err, "retrieving member")
	}
	return m, nil
}

func (db *DB) UserMember(ctx context.Context, userID xid.ID, memberRef string) (m Member, err error) {
	sql, args, err := sq.Select("*").From("members").
		Where("user_id = ? and (id = ? or name = ?)", userID, memberRef, memberRef).ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &m, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return m, ErrMemberNotFound
		}

		return m, errors.Wrap(err, "retrieving member")
	}
	return m, nil
}

func (db *DB) UserMembers(ctx context.Context, userID xid.ID) (ms []Member, err error) {
	sql, args, err := sq.Select("*").From("members").Where("user_id = ?", userID).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &ms, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving members")
	}

	if ms == nil {
		ms = make([]Member, 0)
	}
	return ms, nil
}
