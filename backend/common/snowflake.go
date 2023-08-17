package common

import (
	"strconv"
	"strings"
	"time"
)

// Epoch is the pronouns.cc epoch (January 1st 2022 at 00:00:00 UTC) in milliseconds.
const Epoch = 1_640_995_200_000
const epochDuration = Epoch * time.Millisecond

const NullSnowflake = ^Snowflake(0)

// Snowflake is a 64-bit integer used as a unique ID, with an embedded timestamp.
type Snowflake uint64

// ID is an alias to Snowflake.
type ID = Snowflake

// ParseSnowflake parses a snowflake from a string.
func ParseSnowflake(sf string) (Snowflake, error) {
	if sf == "null" {
		return NullSnowflake, nil
	}

	i, err := strconv.ParseUint(sf, 10, 64)
	if err != nil {
		return 0, err
	}

	return Snowflake(i), nil
}

// NewSnowflake creates a new snowflake from the given time.
func NewSnowflake(t time.Time) Snowflake {
	ts := time.Duration(t.UnixNano()) - epochDuration

	return Snowflake((ts / time.Millisecond) << 22)
}

// String returns the snowflake as a string.
func (s Snowflake) String() string { return strconv.FormatUint(uint64(s), 10) }

// Time returns the creation time of the snowflake.
func (s Snowflake) Time() time.Time {
	ts := time.Duration(s>>22)*time.Millisecond + epochDuration
	return time.Unix(0, int64(ts))
}

func (s Snowflake) IsValid() bool {
	return s != 0 && s != NullSnowflake
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return []byte("null"), nil
	}

	return []byte(`"` + strconv.FormatUint(uint64(s), 10) + `"`), nil
}

func (s *Snowflake) UnmarshalJSON(src []byte) error {
	sf, err := ParseSnowflake(strings.Trim(string(src), `"`))
	if err != nil {
		return err
	}

	*s = sf
	return nil
}

func (s Snowflake) Worker() uint8 {
	return uint8(s & 0x3E0000 >> 17)
}

func (s Snowflake) PID() uint8 {
	return uint8(s & 0x1F000 >> 12)
}

func (s Snowflake) Increment() uint16 {
	return uint16(s & 0xFFF)
}
