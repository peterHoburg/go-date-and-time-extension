package dtegorm

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/peterHoburg/go-date-and-time-extension/dte"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	ErrNewTime             = errors.New("failed to create new time")
	ErrTimeScan            = errors.New("failed to scan value into time struct")
	ErrTimeScanInvalidType = errors.New("invalid type passed to scan")
)

type Time struct { //nolint:recvcheck
	dte.Time `example:"15:04:05Z" format:"time"`
}

func NewTime(s string) (Time, error) {
	timeInstance := Time{}

	err := timeInstance.SetFromString(s)
	if err != nil {
		return Time{}, fmt.Errorf("%w: %w", ErrNewTime, err)
	}

	return timeInstance, nil
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (Time) GormDataType() string {
	return "time"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (Time) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "TIME"
	case "postgres":
		return "TIME with time zone"
	case "sqlserver":
		return "TIME"
	case "sqlite":
		return "TEXT"
	default:
		return ""
	}
}

// Scan implements sql.Scanner interface and scans value into Time,.
func (t *Time) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		err := t.SetFromString(string(v))
		if err != nil {
			return fmt.Errorf("%w: %w", ErrTimeScan, err)
		}
	case string:
		err := t.SetFromString(v)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrTimeScan, err)
		}
	case time.Time:
		err := t.SetFromTime(v)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrTimeScan, err)
		}
	default:
		return ErrTimeScanInvalidType
	}

	return nil
}

// Value implements driver.Valuer interface and returns string format of Time.
func (t Time) Value() (driver.Value, error) {
	return t.String(), nil
}
