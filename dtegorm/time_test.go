package dtegorm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/peterHoburg/go-date-and-time-extension/dtegorm"
)

type TimeExample struct {
	ID       uint `gorm:"primarykey"`
	OnlyTime dtegorm.Time
}

func ExampleTime() {
	dbName, dsn, db := Setup()
	dsn = strings.ReplaceAll(dsn, "dbname="+dbName, "dbname=postgres")
	defer Teardown(dbName, dsn, db)
	// ^^^ Setup for the PG DB. This can be ignored

	onlyTime, err := dtegorm.NewTime("10:04:05-05:00")
	if err != nil {
		return
	}

	example := TimeExample{OnlyTime: onlyTime}

	createResult := db.Create(&example)
	if createResult.Error != nil {
		return
	}

	var exampleResult TimeExample

	getResult := db.First(&exampleResult, example.ID)
	if getResult.Error != nil {
		return
	}

	fmt.Println(exampleResult.OnlyTime.String())

	// Output: 15:04:05Z
}

func TestTime(t *testing.T) {
	t.Parallel()

	dbName, dsn, db := Setup()
	dsn = strings.ReplaceAll(dsn, "dbname="+dbName, "dbname=postgres")
	defer Teardown(dbName, dsn, db)

	type Result struct {
		ColumnName string
		DataType   string
	}

	result := Result{}

	db.Raw(
		"SELECT column_name, data_type " +
			"FROM information_schema.columns " +
			"WHERE table_name = 'time_examples' AND column_name = 'only_time'",
	).Scan(&result)

	if result.ColumnName != "only_time" && result.DataType != "time with time zone" {
		t.Errorf("Column name or data type is not correct")
	}

	onlyTime, err := dtegorm.NewTime("10:04:05-05:00")
	if err != nil {
		t.Errorf("Error creating time")
	}

	example := TimeExample{OnlyTime: onlyTime}

	dbResult := db.Create(&example)
	if dbResult.Error != nil {
		t.Errorf("Error creating example")
	}

	var exampleResult TimeExample

	dbResult = db.First(&exampleResult, example.ID)
	if dbResult.Error != nil {
		t.Errorf("Error getting user")
	}

	if exampleResult.OnlyTime.String() != "15:04:05Z" {
		t.Errorf("Time is not correct")
	}
}
