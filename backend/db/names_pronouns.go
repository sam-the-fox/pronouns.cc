package db

import (
	"context"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type WordStatus int

const (
	StatusUnknown     WordStatus = 0
	StatusFavourite   WordStatus = 1
	StatusOkay        WordStatus = 2
	StatusJokingly    WordStatus = 3
	StatusFriendsOnly WordStatus = 4
	StatusAvoid       WordStatus = 5
	wordStatusMax     WordStatus = 6
)

type Name struct {
	ID     int64      `json:"-"`
	Name   string     `json:"name"`
	Status WordStatus `json:"status"`
}

func (n Name) Validate() string {
	if n.Name == "" {
		return "name cannot be empty"
	}

	if len([]rune(n.Name)) > FieldEntryMaxLength {
		return fmt.Sprintf("name must be %d characters or less, is %d", FieldEntryMaxLength, len([]rune(n.Name)))
	}

	if n.Status == StatusUnknown || n.Status >= wordStatusMax {
		return fmt.Sprintf("status is invalid, must be between 1 and %d, is %d", wordStatusMax-1, n.Status)
	}

	return ""
}

type Pronoun struct {
	ID          int64      `json:"-"`
	DisplayText *string    `json:"display_text"`
	Pronouns    string     `json:"pronouns"`
	Status      WordStatus `json:"status"`
}

func (p Pronoun) Validate() string {
	if p.Pronouns == "" {
		return "pronouns cannot be empty"
	}

	if p.DisplayText != nil {
		if len([]rune(*p.DisplayText)) > FieldEntryMaxLength {
			return fmt.Sprintf("display_text must be %d characters or less, is %d", FieldEntryMaxLength, len([]rune(*p.DisplayText)))
		}
	}

	if len([]rune(p.Pronouns)) > FieldEntryMaxLength {
		return fmt.Sprintf("pronouns must be %d characters or less, is %d", FieldEntryMaxLength, len([]rune(p.Pronouns)))
	}

	if p.Status == StatusUnknown || p.Status >= wordStatusMax {
		return fmt.Sprintf("status is invalid, must be between 1 and %d, is %d", wordStatusMax-1, p.Status)
	}

	return ""
}

func (p Pronoun) String() string {
	if p.DisplayText != nil {
		return *p.DisplayText
	}

	split := strings.Split(p.Pronouns, "/")
	if len(split) <= 2 {
		return strings.Join(split, "/")
	}

	return strings.Join(split[:1], "/")
}

func (db *DB) UserNames(ctx context.Context, userID xid.ID) (ns []Name, err error) {
	sql, args, err := sq.Select("id", "name", "status").From("user_names").Where("user_id = ?", userID).OrderBy("id").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &ns, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	return ns, nil
}

func (db *DB) UserPronouns(ctx context.Context, userID xid.ID) (ps []Pronoun, err error) {
	sql, args, err := sq.
		Select("id", "display_text", "pronouns", "status").
		From("user_pronouns").Where("user_id = ?", userID).
		OrderBy("id").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building sql")
	}

	err = pgxscan.Select(ctx, db, &ps, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing query")
	}
	return ps, nil
}

func (db *DB) SetUserNames(ctx context.Context, tx pgx.Tx, userID xid.ID, names []Name) (err error) {
	sql, args, err := sq.Delete("user_names").Where("user_id = ?", userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting existing names")
	}

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"user_names"},
		[]string{"user_id", "name", "status"},
		pgx.CopyFromSlice(len(names), func(i int) ([]any, error) {
			return []any{
				userID,
				names[i].Name,
				names[i].Status,
			}, nil
		}))
	if err != nil {
		return errors.Wrap(err, "inserting new names")
	}
	return nil
}

func (db *DB) SetUserPronouns(ctx context.Context, tx pgx.Tx, userID xid.ID, names []Pronoun) (err error) {
	sql, args, err := sq.Delete("user_pronouns").Where("user_id = ?", userID).ToSql()
	if err != nil {
		return errors.Wrap(err, "building sql")
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting existing pronouns")
	}

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"user_pronouns"},
		[]string{"user_id", "pronouns", "display_text", "status"},
		pgx.CopyFromSlice(len(names), func(i int) ([]any, error) {
			return []any{
				userID,
				names[i].Pronouns,
				names[i].DisplayText,
				names[i].Status,
			}, nil
		}))
	if err != nil {
		return errors.Wrap(err, "inserting new pronouns")
	}
	return nil
}
