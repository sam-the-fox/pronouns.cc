package db

import (
	"context"

	"codeberg.org/u1f320/pronouns.cc/backend/db/queries"
	"emperror.dev/errors"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

func (db *DB) SetUserNamesPronouns(ctx context.Context, tx pgx.Tx, userID xid.ID, names []FieldEntry, pronouns []PronounEntry) (err error) {
	_, err = queries.NewQuerier(tx).UpdateUserNamesPronouns(ctx, queries.UpdateUserNamesPronounsParams{
		ID:       userID.String(),
		Names:    entriesToDBEntries(names),
		Pronouns: pronounEntriesToDBEntries(pronouns),
	})
	if err != nil {
		return errors.Wrap(err, "executing update names/pronouns query")
	}
	return nil
}

func (db *DB) SetMemberNamesPronouns(ctx context.Context, tx pgx.Tx, memberID xid.ID, names []FieldEntry, pronouns []PronounEntry) (err error) {
	_, err = queries.NewQuerier(tx).UpdateMemberNamesPronouns(ctx, queries.UpdateMemberNamesPronounsParams{
		ID:       memberID.String(),
		Names:    entriesToDBEntries(names),
		Pronouns: pronounEntriesToDBEntries(pronouns),
	})
	if err != nil {
		return errors.Wrap(err, "executing update names/pronouns query")
	}
	return nil
}

func fieldEntriesFromDB(dn []queries.FieldEntry) []FieldEntry {
	names := make([]FieldEntry, len(dn))
	for i := range dn {
		names[i] = FieldEntry{
			Value:  *dn[i].Value,
			Status: WordStatus(*dn[i].Status),
		}
	}
	return names
}

func pronounsFromDB(dn []queries.PronounEntry) []PronounEntry {
	pronouns := make([]PronounEntry, len(dn))
	for i := range dn {
		pronouns[i] = PronounEntry{
			DisplayText: dn[i].DisplayValue,
			Pronouns:    *dn[i].Value,
			Status:      WordStatus(*dn[i].Status),
		}
	}
	return pronouns
}
