package db

import (
	"bytes"
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/rs/xid"
)

type DataExport struct {
	ID        int64
	UserID    xid.ID
	Filename  string
	CreatedAt time.Time
}

func (de DataExport) Path() string {
	return "/exports/" + de.UserID.String() + "/" + de.Filename + ".zip"
}

const ErrNoExport = errors.Sentinel("no data export exists")

const KeepExportTime = 7 * 24 * time.Hour

func (db *DB) UserExport(ctx context.Context, userID xid.ID) (de DataExport, err error) {
	sql, args, err := sq.Select("*").
		From("data_exports").
		Where("user_id = ?", userID).
		OrderBy("id DESC").
		Limit(1).ToSql()
	if err != nil {
		return de, errors.Wrap(err, "building query")
	}

	err = pgxscan.Get(ctx, db, &de, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return de, ErrNoExport
		}

		return de, errors.Wrap(err, "executing sql")
	}
	return de, nil
}

const recentExport = 24 * time.Hour

func (db *DB) HasRecentExport(ctx context.Context, userID xid.ID) (hasExport bool, err error) {
	err = db.QueryRow(ctx,
		"SELECT EXISTS(SELECT * FROM data_exports WHERE user_id = $1 AND created_at > $2)",
		userID, time.Now().Add(-recentExport)).Scan(&hasExport)
	if err != nil {
		return false, errors.Wrap(err, "executing query")
	}
	return hasExport, nil
}

func (db *DB) CreateExport(ctx context.Context, userID xid.ID, filename string, file *bytes.Buffer) (de DataExport, err error) {
	de = DataExport{
		UserID:   userID,
		Filename: filename,
	}

	_, err = db.minio.PutObject(ctx, db.minioBucket, de.Path(), file, int64(file.Len()), minio.PutObjectOptions{
		ContentType:    "application/zip",
		SendContentMd5: true,
	})
	if err != nil {
		return de, errors.Wrap(err, "writing export file")
	}

	sql, args, err := sq.Insert("data_exports").Columns("user_id", "filename").Values(userID, filename).ToSql()
	if err != nil {
		return de, errors.Wrap(err, "building query")
	}

	pgxscan.Get(ctx, db, &de, sql, args...)
	if err != nil {
		return de, errors.Wrap(err, "executing sql")
	}
	return de, nil
}

func (db *DB) DeleteExport(ctx context.Context, de DataExport) (err error) {
	sql, args, err := sq.Delete("data_exports").Where("id = ?", de.ID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building query")
	}

	err = db.minio.RemoveObject(ctx, db.minioBucket, de.Path(), minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "deleting export zip")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing sql")
	}
	return nil
}
