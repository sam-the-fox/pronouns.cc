package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

const (
	MaxMemberCount      = 500
	MaxMemberNameLength = 100
)

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

func (db *DB) getMember(ctx context.Context, q pgxscan.Querier, id xid.ID) (m Member, err error) {
	sql, args, err := sq.Select("*").From("members").Where("id = ?", id).ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, q, &m, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return m, ErrMemberNotFound
		}

		return m, errors.Wrap(err, "retrieving member")
	}
	return m, nil
}

func (db *DB) Member(ctx context.Context, id xid.ID) (m Member, err error) {
	return db.getMember(ctx, db, id)
}

// UserMember returns a member scoped by user.
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

// UserMembers returns all of a user's members, sorted by name.
func (db *DB) UserMembers(ctx context.Context, userID xid.ID) (ms []Member, err error) {
	sql, args, err := sq.Select("*").From("members").Where("user_id = ?", userID).OrderBy("name", "id").ToSql()
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

	err = pgxscan.Get(ctx, tx, &m, sql, args...)
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			// unique constraint violation
			if pge.Code == "23505" {
				return m, ErrMemberNameInUse
			}
		}

		return m, errors.Wrap(err, "executing query")
	}

	return m, nil
}

// DeleteMember deletes a member by the given ID. This is irreversible.
func (db *DB) DeleteMember(ctx context.Context, id xid.ID) (err error) {
	sql, args, err := sq.Delete("members").Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting member")
	}
	return nil
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

func (db *DB) UpdateMember(
	ctx context.Context,
	tx pgx.Tx, id xid.ID,
	name, displayName, bio *string,
	links *[]string,
	avatarURLs []string,
) (m Member, err error) {
	if name == nil && displayName == nil && bio == nil && links == nil && avatarURLs == nil {
		return db.getMember(ctx, tx, id)
	}

	builder := sq.Update("members").Where("id = ?", id)
	if name != nil {
		if *name == "" {
			builder = builder.Set("name", nil)
		} else {
			builder = builder.Set("name", *displayName)
		}
	}
	if displayName != nil {
		if *displayName == "" {
			builder = builder.Set("display_name", nil)
		} else {
			builder = builder.Set("display_name", *displayName)
		}
	}
	if bio != nil {
		if *bio == "" {
			builder = builder.Set("bio", nil)
		} else {
			builder = builder.Set("bio", *bio)
		}
	}
	if links != nil {
		if len(*links) == 0 {
			builder = builder.Set("links", nil)
		} else {
			builder = builder.Set("links", *links)
		}
	}

	if avatarURLs != nil {
		if len(avatarURLs) == 0 {
			builder = builder.Set("avatar_urls", nil)
		} else {
			builder = builder.Set("avatar_urls", avatarURLs)
		}
	}

	sql, args, err := builder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, tx, &m, sql, args...)
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			if pge.Code == "23505" {
				return m, ErrMemberNameInUse
			}
		}

		return m, errors.Wrap(err, "executing sql")
	}

	return m, nil
}
