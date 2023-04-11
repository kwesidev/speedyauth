package controllers

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/go-playground/validator/v10"
	"github.com/kwesidev/authserver/internal/models"
	"github.com/kwesidev/authserver/internal/services"
	"github.com/kwesidev/authserver/internal/utilities"
)

type AuthController struct {
	// Registered Services
	db          *sql.DB
	userService services.UserService
	authService services.AuthService
	validate    *validator.Validate
}

// Creates a new Auth Controller for passing requests
func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		db:          db,
		userService: *services.NewUserService(db),
		authService: *services.NewAuthService(db),
		validate:    validator.New(),
	}
}

// Login Handler To Authenticate user
func (this *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	authRequest := models.AuthenticationRequest{}
	err := utilities.GetJsonInput(&authRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Validates the requests
	err = this.validate.Struct(authRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	authResult, err := this.authService.Login(authRequest.Username, authRequest.Password, "localhost", "App")
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusForbidden)
		return
	}
	utilities.JSONResponse(w, authResult)
}

// Function To Refresh Token
func (this *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenRefreshRequest := models.TokenRefreshRequest{}
	err := utilities.GetJsonInput(&tokenRefreshRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	refreshResult, err := this.authService.GenerateRefreshToken(tokenRefreshRequest.RefreshToken, "localhost", "App")
	if err != nil {
		utilities.JSONError(w, "Failed to generate Token", http.StatusForbidden)
		return
	}
	utilities.JSONResponse(w, refreshResult)
}

// Logout function to logout user
func (this *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenRefreshRequest := models.TokenRefreshRequest{}
	err := utilities.GetJsonInput(&tokenRefreshRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Success bool `json:"success"`
	}{}
	success, err := this.authService.DeleteToken(tokenRefreshRequest.RefreshToken)
	if err != nil {
		utilities.JSONError(w, "Failed to logout ", http.StatusBadRequest)
		return
	}
	response.Success = success
	utilities.JSONResponse(w, response)
}

// Function register User
func (this *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	userRegisterationRequest := models.UserRegistrationRequest{}
	err := utilities.GetJsonInput(&userRegisterationRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validates requests
	err = this.validate.Struct(userRegisterationRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Success bool `json:"success"`
	}{}
	regResult, err := this.userService.Register(userRegisterationRequest)
	if err != nil {
		utilities.JSONError(w, "Failed to register user", http.StatusBadRequest)
		return
	}
	response.Success = regResult
	utilities.JSONResponse(w, response)
}
