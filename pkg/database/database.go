package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbName     = "DB_NAME"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbConn     = "DB_CONNECTION_STRING"
)

type DatabaseConnInfo struct {
	Host          string
	Port          string
	DbName        string
	User          string
	Password      string
	ConnectionStr string
}

func newDatabaseConnInfo() *DatabaseConnInfo {
	return &DatabaseConnInfo{
		Host:          os.Getenv(dbHost),
		Port:          os.Getenv(dbPort),
		DbName:        os.Getenv(dbName),
		User:          os.Getenv(dbUser),
		Password:      os.Getenv(dbPassword),
		ConnectionStr: os.Getenv(dbConn),
	}
}

func (d *DatabaseConnInfo) validate() error {
	if dbConn != "" {
		return nil
	}

	if d.Host == "" {
		return fmt.Errorf("%s is required", dbHost)
	}

	if d.Port == "" {
		return fmt.Errorf("%s is required", dbPort)
	}

	if d.DbName == "" {
		return fmt.Errorf("%s is required", dbName)
	}

	if d.User == "" {
		return fmt.Errorf("%s is required", dbUser)
	}

	if d.Password == "" {
		return fmt.Errorf("%s is required", dbPassword)
	}

	return nil
}

func (d *DatabaseConnInfo) getConnectionString() string {
	if d.ConnectionStr != "" {
		return d.ConnectionStr
	}

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", d.Host, d.Port, d.User, d.DbName, d.Password)
}

func CreateDatabaseConnection() (*gorm.DB, error) {
	connInfo := newDatabaseConnInfo()
	if err := connInfo.validate(); err != nil {
		return nil, fmt.Errorf("error validating database connection info: %v", err)
	}

	db, err := gorm.Open(postgres.Open(connInfo.getConnectionString()), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}
