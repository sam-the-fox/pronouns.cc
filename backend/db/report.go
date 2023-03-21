package db

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/rs/xid"
)

type Report struct {
	ID         int64   `json:"id"`
	UserID     xid.ID  `json:"user_id"`
	MemberID   *xid.ID `json:"member_id"`
	Reason     string  `json:"reason"`
	ReporterID xid.ID  `json:"reporter_id"`

	CreatedAt    time.Time  `json:"created_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	AdminID      *xid.ID    `json:"admin_id"`
	AdminComment *string    `json:"admin_comment"`
}

const ReportPageSize = 100

func (db *DB) Reports(ctx context.Context, closed bool, before int) (rs []Report, err error) {
	builder := sq.Select("*").From("reports").Limit(ReportPageSize).OrderBy("id DESC")
	if before != 0 {
		builder = builder.Where("id < ?", before)
	}
	if closed {
		builder = builder.Where("resolved_at IS NOT NULL")
	} else {
		builder = builder.Where("resolved_at IS NULL")
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &rs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	if len(rs) == 0 {
		return []Report{}, nil
	}
	return rs, nil
}

func (db *DB) ReportsByUser(ctx context.Context, userID xid.ID, before int) (rs []Report, err error) {
	builder := sq.Select("*").From("reports").Where("user_id = ?", userID).Limit(ReportPageSize).OrderBy("id DESC")
	if before != 0 {
		builder = builder.Where("id < ?", before)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &rs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	if len(rs) == 0 {
		return []Report{}, nil
	}
	return rs, nil
}

func (db *DB) ReportsByReporter(ctx context.Context, reporterID xid.ID, before int) (rs []Report, err error) {
	builder := sq.Select("*").From("reports").Where("reporter_id = ?", reporterID).Limit(ReportPageSize).OrderBy("id DESC")
	if before != 0 {
		builder = builder.Where("id < ?", before)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &rs, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	if len(rs) == 0 {
		return []Report{}, nil
	}
	return rs, nil
}

func (db *DB) CreateReport(ctx context.Context, reporterID, userID xid.ID, memberID *xid.ID, reason string) (r Report, err error) {
	sql, args, err := sq.Insert("reports").SetMap(map[string]any{
		"user_id":     userID,
		"reporter_id": reporterID,
		"member_id":   memberID,
		"reason":      reason,
	}).Suffix("RETURNING *").ToSql()
	if err != nil {
		return r, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &r, sql, args...)
	if err != nil {
		return r, errors.Wrap(err, "executing query")
	}
	return r, nil
}
