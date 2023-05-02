package db

import (
	"context"
	"time"

	"codeberg.org/u1f320/pronouns.cc/backend/log"
	"emperror.dev/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
