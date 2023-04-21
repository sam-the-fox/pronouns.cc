package db

import (
	"fmt"
	"strings"
)

type WordStatus string

func (w *WordStatus) UnmarshalJSON(src []byte) error {
	if string(src) == "null" {
		return nil
	}

	s := strings.Trim(string(src), `"`)
	switch s {
	case "1":
		*w = "favourite"
	case "2":
		*w = "okay"
	case "3":
		*w = "jokingly"
	case "4":
		*w = "friends_only"
	case "5":
		*w = "avoid"
	default:
		*w = WordStatus(s)
	}

	return nil
}

func (w WordStatus) Valid(extra CustomPreferences) bool {
	if w == "favourite" || w == "okay" || w == "jokingly" || w == "friends_only" || w == "avoid" {
		return true
	}

	for k := range extra {
		if string(w) == k {
			return true
		}
	}
	return false
}

type FieldEntry struct {
	Value  string     `json:"value"`
	Status WordStatus `json:"status"`
}

func (fe FieldEntry) Validate(custom CustomPreferences) string {
	if fe.Value == "" {
		return "value cannot be empty"
	}

	if len([]rune(fe.Value)) > FieldEntryMaxLength {
		return fmt.Sprintf("name must be %d characters or less, is %d", FieldEntryMaxLength, len([]rune(fe.Value)))
	}

	if !fe.Status.Valid(custom) {
		return "status is invalid"
	}

	return ""
}

type PronounEntry struct {
	Pronouns    string     `json:"pronouns"`
	DisplayText *string    `json:"display_text"`
	Status      WordStatus `json:"status"`
}

func (p PronounEntry) Validate(custom CustomPreferences) string {
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

	if !p.Status.Valid(custom) {
		return "status is invalid"
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
