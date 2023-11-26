package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // глобальная переменная для доступа к базе данных
var err error

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbTimezone := os.Getenv("DB_TIMEZONE")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUser, dbPassword, dbTimezone)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	var dbNameOnServer string
	queryResult := DB.Raw("SELECT datname FROM pg_database WHERE datname = ?", dbName).Scan(&dbNameOnServer)
	if queryResult.Error != nil {
		log.Panicf("Error executing query: %v", queryResult.Error)
	} else if queryResult.RowsAffected == 0 {
		createDatabaseCommand := fmt.Sprintf("CREATE DATABASE \"%s\"", dbName)
		DB.Exec(createDatabaseCommand)
	}

	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbTimezone)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
}
