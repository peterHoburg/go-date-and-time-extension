package dtetime // Package dtetime import github.com/peterHoburg/go-date-and-time-extension/dtetime

import (
	"errors"
	"fmt"
	"time"
)

var ErrTimeParse = errors.New("time does not follow 15:04:05 time only format")

const TimeOnlyWithTimezone = "15:04:05Z07:00"

type Time struct {
	time.Time
}

func Parse(s string) (Time, error) {
	parsedTime, err := time.Parse(TimeOnlyWithTimezone, s)
	if err != nil {
		return Time{}, fmt.Errorf("%w: %w", ErrTimeParse, err)
	}
	parsedTime = parsedTime.UTC()

	return Time{Time: parsedTime}, nil
}

func (t Time) String() string {
	return t.Format(TimeOnlyWithTimezone)
}

func (t Time) UnmarshalText(text []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t Time) MarshalText() (text []byte, err error) {
	//TODO implement me
	panic("implement me")
}

// MarshalJSON implements the [json.Marshaler] interface.
// The time is a quoted string in the hh:mm:ss format
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
	parsedTime := Time{}

	err := tempTime.UnmarshalJSON(data)
	if err != nil {

		if string(data) == "null" {
			return nil
		}

		if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
			return errors.New("Time.UnmarshalJSON: input is not a JSON string")
		}

		data = data[len(`"`) : len(data)-len(`"`)]

		parsedTime, err = Parse(string(data))
		if err != nil {
			return err
		}

	} else {

		parsedTime, err = Parse(tempTime.Format(TimeOnlyWithTimezone))
		if err != nil {
			return fmt.Errorf("failed to format unmarshaled time: %w", err)
		}

	}

	*t = parsedTime
	return nil
}
