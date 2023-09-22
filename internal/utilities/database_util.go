package utilities

import (
	"database/sql"
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host        string
	Username    string
	Password    string
	Port        string
	Database    string
	SSL         bool
	Certificate string
}

// GetMainDatabaseConnections connects to the main Database
func GetMainDatabaseConnection(config DatabaseConfig) (*sql.DB, error) {
	// Build and Establish a database connection
	sslConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full sslrootcert=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.Certificate,
	)
	normalConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
	)
	var connectionString string
	if config.SSL {
		connectionString = sslConnectionString
	} else {
		connectionString = normalConnectionString
	}
	var databaseConnection *sql.DB
	databaseConnection, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := databaseConnection.Ping(); err != nil {
		return nil, err
	}
	databaseConnection.SetMaxIdleConns(5)
	// Maximum Open Connections
	databaseConnection.SetMaxOpenConns(10)
	// Idle Connection Timeout
	databaseConnection.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	databaseConnection.SetConnMaxLifetime(30 * time.Second)
	return databaseConnection, nil
}
