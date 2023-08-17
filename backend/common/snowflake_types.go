package common

import "time"

type UserID Snowflake

func (id UserID) String() string    { return Snowflake(id).String() }
func (id UserID) Time() time.Time   { return Snowflake(id).Time() }
func (id UserID) IsValid() bool     { return Snowflake(id).IsValid() }
func (id UserID) Worker() uint8     { return Snowflake(id).Worker() }
func (id UserID) PID() uint8        { return Snowflake(id).PID() }
func (id UserID) Increment() uint16 { return Snowflake(id).Increment() }

func (id UserID) MarshalJSON() ([]byte, error)    { return Snowflake(id).MarshalJSON() }
func (id *UserID) UnmarshalJSON(src []byte) error { return (*Snowflake)(id).UnmarshalJSON(src) }

type MemberID Snowflake

func (id MemberID) String() string    { return Snowflake(id).String() }
func (id MemberID) Time() time.Time   { return Snowflake(id).Time() }
func (id MemberID) IsValid() bool     { return Snowflake(id).IsValid() }
func (id MemberID) Worker() uint8     { return Snowflake(id).Worker() }
func (id MemberID) PID() uint8        { return Snowflake(id).PID() }
func (id MemberID) Increment() uint16 { return Snowflake(id).Increment() }

func (id MemberID) MarshalJSON() ([]byte, error)    { return Snowflake(id).MarshalJSON() }
func (id *MemberID) UnmarshalJSON(src []byte) error { return (*Snowflake)(id).UnmarshalJSON(src) }

type FlagID Snowflake

func (id FlagID) String() string    { return Snowflake(id).String() }
func (id FlagID) Time() time.Time   { return Snowflake(id).Time() }
func (id FlagID) IsValid() bool     { return Snowflake(id).IsValid() }
func (id FlagID) Worker() uint8     { return Snowflake(id).Worker() }
func (id FlagID) PID() uint8        { return Snowflake(id).PID() }
func (id FlagID) Increment() uint16 { return Snowflake(id).Increment() }

func (id FlagID) MarshalJSON() ([]byte, error)    { return Snowflake(id).MarshalJSON() }
func (id *FlagID) UnmarshalJSON(src []byte) error { return (*Snowflake)(id).UnmarshalJSON(src) }
