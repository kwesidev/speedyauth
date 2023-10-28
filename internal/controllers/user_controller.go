package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kwesidev/speedyauth/internal/models"
	"github.com/kwesidev/speedyauth/internal/services"
	"github.com/kwesidev/speedyauth/internal/utilities"
)

// Registered Service

type UserController struct {
	// Registered Services
	db          *sql.DB
	userService services.UserService
	authService services.AuthService
	validate    *validator.Validate
}

// Creates a new User Controller for passing requests
func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		db:          db,
		userService: *services.NewUserService(db),
		authService: *services.NewAuthService(db),
		validate:    validator.New(),
	}
}

// Index  Welcome user
func (usrCtrl *UserController) Index(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	userDetails := usrCtrl.userService.Get(userId)
	utilities.JSONResponse(w, userDetails)
}

// Update user
func (usrCtrl *UserController) Update(w http.ResponseWriter, r *http.Request) {
	userUpdateRequest := models.UserUpdateRequest{}
	err := utilities.GetJsonInput(&userUpdateRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = usrCtrl.validate.Struct(userUpdateRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	response := models.SuccessResponse{}
	err = usrCtrl.userService.Update(userId, userUpdateRequest)
	if err != nil {
		utilities.JSONError(w, "Failed to Update ", http.StatusBadRequest)
		return
	}
	response.Success = true
	utilities.JSONResponse(w, response)

}

// Logout function to logout user
func (usrCtrl *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenRefreshRequest := models.TokenRefreshRequest{}
	err := utilities.GetJsonInput(&tokenRefreshRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := models.SuccessResponse{}
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	success, err := usrCtrl.userService.DeleteToken(userId, tokenRefreshRequest.RefreshToken)

	if err != nil {
		response.Success = false
		utilities.JSONError(w, "Failed to register", http.StatusBadRequest)
		return
	}
	response.Success = success
	utilities.JSONResponse(w, response)
}

// Enable Time based OTP
func (usrCtrl *UserController) EnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	enableTwoFactorRequest := models.EnableTwoFactorRequest{}
	err := utilities.GetJsonInput(&enableTwoFactorRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	if enableTwoFactorRequest.Type == "TOTP" {
		totpResponse, err := usrCtrl.userService.EnableTwoFactorTOTP(userId)
		if err != nil {
			utilities.JSONError(w, "Failed to Enable Two Factor (TOTP)", http.StatusBadRequest)
			return
		}
		utilities.JSONResponse(w, totpResponse)
		return
	} else {
		err := usrCtrl.userService.EnableTwoFactor(userId, enableTwoFactorRequest.Type)
		if err != nil {
			utilities.JSONError(w, "Failed to Enabled Two Factor EMAIL OR SMS ", http.StatusBadRequest)
			return
		}
	}
	response := models.SuccessResponse{}
	response.Success = true
	utilities.JSONResponse(w, response)
}

// Enable Time based OTP
func (usrCtrl *UserController) VerifyPassCode(w http.ResponseWriter, r *http.Request) {
	verifyPassCodeRequest := models.VerifyPassCodeRequest{}
	err := utilities.GetJsonInput(&verifyPassCodeRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validates requests
	err = usrCtrl.validate.Struct(verifyPassCodeRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("claims").(map[string]interface{})
	userId := claims["userId"].(int)
	response := models.SuccessResponse{}

	if usrCtrl.authService.VerifyPassCode(userId, verifyPassCodeRequest.Code) {
		response.Success = true
	} else {
		response.Success = false
	}
	utilities.JSONResponse(w, response)

}
