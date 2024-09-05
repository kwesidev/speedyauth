package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kwesidev/speedyauth/internal/apiserver"
	"github.com/kwesidev/speedyauth/internal/utilities"
	_ "github.com/lib/pq"
)

var databaseConnection *sql.DB

const SERVER_VERSION = "0.1.1"

// initializes the database connections and other connections
func init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Not loading Config from .env")
	}

	databaseConfig := utilities.DatabaseConfig{
		Host:     os.Getenv("PG_HOST"),
		Username: os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Port:     os.Getenv("PG_PORT"),
		Database: os.Getenv("PG_DB"),
	}
	if strings.ToTitle(os.Getenv("PG_SSL")) == "True" {
		databaseConfig.SSL = true
		databaseConfig.Certificate = os.Getenv("PG_CERT")
	} else {
		databaseConfig.SSL = false
	}
	databaseConnection, err = utilities.GetMainDatabaseConnection(databaseConfig)
	if err != nil {
		log.Fatal("Failed to Connect to the  Database", err)
	}
	// // Apply DB migrations
	// driver, err := postgres.WithInstance(databaseConnection, &postgres.Config{})
	// m, err := migrate.NewWithDatabaseInstance(
	// 	"file://db/migrations",
	// 	"postgres", driver)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := m.Up(); err != nil {
	// 	log.Fatal(err)
	// }
	log.Println("Server Version :", SERVER_VERSION)
	asciiArt := `

███████ ██████  ███████ ███████ ██████  ██    ██      █████  ██    ██ ████████ ██   ██ 
██      ██   ██ ██      ██      ██   ██  ██  ██      ██   ██ ██    ██    ██    ██   ██ 
███████ ██████  █████   █████   ██   ██   ████       ███████ ██    ██    ██    ███████ 
     ██ ██      ██      ██      ██   ██    ██        ██   ██ ██    ██    ██    ██   ██ 
███████ ██      ███████ ███████ ██████     ██        ██   ██  ██████     ██    ██   ██ 
                                                                                       
																																		   																			
	`
	log.Println(asciiArt)

	log.Println("")
}
func main() {
	apiServer := apiserver.NewAPIServer(os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"), databaseConnection)
	apiServer.Run()
}
