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
	port         string
	serverName   string
	db           *sql.DB
	userService  services.UserService
	authService  services.AuthService
	emailService services.EmailService
}

// NewAPIServer initializes the api server
func NewAPIServer(serverName string, port string, db *sql.DB) *APIServer {
	return &APIServer{serverName: serverName, port: port, db: db}
}

func (api *APIServer) setupRoutes() {
	api.registerGlobalFunctions()
	api.registerAdminFunctions()
	api.registerUserFunctions()
}

func (api *APIServer) setUpServices() {
	api.emailService = services.NewEmailService(true)
	api.userService = services.NewUserService(api.db)
	api.authService = services.NewAuthService(api.db, api.userService, api.emailService)
}

// Run  start serving the http requests
func (api *APIServer) Run() {
	api.setUpServices()
	api.setupRoutes()
	api.cleanUp()
	// Listen to incoming connections
	log.Println("Starting SpeedyAuth listening for requests on port " + os.Getenv("SERVER_PORT"))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", api.serverName, api.port), middlewares.LogRequest(http.DefaultServeMux))
	// Exit if fail to start service
	if err != nil {
		log.Fatal("Failed to start Service ")
	}
}

// register Global functions
func (api *APIServer) registerGlobalFunctions() {
	authController := controllers.NewAuthController(api.db, api.authService, api.userService)
	http.HandleFunc("/api/auth/login", middlewares.Method("POST", authController.Login))
	http.HandleFunc("/api/auth/passwordless/login", middlewares.Method("POST", authController.PasswordlessLogin))
	http.HandleFunc("/api/auth/passwordless/complete", middlewares.Method("POST", authController.CompletePasswordlessLogin))
	http.HandleFunc("/api/auth/token/refresh", middlewares.Method("POST", authController.RefreshToken))
	http.HandleFunc("/api/auth/register", middlewares.Method("POST", authController.Register))
	http.HandleFunc("/api/auth/password/reset/request", middlewares.Method("POST", authController.PasswordResetRequest))
	http.HandleFunc("/api/auth/password/reset/verify", middlewares.Method("POST", authController.VerifyAndChangePassword))
	http.HandleFunc("/api/auth/twofactor/verify", middlewares.Method("POST", authController.ValidateTwoFactor))
	http.HandleFunc("/health", authController.Health)
}

// register user functions

func (api *APIServer) registerUserFunctions() {
	userController := controllers.NewUserController(api.db, api.userService, api.authService)
	http.HandleFunc("/api/user", middlewares.Method("GET", middlewares.JwtAuth(userController.Index)))
	http.HandleFunc("/api/user/logout", middlewares.Method("POST", middlewares.JwtAuth(userController.Logout)))
	http.HandleFunc("/api/user/update", middlewares.Method("POST", middlewares.JwtAuth(userController.Update)))
	http.HandleFunc("/api/user/twofactor/enable", middlewares.Method("POST", middlewares.JwtAuth(userController.EnableTwoFactor)))
	http.HandleFunc("/api/user/totpcode/verify", middlewares.Method("POST", middlewares.JwtAuth(userController.VerifyPassCode)))

}

// register admin functions
func (api *APIServer) registerAdminFunctions() {

}

// Cleanup
func (api *APIServer) cleanUp() {
	//Deletes expired tokens after 30 days
	err := api.authService.DeleteExpiredTokens(30)
	if err != nil {
		log.Fatal("There was a problem cleaning up ")
	}
}
