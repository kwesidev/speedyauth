package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kwesidev/authserver/internal/models"
	"github.com/kwesidev/authserver/internal/services"
	"github.com/kwesidev/authserver/internal/utilities"
)

// Registered Service

type UserController struct {
	// Registered Services
	db          *sql.DB
	userService services.UserService
	validate    *validator.Validate
}

// Creates a new User Controller for passing requests
func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		db:          db,
		userService: *services.NewUserService(db),
		validate:    validator.New(),
	}
}

// Index  Welcome user
func (this *UserController) Index(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	userDetails := this.userService.Get(userId)
	utilities.JSONResponse(w, userDetails)
}

// Update user
func (this *UserController) Update(w http.ResponseWriter, r *http.Request) {
	userUpdateRequest := models.UserUpdateRequest{}
	err := utilities.GetJsonInput(&userUpdateRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = this.validate.Struct(userUpdateRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	response := struct {
		Success bool `json:"success"`
	}{}
	err = this.userService.Update(userId, userUpdateRequest)
	if err != nil {
		utilities.JSONError(w, "Failed to Update ", http.StatusBadRequest)
		return
	}
	response.Success = true
	utilities.JSONResponse(w, response)

}

// Logout function to logout user
func (this *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenRefreshRequest := models.TokenRefreshRequest{}
	err := utilities.GetJsonInput(&tokenRefreshRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Success bool `json:"success"`
	}{}
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	success, err := this.userService.DeleteToken(userId, tokenRefreshRequest.RefreshToken)

	if err != nil {
		response.Success = false
		utilities.JSONError(w, "Failed to register", http.StatusBadRequest)
		return
	}
	response.Success = success
	utilities.JSONResponse(w, response)
}
