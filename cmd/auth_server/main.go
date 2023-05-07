package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

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
	databaseConnection = utilities.GetMainDatabaseConnection()
}
func main() {
	initialize()
	apiServer := apiserver.NewServer(os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"), databaseConnection)
	apiServer.Run()
}
