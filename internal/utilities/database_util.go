package utilities

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// GetMainDatabaseConnections connects to the main Database
func GetMainDatabaseConnection() *sql.DB {
	// Build and Establish a database connection
	sslConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full sslrootcert=%s",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB"),
		os.Getenv("PG_CERT"),
	)
	normalConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB"),
	)
	var connectionString string
	if strings.ToTitle(os.Getenv("PG_SSL")) == "True" {
		connectionString = sslConnectionString
	} else {
		connectionString = normalConnectionString
	}
	var databaseConnection *sql.DB
	databaseConnection, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	if err := databaseConnection.Ping(); err != nil {
		log.Fatalln(err)
	}
	databaseConnection.SetMaxIdleConns(5)
	// Maximum Open Connections
	databaseConnection.SetMaxOpenConns(10)
	// Idle Connection Timeout
	databaseConnection.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	databaseConnection.SetConnMaxLifetime(30 * time.Second)
	return databaseConnection
}
