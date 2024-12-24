package dtegorm_test

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dtegorm"
)

type User struct {
	ID       uint `gorm:"primarykey"`
	OnlyTime dtegorm.Time
}

//nolint:funlen
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

func Close(db *gorm.DB) {
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
	err := db.AutoMigrate(&User{})
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

func TestTimeGorm(t *testing.T) {
	dbName, dsn, db := Setup()
	dsn = strings.ReplaceAll(dsn, "dbname="+dbName, "dbname=postgres")
	defer Teardown(dbName, dsn, db)
}
