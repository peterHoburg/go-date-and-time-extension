package dtegorm_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/peterHoburg/go-date-and-time-extension/dtegorm"
)

type Example struct {
	ID       uint `gorm:"primarykey"`
	OnlyTime dtegorm.Time
}

func Open(dsn string, db *gorm.DB) *gorm.DB {
	if db != nil {
		openDB, err := db.DB()
		if err != nil {
			log.Fatal("Error getting DB")
		}

		err = openDB.Close()
		if err != nil {
			log.Fatal("Error closing DB")
		}

		db = nil
	}

	var err error
	var localDB *gorm.DB

	retries := 3

	retry := 0

	for retry < retries {
		localDB, err = gorm.Open(postgres.Open(dsn))
		if err == nil {
			break
		}

		retry++

		time.Sleep(3 * time.Second) //nolint:mnd
	}

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	if localDB == nil {
		log.Fatal("Error connecting to database")
	}

	db = localDB.Debug()
	return db

}

func CreateDB(dbName string, db *gorm.DB) {
	result := db.Exec("CREATE DATABASE " + dbName)
	if result.Error != nil {
		log.Fatal("Error creating database")
	}
}

func DeleteDB(dbName string, db *gorm.DB) {
	result := db.Exec("DROP DATABASE " + dbName)
	if result.Error != nil {
		log.Fatal("Error deleting database")
	}
}

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&Example{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

}

func Setup() (string, string, *gorm.DB) {
	const DSN = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=GMT"
	var db *gorm.DB

	db = Open(DSN, db)

	dbName := strings.ReplaceAll("testdb-"+uuid.New().String(), "-", "")

	CreateDB(dbName, db)
	dsn := strings.ReplaceAll(DSN, "dbname=postgres", "dbname="+dbName)
	db = Open(dsn, db)
	RunMigrations(db)

	return dbName, dsn, db
}

func Teardown(dbName string, dsn string, db *gorm.DB) {
	db = Open(dsn, db)
	DeleteDB(dbName, db)
}

func ExampleTimeGORM() {
	dbName, dsn, db := Setup()
	dsn = strings.ReplaceAll(dsn, "dbname="+dbName, "dbname=postgres")
	defer Teardown(dbName, dsn, db)
	// ^^^ Setup for the PG DB. This can be ignored

	type Example struct {
		ID       uint `gorm:"primarykey"`
		OnlyTime dtegorm.Time
	}

	onlyTime, err := dtegorm.NewTime("10:04:05-05:00")
	if err != nil {
		return
	}

	example := Example{OnlyTime: onlyTime}

	createResult := db.Create(&example)
	if createResult.Error != nil {
		return
	}

	var exampleResult Example

	getResult := db.First(&exampleResult, example.ID)
	if getResult.Error != nil {
		return
	}
	fmt.Println(exampleResult.OnlyTime.String())

	// Output: 15:04:05Z
}

func TestTimeGORM(t *testing.T) {
	dbName, dsn, db := Setup()
	dsn = strings.ReplaceAll(dsn, "dbname="+dbName, "dbname=postgres")
	defer Teardown(dbName, dsn, db)

	type Result struct {
		ColumnName string
		DataType   string
	}
	result := Result{}

	db.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'examples' AND column_name = 'only_time'").Scan(&result)
	if result.ColumnName != "only_time" && result.DataType != "time with time zone" {
		t.Errorf("Column name or data type is not correct")
	}

	onlyTime, err := dtegorm.NewTime("10:04:05-05:00")
	if err != nil {
		t.Errorf("Error creating time")
	}

	example := Example{OnlyTime: onlyTime}
	dbResult := db.Create(&example)
	if dbResult.Error != nil {
		t.Errorf("Error creating example")
	}
	var exampleResult Example

	dbResult = db.First(&exampleResult, example.ID)
	if dbResult.Error != nil {
		t.Errorf("Error getting user")
	}

	if exampleResult.OnlyTime.String() != "15:04:05Z" {
		t.Errorf("Time is not correct")
	}
}