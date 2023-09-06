package main

import (
	"database/sql"
	"log"
	"os"

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
	databaseConnection = utilities.GetMainDatabaseConnection()
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
