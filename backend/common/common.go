// Package common contains functions and types common to all (or most) packages.
package common

import "unicode/utf8"

func StringLength(s *string) int {
	if s == nil {
		return -1
	}
	return utf8.RuneCountInString(*s)
}
