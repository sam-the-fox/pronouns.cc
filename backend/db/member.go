package db

import (
	"context"
	"regexp"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	Avatar      *string
	Links       []string
	Names       []FieldEntry
	Pronouns    []PronounEntry
	Unlisted    bool
}

const (
	ErrMemberNotFound  = errors.Sentinel("member not found")
	ErrMemberNameInUse = errors.Sentinel("member name already in use")
)

// member names must match this regex
var memberNameRegex = regexp.MustCompile("^[^@\\?!#/\\\\[\\]\"'$%&()+<=>^|~`,]{1,100}$")

func MemberNameValid(name string) bool {
	return memberNameRegex.MatchString(name)
}

func (db *DB) Member(ctx context.Context, id xid.ID) (m Member, err error) {
	sql, args, err := sq.Select("*").From("members").Where("id = ?", id).ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &m, sql, args...)
	if err != nil {
		return m, errors.Wrap(err, "executing query")
	}
	return m, nil
}

// UserMember returns a member scoped by user.
func (db *DB) UserMember(ctx context.Context, userID xid.ID, memberRef string) (m Member, err error) {
	sql, args, err := sq.Select("*").From("members").Where("user_id = ?", userID).Where("(id = ? or name = ?)", memberRef, memberRef).ToSql()
	if err != nil {
		return m, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &m, sql, args...)
	if err != nil {
		return m, errors.Wrap(err, "executing query")
	}
	return m, nil
}

// UserMembers returns all of a user's members, sorted by name.
func (db *DB) UserMembers(ctx context.Context, userID xid.ID, showHidden bool) (ms []Member, err error) {
	builder := sq.Select("*").
		From("members").Where("user_id = ?", userID).
		OrderBy("name", "id")
	if !showHidden {
		builder = builder.Where("unlisted = ?", false)
	}

	sql, args, err := builder.ToSql()
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
func (db *DB) CreateMember(
	ctx context.Context, tx pgx.Tx, userID xid.ID,
	name string, displayName *string, bio string, links []string,
) (m Member, err error) {
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
	unlisted *bool,
	links *[]string,
	avatar *string,
) (m Member, err error) {
	if name == nil && displayName == nil && bio == nil && links == nil && avatar == nil {
		// get member
		sql, args, err := sq.Select("*").From("members").Where("id = ?", id).ToSql()
		if err != nil {
			return m, errors.Wrap(err, "building sql")
		}

		err = pgxscan.Get(ctx, tx, &m, sql, args...)
		if err != nil {
			return m, errors.Wrap(err, "executing query")
		}
		return m, nil
	}

	builder := sq.Update("members").Where("id = ?", id).Suffix("RETURNING *")
	if name != nil {
		if *name == "" {
			return m, errors.Wrap(err, "name was empty")
		} else {
			builder = builder.Set("name", *name)
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
		builder = builder.Set("links", *links)
	}
	if unlisted != nil {
		builder = builder.Set("unlisted", *unlisted)
	}

	if avatar != nil {
		if *avatar == "" {
			builder = builder.Set("avatar", nil)
		} else {
			builder = builder.Set("avatar", avatar)
		}
	}

	sql, args, err := builder.ToSql()
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
