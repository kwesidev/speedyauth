package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/kwesidev/authserver/internal/controllers"
	"github.com/kwesidev/authserver/internal/middlewares"
	"github.com/kwesidev/authserver/internal/utilities"
)

var databaseConnection *sql.DB

// initializes the database connections and other connections
func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseConnection = utilities.GetMainDatabaseConnection()
}

// register Global functions
func registerGlobalFunctions() {
	authController := controllers.NewAuthController(databaseConnection)
	http.HandleFunc("/api/auth/login", authController.Login)
	http.HandleFunc("/api/auth/tokenRefresh", authController.RefreshToken)
	http.HandleFunc("/api/auth/register", authController.Register)
	http.HandleFunc("/api/auth/logout", authController.Logout)
	http.HandleFunc("/api/auth/passwordResetRequest", authController.PasswordResetRequest)
	http.HandleFunc("/api/auth/verifyAndResetPassword", authController.VerifyAndChangePassword)
}

// register user functions
func registerUserFunctions() {
	userController := controllers.NewUserController(databaseConnection)
	http.HandleFunc("/api/user", middlewares.JwtAuth(userController.Index))
}

// register admin functions
func registerAdminFunctions() {

}
func main() {
	initialize()
	registerGlobalFunctions()
	registerAdminFunctions()
	registerUserFunctions()
	log.Println("Starting Auth Server listening for requests on port 8080")
	// Listen to incoming connections
	err := http.ListenAndServe(":8080", middlewares.LogRequest(http.DefaultServeMux))
	// Exit if fail to start service
	if err != nil {
		log.Fatal("Failed to start Service ")
	}

}
