package apiserver

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kwesidev/authserver/internal/controllers"
	"github.com/kwesidev/authserver/internal/middlewares"
)

type APIServer struct {
	port       string
	serverName string
	db         *sql.DB
}

// NewAPIServer initializes the api server
func NewAPIServer(serverName string, port string, db *sql.DB) *APIServer {
	return &APIServer{serverName: serverName, port: port, db: db}
}

func (this *APIServer) setupRoutes() {
	this.registerGlobalFunctions()
	this.registerAdminFunctions()
	this.registerUserFunctions()
}

// Run  start serving the http requests
func (this *APIServer) Run() {
	this.setupRoutes()
	// Listen to incoming connections
	log.Println("Starting Auth Server listening for requests on port " + os.Getenv("SERVER_PORT"))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", this.serverName, this.port), middlewares.LogRequest(http.DefaultServeMux))

	// Exit if fail to start service
	if err != nil {
		log.Fatal("Failed to start Service ")
	}
}

// register Global functions
func (this *APIServer) registerGlobalFunctions() {
	authController := controllers.NewAuthController(this.db)
	http.HandleFunc("/api/auth/login", middlewares.Method("POST", authController.Login))
	http.HandleFunc("/api/auth/tokenRefresh", middlewares.Method("POST", authController.RefreshToken))
	http.HandleFunc("/api/auth/register", middlewares.Method("POST", authController.Register))
	http.HandleFunc("/api/auth/logout", middlewares.Method("POST", authController.Logout))
	http.HandleFunc("/api/auth/passwordResetRequest", middlewares.Method("POST", authController.PasswordResetRequest))
	http.HandleFunc("/api/auth/verifyAndResetPassword", middlewares.Method("POST", authController.VerifyAndChangePassword))
	http.HandleFunc("/health", authController.Health)
}

// register user functions
func (this *APIServer) registerUserFunctions() {
	userController := controllers.NewUserController(this.db)
	http.HandleFunc("/api/user", middlewares.Method("GET", middlewares.JwtAuth(userController.Index)))
	http.HandleFunc("/api/user/update", middlewares.Method("POST", middlewares.JwtAuth(userController.Update)))
}

// register admin functions
func (this *APIServer) registerAdminFunctions() {

}
