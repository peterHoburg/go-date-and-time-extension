package dtegorm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/peterHoburg/go-date-and-time-extension/dtegorm"
)

type DateExample struct {
	ID       uint `gorm:"primarykey"`
	OnlyDate dtegorm.Date
}

func ExampleDate() {
	dbName, dsn, db := Setup()
	dsn = strings.ReplaceAll(dsn, "dbname="+dbName, "dbname=postgres")
	defer Teardown(dbName, dsn, db)
	// ^^^ Setup for the PG DB. This can be ignored

	onlyDate, err := dtegorm.NewDate("2006-01-02")
	if err != nil {
		return
	}

	example := DateExample{OnlyDate: onlyDate}

	createResult := db.Create(&example)
	if createResult.Error != nil {
		return
	}

	var exampleResult DateExample

	getResult := db.First(&exampleResult, example.ID)
	if getResult.Error != nil {
		return
	}

	fmt.Println(exampleResult.OnlyDate.String())

	// Output: 2006-01-02
}

func TestDate(t *testing.T) {
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
			"WHERE table_name = 'date_examples' AND column_name = 'only_date'",
	).Scan(&result)

	if result.ColumnName != "only_date" && result.DataType != "date" {
		t.Errorf("Column name or data type is not correct")
	}

	onlyDate, err := dtegorm.NewDate("2006-01-02T15:04:05Z")
	if err != nil {
		t.Errorf("Error creating date")
	}

	example := DateExample{OnlyDate: onlyDate}

	dbResult := db.Create(&example)
	if dbResult.Error != nil {
		t.Errorf("Error creating example")
	}

	var exampleResult DateExample

	dbResult = db.First(&exampleResult, example.ID)
	if dbResult.Error != nil {
		t.Errorf("Error getting user")
	}

	if exampleResult.OnlyDate.String() != "2006-01-02" {
		t.Errorf("Date is not correct, %s, %s", exampleResult.OnlyDate.String(), "2006-01-02")
	}
}
