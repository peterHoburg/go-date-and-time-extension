package dte // Package dte import github.com/peterHoburg/go-date-and-time-extension/dte

import (
	"errors"
	"fmt"
	"time"
)

var ErrDateParse = errors.New("date does not follow yyyy-mm-dd date only format")

const (
	DateOnly = "2006-01-02"
)

var dateAcceptableFormats = []string{ //nolint:gochecknoglobals
	DateOnly,
	time.RFC3339,
}

type Date struct { //nolint:recvcheck
	time.Time `example:"2006-01-02"`
}

func NewDate(s string) (Date, error) {
	timeInstance := Date{}

	err := timeInstance.SetFromString(s)
	if err != nil {
		return Date{}, err
	}

	return timeInstance, nil
}

func (d *Date) SetFromString(s string) error {
	var err error

	parsedTime := time.Time{}

	for _, layout := range dateAcceptableFormats {
		if parsedTime, err = time.Parse(layout, s); err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrDateParse, err)
	}

	*d = Date{parsedTime}

	return nil
}

func (d *Date) SetFromTime(inputTime time.Time) error {
	timeSting := inputTime.Format(DateOnly)

	err := d.SetFromString(timeSting)
	if err != nil {
		return err
	}

	return nil
}

func (d Date) String() string {
	return d.Format(DateOnly)
}

// MarshalJSON implements the [json.Marshaler] interface.
// The date is a quoted string in the yyyy-mm-dd format.
func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateOnly)+len(`""`))

	b = append(b, '"')

	formatedTime := d.Format(DateOnly)

	b = append(b, formatedTime...)
	b = append(b, '"')

	return b, nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
// The time must be a quoted string in the RFC 3339 format or yyyy-mm-dd.
func (d *Date) UnmarshalJSON(data []byte) error {
	tempDate := time.Time{}

	var parsedTime Date

	err := tempDate.UnmarshalJSON(data)
	if err != nil { //nolint:nestif
		if string(data) == "null" || string(data) == "\"null\"" {
			return nil
		}

		if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
			return fmt.Errorf("Date.UnmarshalJSON: input is not a JSON string: %w", err)
		}

		data = data[len(`"`) : len(data)-len(`"`)]

		parsedTime, err = NewDate(string(data))
		if err != nil {
			return err
		}
	} else {
		parsedTime, err = NewDate(tempDate.Format(DateOnly))
		if err != nil {
			return fmt.Errorf("failed to format unmarshaled time: %w", err)
		}
	}

	*d = parsedTime

	return nil
}
