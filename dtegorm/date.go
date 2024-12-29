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
	ErrNewDate             = errors.New("failed to create new date")
	ErrDateScan            = errors.New("failed to scan value into date struct")
	ErrDateScanInvalidType = errors.New("invalid type passed to scan")
)

type Date struct { //nolint:recvcheck
	dte.Date
}

func NewDate(s string) (Date, error) {
	timeInstance := Date{}

	err := timeInstance.SetFromString(s)
	if err != nil {
		return Date{}, fmt.Errorf("%w: %w", ErrNewDate, err)
	}

	return timeInstance, nil
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (Date) GormDataType() string {
	return "date"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (Date) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	const date = "DATE"

	switch db.Dialector.Name() {
	case "mysql":
		return date
	case "postgres":
		return date
	case "sqlserver":
		return date
	case "sqlite":
		return date
	default:
		return ""
	}
}

// Scan implements sql.Scanner interface and scans value into Date,.
func (d *Date) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		err := d.SetFromString(string(v))
		if err != nil {
			return fmt.Errorf("%w: %w", ErrDateScan, err)
		}
	case string:
		err := d.SetFromString(v)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrDateScan, err)
		}
	case time.Time:
		err := d.SetFromTime(v)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrDateScan, err)
		}
	default:
		return ErrDateScanInvalidType
	}

	return nil
}

// Value implements driver.Valuer interface and returns string format of Date.
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}
