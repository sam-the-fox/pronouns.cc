package db

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/rs/xid"
)

type Report struct {
	ID         int64
	UserID     xid.ID
	MemberID   *xid.ID
	Reason     string
	ReporterID xid.ID

	CreatedAt    time.Time
	ResolvedAt   *time.Time
	AdminID      *xid.ID
	AdminComment *string
}

const reportPageSize = 100

func (db *DB) Reports(ctx context.Context, closed bool, page int) (rs []Report, err error) {
	builder := sq.Select("*").From("reports").Offset(uint64(reportPageSize * page)).Limit(reportPageSize).OrderBy("id ASC")
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
