package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/rs/xid"
)

type UserSettings struct {
	ReadChangelog      string `json:"read_changelog"`
	ReadSettingsNotice string `json:"read_settings_notice"`
}

func (db *DB) UpdateUserSettings(ctx context.Context, id xid.ID, us UserSettings) error {
	sql, args, err := sq.Update("users").Set("settings", us).Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}
	return nil
}
