package db

import (
	"context"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"emperror.dev/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/xid"
)

func (db *DB) initMetrics() (err error) {
	err = prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "pronouns_users_total",
		Help: "The total number of registered users",
	}, func() float64 {
		count, err := db.TotalUserCount(context.Background())
		if err != nil {
			log.Errorf("getting user count for metrics: %v", err)
		}
		return float64(count)
	}))
	if err != nil {
		return errors.Wrap(err, "registering user count gauge")
	}

	err = prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "pronouns_members_total",
		Help: "The total number of registered members",
	}, func() float64 {
		count, err := db.TotalMemberCount(context.Background())
		if err != nil {
			log.Errorf("getting member count for metrics: %v", err)
		}
		return float64(count)
	}))
	if err != nil {
		return errors.Wrap(err, "registering member count gauge")
	}

	err = prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "pronouns_users_active",
		Help: "The number of users active in the past 30 days",
	}, func() float64 {
		count, err := db.ActiveUsers(context.Background())
		if err != nil {
			log.Errorf("getting active user count for metrics: %v", err)
		}
		return float64(count)
	}))
	if err != nil {
		return errors.Wrap(err, "registering active user count gauge")
	}

	err = prometheus.Register(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "pronouns_database_latency",
		Help: "The latency to the database in nanoseconds",
	}, func() float64 {
		start := time.Now()
		_, err = db.Exec(context.Background(), "SELECT 1")
		if err != nil {
			log.Errorf("pinging database: %v", err)
			return -1
		}
		return float64(time.Since(start))
	}))
	if err != nil {
		return errors.Wrap(err, "registering database latency gauge")
	}

	db.TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pronouns_api_requests_total",
		Help: "The total number of API requests since the last restart",
	})

	return nil
}

func (db *DB) TotalUserCount(ctx context.Context) (numUsers int64, err error) {
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&numUsers)
	if err != nil {
		return 0, errors.Wrap(err, "querying user count")
	}
	return numUsers, nil
}

func (db *DB) TotalMemberCount(ctx context.Context) (numMembers int64, err error) {
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM members WHERE unlisted = false AND user_id = ANY(SELECT id FROM users WHERE deleted_at IS NULL)").Scan(&numMembers)
	if err != nil {
		return 0, errors.Wrap(err, "querying member count")
	}
	return numMembers, nil
}

const activeTime = 30 * 24 * time.Hour

func (db *DB) ActiveUsers(ctx context.Context) (numUsers int64, err error) {
	t := time.Now().Add(-activeTime)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND last_active > $1", t).Scan(&numUsers)
	if err != nil {
		return 0, errors.Wrap(err, "querying active user count")
	}
	return numUsers, nil
}

type connOrTx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

// UpdateActiveTime is called on create and update endpoints (PATCH /users/@me, POST/PATCH/DELETE /members)
func (db *DB) UpdateActiveTime(ctx context.Context, tx connOrTx, userID xid.ID) (err error) {
	sql, args, err := sq.Update("users").Set("last_active", time.Now().UTC()).Where("id = ?", userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}
