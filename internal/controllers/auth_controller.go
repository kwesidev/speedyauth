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
	authResult, err := this.authService.Login(authRequest.Username, authRequest.Password, "", "")
	if err != nil {
		if err.Error() == services.ERROR_INVALID_PASSWORD || err.Error() == services.ERROR_ACCOUNT_NOT_ACTIVE || err.Error() == services.ERROR_INVALID_USERNAME {
			utilities.JSONError(w, err.Error(), http.StatusUnauthorized)
		} else {
			utilities.JSONError(w, services.ERROR_SERVER_ERROR, http.StatusInternalServerError)
		}
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
	refreshResult, err := this.authService.GenerateRefreshToken(tokenRefreshRequest.RefreshToken, r.RemoteAddr, r.UserAgent())
	if err != nil {
		if err.Error() == services.ERROR_TOKEN_INVALID {
			utilities.JSONError(w, services.ERROR_TOKEN_INVALID, http.StatusUnauthorized)
		} else {
			utilities.JSONError(w, services.ERROR_SERVER_ERROR, http.StatusUnauthorized)
		}
		return
	}
	utilities.JSONResponse(w, refreshResult)
}

// Reset Password Request
func (this *AuthController) PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	passwordResetRequest := models.PasswordResetRequest{}
	err := utilities.GetJsonInput(&passwordResetRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Success bool `json:"success"`
	}{}
	success, err := this.authService.ResetPasswordRequest(passwordResetRequest.Username)
	if err != nil {
		utilities.JSONError(w, "Failed to Send Reset password Request", http.StatusBadRequest)
		return
	}
	response.Success = success
	utilities.JSONResponse(w, response)
}

// Verify and update the password
func (this *AuthController) VerifyAndChangePassword(w http.ResponseWriter, r *http.Request) {
	verifyAndChangePasswordRequest := models.VerifyChangePasswordRequest{}
	err := utilities.GetJsonInput(&verifyAndChangePasswordRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validates requests
	err = this.validate.Struct(verifyAndChangePasswordRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := struct {
		Success bool `json:"success"`
	}{}
	success, err := this.authService.VerifyAndSetNewPassword(verifyAndChangePasswordRequest.Code, verifyAndChangePasswordRequest.Password)
	if err != nil {
		utilities.JSONError(w, services.ERROR_PASSWORD_UPDATE, http.StatusBadRequest)
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
		utilities.JSONError(w, services.ERROR_USER_REGISTRATION_FAILED, http.StatusBadRequest)
		return
	}
	response.Success = regResult
	utilities.JSONResponse(w, response)
}

// Validates Two Factor this function is only called when two factor is required
func (this *AuthController) ValidateTwoFactor(w http.ResponseWriter, r *http.Request) {
	twoFactorRequest := models.VerifyTwoFactorRequest{}
	err := utilities.GetJsonInput(&twoFactorRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validates requests
	err = this.validate.Struct(twoFactorRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	authResult, err := this.authService.ValidateTwoFactor(twoFactorRequest.Code, twoFactorRequest.RequestId, "", "")
	if err != nil {
		utilities.JSONError(w, "Failed to Complete the authentication", http.StatusBadRequest)
		return
	}
	utilities.JSONResponse(w, authResult)

}

// Health
func (this *AuthController) Health(w http.ResponseWriter, r *http.Request) {
	utilities.JSONResponse(w, "OKAY")
}
