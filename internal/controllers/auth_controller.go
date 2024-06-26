package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/go-playground/validator/v10"
	"github.com/kwesidev/speedyauth/internal/models"
	"github.com/kwesidev/speedyauth/internal/services"
	"github.com/kwesidev/speedyauth/internal/utilities"
)

type AuthController struct {
	// Registered Services
	db          *sql.DB
	userService services.UserService
	authService services.AuthService
	validate    *validator.Validate
}

// Creates a new Auth Controller for passing requests
func NewAuthController(db *sql.DB, authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		db:          db,
		userService: userService,
		authService: authService,
		validate:    validator.New(),
	}
}

// Login Handler To Authenticate user
func (authCtrl *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	authRequest := models.AuthenticationRequest{}
	err := utilities.GetJsonInput(&authRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = authCtrl.validate.Struct(authRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authResult *models.AuthenticationResponse
	authResult, err = authCtrl.authService.LoginByUsernamePassword(authRequest.Username, authRequest.Password, "", "")

	if err != nil {
		if errors.Is(err, services.ErrInvalidUsername) || errors.Is(err, services.ErrInvalidPassword) || errors.Is(err, services.ErrAccountNotActive) {
			utilities.JSONError(w, err.Error(), http.StatusUnauthorized)
		} else {
			utilities.JSONError(w, services.ErrServer.Error(), http.StatusInternalServerError)
		}
		return
	}
	utilities.JSONResponse(w, authResult)
}

// Login Handler To Authenticate user without passsword
func (authCtrl *AuthController) PasswordlessLogin(w http.ResponseWriter, r *http.Request) {
	passwordLessAuthRequest := models.PasswordLessAuthRequest{}
	err := utilities.GetJsonInput(&passwordLessAuthRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = authCtrl.validate.Struct(passwordLessAuthRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var passwordLessAuthResponse *models.PasswordLessAuthResponse
	passwordLessAuthResponse, err = authCtrl.authService.PasswordlessLogin(passwordLessAuthRequest.Username, passwordLessAuthRequest.SendMethod, "", "")
	if err != nil {
		if errors.Is(err, services.ErrInvalidUsername) || errors.Is(err, services.ErrInvalidPassword) || errors.Is(err, services.ErrAccountNotActive) {
			utilities.JSONError(w, err.Error(), http.StatusUnauthorized)
		} else {
			utilities.JSONError(w, services.ErrServer.Error(), http.StatusInternalServerError)
		}
		return
	}
	utilities.JSONResponse(w, passwordLessAuthResponse)
}

// Completes passwordless login
func (authCtrl *AuthController) CompletePasswordlessLogin(w http.ResponseWriter, r *http.Request) {
	completePasswordLessLogin := models.CompletePasswordLessRequest{}
	err := utilities.GetJsonInput(&completePasswordLessLogin, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validates requests
	err = authCtrl.validate.Struct(completePasswordLessLogin)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var (
		authResult *models.AuthenticationResponse
	)
	authResult, err = authCtrl.authService.CompletePasswordLessLogin(completePasswordLessLogin.Code, completePasswordLessLogin.RequestId)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCode) {
			utilities.JSONError(w, err.Error(), http.StatusUnauthorized)
		} else {
			utilities.JSONError(w, services.ErrServer.Error(), http.StatusInternalServerError)
		}
		return
	}
	utilities.JSONResponse(w, authResult)
}

// Function To Refresh Token
func (authCtrl *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenRefreshRequest := models.TokenRefreshRequest{}
	err := utilities.GetJsonInput(&tokenRefreshRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	refreshResult, err := authCtrl.authService.GenerateRefreshToken(tokenRefreshRequest.RefreshToken, r.RemoteAddr, r.UserAgent())
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			utilities.JSONError(w, err.Error(), http.StatusUnauthorized)
		} else {
			utilities.JSONError(w, services.ErrServer.Error(), http.StatusUnauthorized)
		}
		return
	}
	utilities.JSONResponse(w, refreshResult)
}

// Reset Password Request
func (authCtrl *AuthController) PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	passwordResetRequest := models.PasswordResetRequest{}
	err := utilities.GetJsonInput(&passwordResetRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := models.SuccessResponse{}
	success, err := authCtrl.authService.ResetPasswordRequest(passwordResetRequest.Username)
	if err != nil {
		utilities.JSONError(w, "Failed to Send Reset password Request", http.StatusBadRequest)
		return
	}
	response.Success = success
	utilities.JSONResponse(w, response)
}

// Verify and update the password
func (authCtrl *AuthController) VerifyAndChangePassword(w http.ResponseWriter, r *http.Request) {
	verifyAndChangePasswordRequest := models.VerifyChangePasswordRequest{}
	err := utilities.GetJsonInput(&verifyAndChangePasswordRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validates requests
	err = authCtrl.validate.Struct(verifyAndChangePasswordRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := models.SuccessResponse{}
	success, err := authCtrl.authService.VerifyAndSetNewPassword(verifyAndChangePasswordRequest.Code, verifyAndChangePasswordRequest.Password)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success = success
	utilities.JSONResponse(w, response)
}

// Function register User
func (authCtrl *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	userRegisterationRequest := models.UserRegistrationRequest{}
	err := utilities.GetJsonInput(&userRegisterationRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validates requests
	err = authCtrl.validate.Struct(userRegisterationRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := models.SuccessResponse{}
	regResult, err := authCtrl.userService.Register(userRegisterationRequest)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success = regResult
	utilities.JSONResponse(w, response)
}

// Validates Two Factor authCtrl function is only called when two factor is required
func (authCtrl *AuthController) ValidateTwoFactor(w http.ResponseWriter, r *http.Request) {
	twoFactorRequest := models.VerifyTwoFactorRequest{}
	err := utilities.GetJsonInput(&twoFactorRequest, r)
	if err != nil {
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validates requests
	err = authCtrl.validate.Struct(twoFactorRequest)
	if err != nil {
		log.Println(err)
		utilities.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check the type of verification
	var (
		authResult    *models.AuthenticationResponse
		userDataToken map[string]interface{}
	)
	if twoFactorRequest.Method == "TOTP" {
		userDataToken, err = utilities.ValidateJwtAndGetClaims(twoFactorRequest.Token)
		if err != nil {
			utilities.JSONError(w, "Invalid Token", http.StatusBadRequest)
			return
		}
		userId := userDataToken["userId"].(int)
		authResult, err = authCtrl.authService.VerifyTOTP(userId, twoFactorRequest.Code, "", "")

	} else {
		authResult, err = authCtrl.authService.ValidateTwoFactor(twoFactorRequest.Code, twoFactorRequest.Token, "", "")
	}
	if err != nil {
		utilities.JSONError(w, "Failed to Complete the authentication, Code provided is wrong", http.StatusBadRequest)
		return
	}
	utilities.JSONResponse(w, authResult)
}

// Health
func (authCtrl *AuthController) Health(w http.ResponseWriter, r *http.Request) {
	utilities.JSONResponse(w, "OKAY")
}
