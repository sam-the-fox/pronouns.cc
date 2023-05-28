package db

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"emperror.dev/errors"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minio/minio-go/v7"
	"github.com/rs/xid"
)

type PrideFlag struct {
	ID          xid.ID  `json:"id"`
	UserID      xid.ID  `json:"-"`
	Hash        string  `json:"hash"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type UserFlag struct {
	ID     int64  `json:"-"`
	UserID xid.ID `json:"-"`
	FlagID xid.ID `json:"id"`

	Hash        string  `json:"hash"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type MemberFlag struct {
	ID       int64  `json:"-"`
	MemberID xid.ID `json:"-"`
	FlagID   xid.ID `json:"id"`

	Hash        string  `json:"hash"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

const (
	MaxPrideFlags           = 100
	MaxPrideFlagTitleLength = 100
	MaxPrideFlagDescLength  = 500
)

const (
	ErrInvalidFlagID = errors.Sentinel("invalid flag ID")
	ErrFlagNotFound  = errors.Sentinel("flag not found")
)

func (db *DB) AccountFlags(ctx context.Context, userID xid.ID) (fs []PrideFlag, err error) {
	sql, args, err := sq.Select("*").From("pride_flags").Where("user_id = ?", userID).OrderBy("id").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building query")
	}

	err = pgxscan.Select(ctx, db, &fs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	return NotNull(fs), nil
}

func (db *DB) UserFlag(ctx context.Context, flagID xid.ID) (f PrideFlag, err error) {
	sql, args, err := sq.Select("*").From("pride_flags").Where("id = ?", flagID).ToSql()
	if err != nil {
		return f, errors.Wrap(err, "building query")
	}

	err = pgxscan.Get(ctx, db, &f, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return f, ErrFlagNotFound
		}

		return f, errors.Wrap(err, "executing query")
	}
	return f, nil
}

func (db *DB) UserFlags(ctx context.Context, userID xid.ID) (fs []UserFlag, err error) {
	sql, args, err := sq.Select("u.id", "u.flag_id", "f.user_id", "f.hash", "f.name", "f.description").
		From("user_flags AS u").
		Where("u.user_id = $1", userID).
		Join("pride_flags AS f ON u.flag_id = f.id").
		OrderBy("u.id ASC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building query")
	}

	err = pgxscan.Select(ctx, db, &fs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	return NotNull(fs), nil
}

func (db *DB) MemberFlags(ctx context.Context, memberID xid.ID) (fs []MemberFlag, err error) {
	sql, args, err := sq.Select("m.id", "m.flag_id", "m.member_id", "f.hash", "f.name", "f.description").
		From("member_flags AS m").
		Where("m.member_id = $1", memberID).
		Join("pride_flags AS f ON m.flag_id = f.id").
		OrderBy("m.id ASC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building query")
	}

	err = pgxscan.Select(ctx, db, &fs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	return NotNull(fs), nil
}

func (db *DB) SetUserFlags(ctx context.Context, tx pgx.Tx, userID xid.ID, flags []xid.ID) (err error) {
	sql, args, err := sq.Delete("user_flags").Where("user_id = ?", userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting existing flags")
	}

	n, err := tx.CopyFrom(ctx, pgx.Identifier{"user_flags"}, []string{"user_id", "flag_id"},
		pgx.CopyFromSlice(len(flags), func(i int) ([]any, error) {
			return []any{userID, flags[i]}, nil
		}))
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			if pge.Code == foreignKeyViolation {
				return ErrInvalidFlagID
			}
		}

		return errors.Wrap(err, "copying new flags")
	}
	if n > 0 {
		log.Debugf("set %v flags for user %v", n, userID)
	}
	return nil
}

func (db *DB) SetMemberFlags(ctx context.Context, tx pgx.Tx, memberID xid.ID, flags []xid.ID) (err error) {
	sql, args, err := sq.Delete("member_flags").Where("member_id = ?", memberID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting existing flags")
	}

	n, err := tx.CopyFrom(ctx, pgx.Identifier{"member_flags"}, []string{"member_id", "flag_id"},
		pgx.CopyFromSlice(len(flags), func(i int) ([]any, error) {
			return []any{memberID, flags[i]}, nil
		}))
	if err != nil {
		pge := &pgconn.PgError{}
		if errors.As(err, &pge) {
			if pge.Code == foreignKeyViolation {
				return ErrInvalidFlagID
			}
		}

		return errors.Wrap(err, "copying new flags")
	}
	if n > 0 {
		log.Debugf("set %v flags for member %v", n, memberID)
	}
	return nil
}

func (db *DB) CreateFlag(ctx context.Context, tx pgx.Tx, userID xid.ID, name, desc string) (f PrideFlag, err error) {
	description := &desc
	if desc == "" {
		description = nil
	}

	sql, args, err := sq.Insert("pride_flags").
		SetMap(map[string]any{
			"id":          xid.New(),
			"hash":        "",
			"user_id":     userID.String(),
			"name":        name,
			"description": description,
		}).Suffix("RETURNING *").ToSql()
	if err != nil {
		return f, errors.Wrap(err, "building query")
	}

	err = pgxscan.Get(ctx, tx, &f, sql, args...)
	if err != nil {
		return f, errors.Wrap(err, "executing query")
	}
	return f, nil
}

func (db *DB) EditFlag(ctx context.Context, tx pgx.Tx, flagID xid.ID, name, desc, hash *string) (f PrideFlag, err error) {
	b := sq.Update("pride_flags").
		Where("id = ?", flagID)
	if name != nil {
		b = b.Set("name", *name)
	}
	if desc != nil {
		if *desc == "" {
			b = b.Set("description", nil)
		} else {
			b = b.Set("description", *desc)
		}
	}
	if hash != nil {
		b = b.Set("hash", *hash)
	}

	sql, args, err := b.Suffix("RETURNING *").ToSql()
	if err != nil {
		return f, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, tx, &f, sql, args...)
	if err != nil {
		return f, errors.Wrap(err, "executing query")
	}
	return f, nil
}

func (db *DB) WriteFlag(ctx context.Context, flagID xid.ID, flag *bytes.Buffer) (hash string, err error) {
	hasher := sha256.New()
	_, err = hasher.Write(flag.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "hashing flag")
	}
	hash = hex.EncodeToString(hasher.Sum(nil))

	_, err = db.minio.PutObject(ctx, db.minioBucket, "/flags/"+hash+".webp", flag, -1, minio.PutObjectOptions{
		ContentType: "image/webp",
	})
	if err != nil {
		return "", errors.Wrap(err, "uploading flag")
	}

	return hash, nil
}

func (db *DB) DeleteFlag(ctx context.Context, flagID xid.ID, hash string) error {
	sql, args, err := sq.Delete("pride_flags").Where("id = ?", flagID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	err = db.minio.RemoveObject(ctx, db.minioBucket, "/flags/"+flagID.String()+"/"+hash+".webp", minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "deleting flag")
	}

	return nil
}

func (db *DB) FlagObject(ctx context.Context, flagID xid.ID, hash string) (io.ReadCloser, error) {
	obj, err := db.minio.GetObject(ctx, db.minioBucket, "/flags/"+flagID.String()+"/"+hash+".webp", minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "getting object")
	}
	return obj, nil
}

// ConvertFlag parses a flag from a data URI, converts it to WebP, and returns the result.
func (db *DB) ConvertFlag(data string) (webpOut *bytes.Buffer, err error) {
	defer vips.ShutdownThread()

	data = strings.TrimSpace(data)
	if !strings.Contains(data, ",") || !strings.Contains(data, ":") || !strings.Contains(data, ";") {
		return nil, ErrInvalidDataURI
	}
	split := strings.Split(data, ",")

	rawData, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid base64 data")
	}

	image, err := vips.LoadImageFromBuffer(rawData, nil)
	if err != nil {
		return nil, errors.Wrap(err, "decoding image")
	}

	err = image.ThumbnailWithSize(256, 256, vips.InterestingNone, vips.SizeBoth)
	if err != nil {
		return nil, errors.Wrap(err, "resizing image")
	}

	webpExport := vips.NewWebpExportParams()
	webpExport.Lossless = true
	webpB, _, err := image.ExportWebp(webpExport)
	if err != nil {
		return nil, errors.Wrap(err, "exporting webp image")
	}
	webpOut = bytes.NewBuffer(webpB)

	return webpOut, nil
}
