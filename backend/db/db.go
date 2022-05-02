package db

import (
	"context"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mediocregopher/radix/v4"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type DB struct {
	*pgxpool.Pool

	Redis radix.Client
}

func New(dsn string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	redis, err := (&radix.PoolConfig{}).New(context.Background(), "tcp", os.Getenv("REDIS"))
	if err != nil {
		return nil, err
	}

	db := &DB{
		Pool:  pool,
		Redis: redis,
	}

	return db, nil
}

// MultiCmd executes the given Redis commands in order.
// If any return an error, the function is aborted.
func (db *DB) MultiCmd(ctx context.Context, cmds ...radix.Action) error {
	for _, cmd := range cmds {
		err := db.Redis.Do(ctx, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
