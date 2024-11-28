package datatype

import (
	"database/sql/driver"
	"date-apps-be/pkg/derrors"
	"errors"
	"time"
)

const dateFormat = "2006-01-02"

func ParseDate(str string, location string) (Date, error) {

	loc, err := time.LoadLocation(location)
	if err != nil {
		return Date{}, derrors.New(derrors.InvalidArgument, "Invalid Time Location")
	}
	tmp, err := time.ParseInLocation(dateFormat, str, loc)

	return Date{
		value: &tmp,
	}, err
}

func NewDateNow() Date {
	value := time.Now()
	return Date{
		value: &value,
	}
}

type Date struct {
	value *time.Time
}

func (t Date) IsNil() bool {
	return t.value == nil
}

func (t Date) Time() *time.Time {
	return t.value
}

func (t Date) AddDate(years, month, day int) Date {
	if t.value == nil {
		return Date{}
	}
	addedTime := t.value.AddDate(years, month, day)
	t.value = &addedTime
	return t
}

func (t Date) String() string {
	if t.value == nil {
		return ""
	}

	return time.Time(*t.value).Format(time.RFC3339)
}

func (t Date) MarshalText() ([]byte, error) {
	if t.value == nil {
		return []byte(""), nil
	}

	return []byte(time.Time(*t.value).Format(dateFormat)), nil
}

func (t *Date) UnmarshalText(b []byte) error {
	tmp, err := time.Parse(dateFormat, string(b))
	if err != nil {
		return err
	}
	t.value = &tmp

	return err
}

// Scan implements the Scanner interface.
func (t *Date) Scan(value interface{}) error {
	if value == nil {
		t.value = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	tmp, err := time.Parse(dateFormat, string(b))
	if err != nil {
		return err
	}
	t.value = &tmp
	return nil
}

// Value implements the driver Value interface.
func (t *Date) Value() (driver.Value, error) {
	if t.value == nil {
		return nil, nil
	}
	return time.Time(*t.value).Format(dateFormat), nil
}

func (t *Date) IsBefore(d Date) bool {
	if t.IsNil() {
		return false
	}
	return t.value.Before(*d.value)
}

func (t *Date) IsAfter(d Date) bool {
	if t.IsNil() {
		return false
	}
	return t.value.After(*d.value)
}
