package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"emperror.dev/errors"
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

// SetJSON sets the given key to v marshaled as JSON.
func (db *DB) SetJSON(ctx context.Context, key string, v any, args ...string) error {
	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "marshaling json")
	}

	cmdArgs := make([]string, 0, len(args)+2)
	cmdArgs = append(cmdArgs, key, string(b))
	cmdArgs = append(cmdArgs, args...)

	err = db.Redis.Do(ctx, radix.Cmd(nil, "SET", cmdArgs...))
	if err != nil {
		return errors.Wrap(err, "writing to Redis")
	}
	return nil
}

// GetJSON gets the given key as a JSON object.
func (db *DB) GetJSON(ctx context.Context, key string, v any) error {
	var b []byte

	err := db.Redis.Do(ctx, radix.Cmd(&b, "GET", key))
	if err != nil {
		return errors.Wrap(err, "reading from Redis")
	}

	if b == nil {
		return nil
	}

	if v == nil {
		return fmt.Errorf("nil pointer passed into GetJSON")
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return errors.Wrap(err, "unmarshaling json")
	}
	return nil
}

// GetDelJSON gets the given key as a JSON object and deletes it.
func (db *DB) GetDelJSON(ctx context.Context, key string, v any) error {
	var b []byte

	err := db.Redis.Do(ctx, radix.Cmd(&b, "GETDEL", key))
	if err != nil {
		return errors.Wrap(err, "reading from Redis")
	}

	if b == nil {
		return nil
	}

	if v == nil {
		return fmt.Errorf("nil pointer passed into GetDelJSON")
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return errors.Wrap(err, "unmarshaling json")
	}
	return nil
}
