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

type Time struct {
	dte.Time
}

func NewTime(s string) (Time, error) {
	timeInstance := Time{}

	err := timeInstance.SetFromString(s)
	if err != nil {
		return Time{}, err
	}
	return timeInstance, nil
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (Time) GormDataType() string {
	return "time"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
			return err
		}
	case string:
		err := t.SetFromString(v)
		if err != nil {
			return err
		}
	case time.Time:
		err := t.SetFromTime(v)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("failed to scan value: %v", v))
	}

	return nil
}

// Value implements driver.Valuer interface and returns string format of Time.
func (t Time) Value() (driver.Value, error) {
	return t.String(), nil
}
