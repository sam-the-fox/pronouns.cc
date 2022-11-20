package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

const MaxMemberCount = 500

type Member struct {
	ID          xid.ID
	UserID      xid.ID
	Name        string
	DisplayName *string
	Bio         *string
	AvatarURLs  []string `db:"avatar_urls"`
	Links       []string
}

const (
	ErrMemberNotFound  = errors.Sentinel("member not found")
	ErrMemberNameInUse = errors.Sentinel("member name already in use")
)

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

// CreateMember creates a member.
func (db *DB) CreateMember(ctx context.Context, tx pgx.Tx, userID xid.ID, name string, displayName *string, bio string, links []string) (m Member, err error) {
	sql, args, err := sq.Insert("members").
		Columns("user_id", "id", "name", "display_name", "bio", "links").
		Values(userID, xid.New(), name, displayName, bio, links).
		Suffix("RETURNING *").ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &m, sql, args...)
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			if pge.Code == "23505" {
				return m, ErrMemberNameInUse
			}
		}

		return m, errors.Wrap(err, "executing query")
	}

	return m, nil
}

// MemberCount returns the number of members that the given user has.
func (db *DB) MemberCount(ctx context.Context, userID xid.ID) (n int64, err error) {
	sql, args, err := sq.Select("count(id)").From("members").Where("user_id = ?", userID).ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "building sql")
	}

	err = db.QueryRow(ctx, sql, args...).Scan(&n)
	if err != nil {
		return 0, errors.Wrap(err, "executing query")
	}

	return n, nil
}
