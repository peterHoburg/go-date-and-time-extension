package dte // Package dte import github.com/peterHoburg/go-date-and-time-extension/dte

import (
	"errors"
	"fmt"
	"time"
)

var ErrTimeParse = errors.New("time does not follow 15:04:05 time only format")

const (
	TimeOnlyWithTimezone          = "15:04:05Z07:00"
	TimeOnlyWithTimezoneWithSpace = "15:04:05 Z07:00"
	TimeOnlyWithTimezoneShort     = "15:04:05Z07"
)

var timeAcceptableFormats = []string{ //nolint:gochecknoglobals
	TimeOnlyWithTimezone,
	TimeOnlyWithTimezoneWithSpace,
	TimeOnlyWithTimezoneShort,
}

type Time struct { //nolint:recvcheck
	time.Time `example:"15:04:05Z" format:"time"`
}

func NewTime(s string) (Time, error) {
	timeInstance := Time{}

	err := timeInstance.SetFromString(s)
	if err != nil {
		return Time{}, err
	}

	return timeInstance, nil
}

func (t *Time) SetFromString(s string) error {
	var err error

	parsedTime := time.Time{}

	for _, layout := range timeAcceptableFormats {
		if parsedTime, err = time.Parse(layout, s); err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrTimeParse, err)
	}

	parsedTime = parsedTime.UTC()
	*t = Time{parsedTime}

	return nil
}

func (t *Time) SetFromTime(inputTime time.Time) error {
	timeSting := inputTime.Format(TimeOnlyWithTimezone)

	err := t.SetFromString(timeSting)
	if err != nil {
		return err
	}

	return nil
}

func (t Time) String() string {
	return t.Format(TimeOnlyWithTimezone)
}

// MarshalJSON implements the [json.Marshaler] interface.
// The time is a quoted string in the hh:mm:ss format.
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeOnlyWithTimezone)+len(`""`))

	b = append(b, '"')

	formatedTime := t.Format(TimeOnlyWithTimezone)

	b = append(b, formatedTime...)
	b = append(b, '"')

	return b, nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
// The time must be a quoted string in the RFC 3339 format or hh:mm:ss.
func (t *Time) UnmarshalJSON(data []byte) error {
	tempTime := time.Time{}

	var parsedTime Time

	err := tempTime.UnmarshalJSON(data)
	if err != nil { //nolint:nestif
		if string(data) == "null" || string(data) == "\"null\"" {
			return nil
		}

		if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
			return fmt.Errorf("Time.UnmarshalJSON: input is not a JSON string: %w", err)
		}

		data = data[len(`"`) : len(data)-len(`"`)]

		parsedTime, err = NewTime(string(data))
		if err != nil {
			return err
		}
	} else {
		parsedTime, err = NewTime(tempTime.Format(TimeOnlyWithTimezone))
		if err != nil {
			return fmt.Errorf("failed to format unmarshaled time: %w", err)
		}
	}

	*t = parsedTime

	return nil
}
