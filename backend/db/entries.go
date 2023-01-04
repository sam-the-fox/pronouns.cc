package db

import "codeberg.org/u1f320/pronouns.cc/backend/db/queries"

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

type FieldEntry struct {
	Value  string     `json:"value"`
	Status WordStatus `json:"status"`
}

type PronounEntry struct {
	Pronouns    string     `json:"pronouns"`
	DisplayText *string    `json:"display_text"`
	Status      WordStatus `json:"status"`
}

func dbEntriesToFieldEntries(entries []queries.FieldEntry) []FieldEntry {
	out := make([]FieldEntry, len(entries))
	for i := range entries {
		out[i] = FieldEntry{
			*entries[i].Value, WordStatus(*entries[i].Status),
		}
	}
	return out
}

func dbPronounEntriesToPronounEntries(entries []queries.PronounEntry) []PronounEntry {
	out := make([]PronounEntry, len(entries))
	for i := range entries {
		out[i] = PronounEntry{
			*entries[i].Value, entries[i].DisplayValue, WordStatus(*entries[i].Status),
		}
	}
	return out
}

func entriesToDBEntries(entries []FieldEntry) []queries.FieldEntry {
	out := make([]queries.FieldEntry, len(entries))
	for i := range entries {
		status := int32(entries[i].Status)
		out[i] = queries.FieldEntry{
			&entries[i].Value, &status,
		}
	}
	return out
}

func pronounEntriesToDBEntries(entries []PronounEntry) []queries.PronounEntry {
	out := make([]queries.PronounEntry, len(entries))
	for i := range entries {
		status := int32(entries[i].Status)
		out[i] = queries.PronounEntry{
			&entries[i].Pronouns, entries[i].DisplayText, &status,
		}
	}
	return out
}
