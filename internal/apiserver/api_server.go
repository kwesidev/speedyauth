package apiserver

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kwesidev/speedyauth/internal/controllers"
	"github.com/kwesidev/speedyauth/internal/middlewares"
	"github.com/kwesidev/speedyauth/internal/services"
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

func (ap *APIServer) setupRoutes() {
	ap.registerGlobalFunctions()
	ap.registerAdminFunctions()
	ap.registerUserFunctions()
}

// Run  start serving the http requests
func (ap *APIServer) Run() {
	ap.cleanUp()
	ap.setupRoutes()
	// Listen to incoming connections
	log.Println("Starting SpeedyAuth listening for requests on port " + os.Getenv("SERVER_PORT"))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", ap.serverName, ap.port), middlewares.LogRequest(http.DefaultServeMux))

	// Exit if fail to start service
	if err != nil {
		log.Fatal("Failed to start Service ")
	}
}

// register Global functions
func (ap *APIServer) registerGlobalFunctions() {
	authController := controllers.NewAuthController(ap.db)
	http.HandleFunc("/api/auth/login", middlewares.Method("POST", authController.Login))
	http.HandleFunc("/api/auth/passwordLesslogin", middlewares.Method("POST", authController.PasswordlessLogin))
	http.HandleFunc("/api/auth/completePasswordlessLogin", middlewares.Method("POST", authController.CompletePasswordlessLogin))
	http.HandleFunc("/api/auth/tokenRefresh", middlewares.Method("POST", authController.RefreshToken))
	http.HandleFunc("/api/auth/register", middlewares.Method("POST", authController.Register))
	http.HandleFunc("/api/auth/passwordResetRequest", middlewares.Method("POST", authController.PasswordResetRequest))
	http.HandleFunc("/api/auth/verifyAndResetPassword", middlewares.Method("POST", authController.VerifyAndChangePassword))
	http.HandleFunc("/api/auth/verifyTwoFactor", middlewares.Method("POST", authController.ValidateTwoFactor))
	http.HandleFunc("/health", authController.Health)
}

// register user functions
func (ap *APIServer) registerUserFunctions() {
	userController := controllers.NewUserController(ap.db)
	http.HandleFunc("/api/user", middlewares.Method("GET", middlewares.JwtAuth(userController.Index)))
	http.HandleFunc("/api/user/logout", middlewares.Method("POST", middlewares.JwtAuth(userController.Logout)))
	http.HandleFunc("/api/user/update", middlewares.Method("POST", middlewares.JwtAuth(userController.Update)))
	http.HandleFunc("/api/user/enableTwoFactor", middlewares.Method("POST", middlewares.JwtAuth(userController.EnableTwoFactor)))
	http.HandleFunc("/api/user/verifyPassCode", middlewares.Method("POST", middlewares.JwtAuth(userController.VerifyPassCode)))

}

// register admin functions
func (ap *APIServer) registerAdminFunctions() {

}

// Cleanup
func (ap *APIServer) cleanUp() {
	authService := services.NewAuthService(ap.db)
	// Deletes expired tokens after 30 days
	err := authService.DeleteExpiredTokens(30)
	if err != nil {
		log.Fatal("There was a problem cleaning up ")
	}
}
