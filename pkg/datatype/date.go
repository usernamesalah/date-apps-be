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

func (t *Date) IsNil() bool {
	return t == nil || t.value == nil
}

func (t *Date) Time() *time.Time {
	if t == nil {
		return nil
	}
	return t.value
}

func (t *Date) AddDate(years, months, days int) Date {
	if t == nil || t.value == nil {
		return Date{}
	}
	addedTime := t.value.AddDate(years, months, days)
	return Date{value: &addedTime}
}

func (t *Date) String() string {
	if t == nil || t.value == nil {
		return ""
	}

	return t.value.Format(time.RFC3339)
}

func (t *Date) MarshalText() ([]byte, error) {
	if t == nil || t.value == nil {
		return []byte(""), nil
	}

	return []byte(t.value.Format(dateFormat)), nil
}

func (t *Date) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		t.value = nil
		return nil
	}
	tmp, err := time.Parse(dateFormat, string(b))
	if err != nil {
		return err
	}
	t.value = &tmp

	return nil
}

// Scan implements the Scanner interface.
func (t *Date) Scan(value interface{}) error {
	if value == nil {
		t.value = nil
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.value = &v
	case []byte:
		tmp, err := time.Parse(dateFormat, string(v))
		if err != nil {
			return err
		}
		t.value = &tmp
	default:
		return errors.New("datatype.Date: unsupported scan source type")
	}

	return nil
}

// Value implements the driver Valuer interface.
func (t *Date) Value() (driver.Value, error) {
	if t == nil || t.value == nil {
		return nil, nil
	}
	return t.value.Format(dateFormat), nil
}

func (t *Date) IsBefore(d Date) bool {
	if t == nil || t.value == nil || d.value == nil {
		return false
	}
	return t.value.Before(*d.value)
}

func (t *Date) IsAfter(d Date) bool {
	if t == nil || t.value == nil || d.value == nil {
		return false
	}
	return t.value.After(*d.value)
}
