package datatype

import (
	"database/sql/driver"
	"errors"
	"time"
)

func NewTime(value *time.Time) Time {
	return Time{value}
}

func NewTimeNow() Time {
	value := time.Now()
	return Time{value: &value}
}

func ParseTime(t string) (_ Time, err error) {
	tmp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return
	}
	return Time{value: &tmp}, err
}

type Time struct {
	value *time.Time
}

func (t Time) IsNil() bool {
	return t.value == nil
}

func (t Time) Time() *time.Time {
	return t.value
}

func (t Time) String() string {
	if t.value == nil {
		return ""
	}
	return time.Time(*t.value).Format(time.RFC3339)
}

func (t Time) MarshalText() ([]byte, error) {
	if t.value == nil {
		return []byte(""), nil
	}
	return []byte(time.Time(*t.value).Format(time.RFC3339)), nil
}

func (t *Time) UnmarshalText(b []byte) error {
	tmp, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		return err
	}
	t.value = &tmp

	return err
}

// Scan implements the Scanner interface.
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		t.value = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	tmp, err := time.Parse("2006-01-02 15:04:05", string(b))
	if err != nil {
		return err
	}
	t.value = &tmp
	return nil
}

// Value implements the driver Value interface.
func (t *Time) Value() (driver.Value, error) {
	if t.value == nil {
		return nil, nil
	}
	return time.Time(*t.value).Format("2006-01-02 15:04:05"), nil
}

func (t *Time) IsBefore(time Time) bool {
	if t.IsNil() {
		return false
	}
	return t.value.Before(*time.value)
}

func (t *Time) IsAfter(time Time) bool {
	if t.IsNil() {
		return false
	}
	return t.value.After(*time.value)
}
