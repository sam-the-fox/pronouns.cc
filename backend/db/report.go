package db

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
)

type Report struct {
	ID         int64   `json:"id"`
	UserID     xid.ID  `json:"user_id"`
	UserName   string  `json:"user_name"`
	MemberID   xid.ID  `json:"member_id"`
	MemberName *string `json:"member_name"`
	Reason     string  `json:"reason"`
	ReporterID xid.ID  `json:"reporter_id"`

	CreatedAt    time.Time  `json:"created_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	AdminID      xid.ID     `json:"admin_id"`
	AdminComment *string    `json:"admin_comment"`
}

const ReportPageSize = 100
const ErrReportNotFound = errors.Sentinel("report not found")

func (db *DB) Reports(ctx context.Context, closed bool, before int) (rs []Report, err error) {
	builder := sq.Select("*",
		"(SELECT username FROM users WHERE id = reports.user_id) AS user_name",
		"(SELECT name FROM members WHERE id = reports.member_id) AS member_name").
		From("reports").
		Limit(ReportPageSize).
		OrderBy("id DESC")
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

func (db *DB) Report(ctx context.Context, tx pgx.Tx, id int64) (r Report, err error) {
	sql, args, err := sq.Select("*").From("reports").Where("id = ?", id).ToSql()
	if err != nil {
		return r, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, tx, &r, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return r, ErrReportNotFound
		}

		return r, errors.Wrap(err, "executing query")
	}
	return r, nil
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

func (db *DB) ResolveReport(ctx context.Context, ex Execer, id int64, adminID xid.ID, comment string) error {
	sql, args, err := sq.Update("reports").
		Set("admin_id", adminID).
		Set("admin_comment", comment).
		Set("resolved_at", time.Now().UTC()).
		Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = ex.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}

type Warning struct {
	ID        int64      `json:"id"`
	UserID    xid.ID     `json:"-"`
	Reason    string     `json:"reason"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"-"`
}

func (db *DB) CreateWarning(ctx context.Context, tx pgx.Tx, userID xid.ID, reason string) (w Warning, err error) {
	sql, args, err := sq.Insert("warnings").SetMap(map[string]any{
		"user_id": userID,
		"reason":  reason,
	}).ToSql()
	if err != nil {
		return w, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, tx, &w, sql, args...)
	if err != nil {
		return w, errors.Wrap(err, "executing query")
	}
	return w, nil
}

func (db *DB) Warnings(ctx context.Context, userID xid.ID, unread bool) (ws []Warning, err error) {
	builder := sq.Select("*").From("warnings").Where("user_id = ?", userID).OrderBy("id DESC")
	if unread {
		builder = builder.Where("read_at IS NULL")
	}
	sql, args, err := builder.ToSql()

	if err != nil {
		return ws, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &ws, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	if len(ws) == 0 {
		return []Warning{}, nil
	}
	return ws, nil
}

func (db *DB) AckWarning(ctx context.Context, userID xid.ID, id int64) (ok bool, err error) {
	sql, args, err := sq.Update("warnings").
		Set("read_at", time.Now().UTC()).
		Where("user_id = ?", userID).
		Where("id = ?", id).
		Where("read_at IS NULL").ToSql()
	if err != nil {
		return false, errors.Wrap(err, "building sql")
	}

	ct, err := db.Exec(ctx, sql, args...)
	if err != nil {
		return false, errors.Wrap(err, "executing query")
	}

	if ct.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}
