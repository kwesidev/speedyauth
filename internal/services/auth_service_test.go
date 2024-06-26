package services

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kwesidev/speedyauth/internal/utilities"
	_ "github.com/lib/pq"
)

var db *sql.DB

func setUp() {
	var err error
	db, err = utilities.GetMainDatabaseConnection(utilities.DatabaseConfig{
		Host:     os.Getenv("PG_HOST"),
		Username: os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Database: os.Getenv("PG_DB"),
		Port:     os.Getenv("PG_PORT"),
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
		return
	}
	// userService := NewUserService(db)
	// _, err = userService.Register(models.UserRegistrationRequest{
	// 	Username:     "john.doe",
	// 	Password:     "password_2030333",
	// 	EmailAddress: "johndoe@localhost",
	// 	FirstName:    "john",
	// 	LastName:     "doe",
	// 	CellNumber:   "0731482947",
	// })

}
func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}
func TestLoginByUsernamePassword(t *testing.T) {
	userService := NewUserService(db)
	EmailService := NewEmailService(true)
	authService := NewAuthService(db, userService, EmailService)
	_, err := authService.LoginByUsernamePassword("john.doe", "password", "", "")
	if err != nil {
		t.Error("Failed to authenticate")
	}
}
