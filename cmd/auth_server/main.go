package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kwesidev/authserver/internal/apiserver"
	"github.com/kwesidev/authserver/internal/utilities"
)

var databaseConnection *sql.DB

// initializes the database connections and other connections
func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Not loading Config from .env")
	}

	databaseConfig := utilities.DatabaseConfig{
		Host:     os.Getenv("PG_HOST"),
		Userame:  os.Getenv("PG_USER"),
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

	asciiArt := `

	█████╗ ██╗   ██╗████████╗██╗  ██╗    ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗     
	██╔══██╗██║   ██║╚══██╔══╝██║  ██║    ██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗    
	███████║██║   ██║   ██║   ███████║    ███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝    
	██╔══██║██║   ██║   ██║   ██╔══██║    ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗    
	██║  ██║╚██████╔╝   ██║   ██║  ██║    ███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║    
	╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝    ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝    
																							   
	
	`
	log.Println(asciiArt)
}
func main() {
	initialize()
	apiServer := apiserver.NewAPIServer(os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"), databaseConnection)
	apiServer.Run()
}
