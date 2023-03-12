package db

import (
	"fmt"
	"strings"
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

type FieldEntry struct {
	Value  string     `json:"value"`
	Status WordStatus `json:"status"`
}

func (fe FieldEntry) Validate() string {
	if fe.Value == "" {
		return "value cannot be empty"
	}

	if len([]rune(fe.Value)) > FieldEntryMaxLength {
		return fmt.Sprintf("name must be %d characters or less, is %d", FieldEntryMaxLength, len([]rune(fe.Value)))
	}

	if fe.Status == StatusUnknown || fe.Status >= wordStatusMax {
		return fmt.Sprintf("status is invalid, must be between 1 and %d, is %d", wordStatusMax-1, fe.Status)
	}

	return ""
}

type PronounEntry struct {
	Pronouns    string     `json:"pronouns"`
	DisplayText *string    `json:"display_text"`
	Status      WordStatus `json:"status"`
}

func (p PronounEntry) Validate() string {
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

func (p PronounEntry) String() string {
	if p.DisplayText != nil {
		return *p.DisplayText
	}

	split := strings.Split(p.Pronouns, "/")
	if len(split) <= 2 {
		return strings.Join(split, "/")
	}

	return strings.Join(split[:1], "/")
}
