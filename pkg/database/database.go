package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbConn = "DB_CONNECTION_STRING"
)

type DatabaseConnInfo struct {
	ConnectionStr string
}

func newDatabaseConnInfo() *DatabaseConnInfo {
	return &DatabaseConnInfo{
		ConnectionStr: os.Getenv(dbConn),
	}
}

func (d *DatabaseConnInfo) validate() error {
	if d.ConnectionStr == "" {
		return fmt.Errorf("database connection string not found")
	}
	return nil
}

func CreateDatabaseConnection() (*gorm.DB, error) {
	connInfo := newDatabaseConnInfo()
	if err := connInfo.validate(); err != nil {
		return nil, fmt.Errorf("error validating database connection info: %v", err)
	}

	db, err := gorm.Open(postgres.Open(connInfo.ConnectionStr), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}
