package db

import (
	"context"
	"os"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"golang.org/x/oauth2"
)

type FediverseApp struct {
	ID int64
	// Instance is the instance's base API url, excluding schema
	Instance     string
	ClientID     string
	ClientSecret string
	InstanceType string
}

func (f FediverseApp) ClientConfig() *oauth2.Config {
	if f.MastodonCompatible() {
		return &oauth2.Config{
			ClientID:     f.ClientID,
			ClientSecret: f.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://" + f.Instance + "/oauth/authorize",
				TokenURL:  "https://" + f.Instance + "/oauth/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
			Scopes:      []string{"read:accounts"},
			RedirectURL: os.Getenv("BASE_URL") + "/auth/login/mastodon/" + f.Instance,
		}
	}

	return &oauth2.Config{
		ClientID:     f.ClientID,
		ClientSecret: f.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://" + f.Instance + "/auth",
			TokenURL:  "https://" + f.Instance + "/api/auth/session/oauth",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		Scopes:      []string{"read:account"},
		RedirectURL: os.Getenv("BASE_URL") + "/auth/login/misskey/" + f.Instance,
	}
}

func (f FediverseApp) MastodonCompatible() bool {
	return f.InstanceType == "mastodon" || f.InstanceType == "pleroma" || f.InstanceType == "akkoma" || f.InstanceType == "pixelfed"
}

const ErrNoInstanceApp = errors.Sentinel("instance doesn't have an app")

func (db *DB) FediverseApp(ctx context.Context, instance string) (fa FediverseApp, err error) {
	sql, args, err := sq.Select("*").From("fediverse_apps").Where("instance = ?", instance).ToSql()
	if err != nil {
		return fa, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &fa, sql, args...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return fa, ErrNoInstanceApp
		}
		return fa, errors.Wrap(err, "executing query")
	}
	return fa, nil
}

func (db *DB) FediverseAppByID(ctx context.Context, id int64) (fa FediverseApp, err error) {
	sql, args, err := sq.Select("*").From("fediverse_apps").Where("id = ?", id).ToSql()
	if err != nil {
		return fa, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Get(ctx, db, &fa, sql, args...)
	if err != nil {
		return fa, errors.Wrap(err, "executing query")
	}
	return fa, nil
}

func (db *DB) CreateFediverseApp(ctx context.Context, instance, instanceType, clientID, clientSecret string) (fa FediverseApp, err error) {
	sql, args, err := sq.Insert("fediverse_apps").
		Columns("instance", "instance_type", "client_id", "client_secret").
		Values(instance, instanceType, clientID, clientSecret).
		Suffix("RETURNING *").ToSql()
	if err != nil {
		return fa, errors.Wrap(err, "building query")
	}

	err = pgxscan.Get(ctx, db, &fa, sql, args...)
	if err != nil {
		return fa, errors.Wrap(err, "executing query")
	}
	return fa, nil
}
