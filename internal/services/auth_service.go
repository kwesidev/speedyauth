package services

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/kwesidev/speedyauth/internal/models"
	"github.com/kwesidev/speedyauth/internal/utilities"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginByUsernamePassword(username, password, ipAddress, userAgent string) (*models.AuthenticationResponse, error)
	PasswordlessLogin(username, sendMethod, ipAddress, userAgent string) (*models.PasswordLessAuthResponse, error)
	CompletePasswordLessLogin(code, requestId string) (*models.AuthenticationResponse, error)
	ResetPasswordRequest(username string) (bool, error)
	GenerateRefreshToken(oldRefreshToken, ipAddress, userAgent string) (*models.AuthenticationResponse, error)
	VerifyAndSetNewPassword(code string, password string) (bool, error)
	DeleteExpiredTokens(days int) error
	VerifyPassCode(userId int, passCode string) bool
	VerifyTOTP(userId int, passCode, ipAddress, userAgent string) (*models.AuthenticationResponse, error)
	ValidateTwoFactor(code, requestId string, ipAddress, userAgent string) (*models.AuthenticationResponse, error)
}
type authService struct {
	db           *sql.DB
	userService  UserService
	emailService EmailService
	tokenTime    time.Duration
}

func NewAuthService(db *sql.DB, userService UserService, emailService EmailService) AuthService {
	tokenTime, _ := time.ParseDuration(os.Getenv("TOKEN_EXPIRY_TIME"))
	return &authService{
		db:           db,
		userService:  userService,
		emailService: emailService,
		tokenTime:    tokenTime,
	}

}

// Login function to authenticate user by username and password
func (authSrv *authService) LoginByUsernamePassword(username, password, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId       int
		passwordHash string
		err          error
	)
	row := authSrv.db.QueryRow("SELECT id, password FROM users WHERE username = $1  LIMIT 1 ", username)
	row.Scan(&userId, &passwordHash)
	// Check if username is valid
	if userId == 0 {
		return nil, ErrInvalidUsername
	}
	userDetails := authSrv.userService.Get(userId)
	if !userDetails.Active {
		return nil, ErrAccountNotActive
	}
	// Validates password
	if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}
	return authSrv.generateAuthResponse(*userDetails, ipAddress, userAgent)
}
func (authSrv *authService) generateAuthResponse(userDetails models.User, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	// Check if two authentication is required
	if userDetails.TwoFactorEnabled {
		if userDetails.TwoFactorMethod != "TOTP" {
			return authSrv.twoFactorRequest(userDetails, ipAddress, userAgent)
		}
		// Otherwise its TOTP then
		authResult := &models.AuthenticationResponse{}
		// Generate a short token which expires after 5minutes
		shortToken, _ := utilities.GenerateJwtToken(userDetails.ID, userDetails.Roles, (time.Second * 300))
		authResult.TwoFactorEnabled = true
		authResult.Token = shortToken
		authResult.TwoFactorMethod = userDetails.TwoFactorMethod
		return authResult, nil
	}
	// Get user roles
	roles, err := authSrv.userService.GetRoles(userDetails.ID)
	userDetails.Roles = roles
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return authSrv.generateTokenDetails(userDetails, ipAddress, userAgent)
}

// Func loginByUsername this will send an otp to the user which then be verified
func (authSrv *authService) PasswordlessLogin(username, sendMethod, ipAddress, userAgent string) (*models.PasswordLessAuthResponse, error) {
	userDetails := authSrv.userService.GetByUsername(username)
	if userDetails == nil {
		return nil, ErrInvalidUsername
	}
	if !userDetails.Active {
		return nil, ErrAccountNotActive
	}
	tx, err := authSrv.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Generates request ID
	requestId := utilities.GenerateOpaqueToken(45)
	// Generate 6 random code
	randomCodes := utilities.GenerateRandomDigits(6)
	if _, err = tx.Exec("INSERT INTO otp_requests(user_id, request_id, code, send_method, expiry_time, ip_address, user_agent, created_at) values($1, $2, $3, $4, $5, $6, $7, NOW())", userDetails.ID, requestId, randomCodes, "EMAIL", time.Now().Add(1*time.Minute), ipAddress, userAgent); err != nil {
		log.Println(err)
		return nil, err
	}
	if err = authSrv.emailService.SendEmailLoginRequest(randomCodes, *userDetails); err != nil {
		log.Println("Email Error", err)
		return nil, ErrSendingMail
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	passwordLessAuthResponse := models.PasswordLessAuthResponse{}
	passwordLessAuthResponse.RequestId = requestId
	passwordLessAuthResponse.SendMethod = "EMAIL"
	return &passwordLessAuthResponse, nil
}

// Func completePasswordLessLogin
func (authSrv *authService) CompletePasswordLessLogin(code, requestId string) (*models.AuthenticationResponse, error) {
	var (
		userId               int
		userAgent, ipAddress string
	)
	tx, err := authSrv.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	row := tx.QueryRow("SELECT user_id, ip_address,user_agent FROM otp_requests WHERE code = $1 AND request_id = $2 AND expiry_time >= NOW()", code, requestId)
	row.Scan(&userId, &ipAddress, &userAgent)
	if userId == 0 {
		log.Println("Invalid Code or Request Id Invalid")
		return nil, ErrInvalidCode
	}
	userDetails := authSrv.userService.Get(userId)
	// Deletes the otp requests
	if _, err = tx.Exec("DELETE FROM otp_requests WHERE request_id = $1", requestId); err != nil {
		log.Println(err)
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return authSrv.generateAuthResponse(*userDetails, ipAddress, userAgent)
}

// Refresh Token generates a new refresh token that will be used to get a new access token and a refresh token
func (authSrv *authService) GenerateRefreshToken(oldRefreshToken, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId int
	)
	authResult := &models.AuthenticationResponse{}
	row := authSrv.db.QueryRow("SELECT user_id FROM user_refresh_tokens WHERE token = $1 AND expiry_time > NOW() ", oldRefreshToken)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Refresh Token is not there")
		return nil, ErrInvalidToken
	}
	// Check if account is active before refreshing token
	userDetails := authSrv.userService.Get(userId)
	if !userDetails.Active {
		return nil, ErrAccountNotActive
	}
	roles, _ := authSrv.userService.GetRoles(userId)
	tokenExpiry := time.Duration(authSrv.tokenTime)

	jwtToken, err := utilities.GenerateJwtToken(userId, roles, time.Duration(tokenExpiry))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Delete the old token and generate new access token and refresh token
	refreshToken := utilities.GenerateOpaqueToken(45)
	tx, err := authSrv.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if _, err = tx.Exec("DELETE FROM user_refresh_tokens WHERE user_id = $1 AND token = $2", userId, oldRefreshToken); err != nil {
		log.Println(err)
		return nil, err
	}
	queryString :=
		`INSERT 
			INTO user_refresh_tokens
				(user_id, token, created, ip_address, user_agent, expiry_time)
	        VALUES
	        	($1, $2 ,NOW() ,$3 ,$4, $5)
	    `
	// Generate a jwt and refresh token
	if _, err = tx.Exec(queryString, userId, refreshToken, ipAddress, userAgent, time.Now().Add(tokenExpiry)); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = roles
	authResult.Expires = int(tokenExpiry.Seconds())
	return authResult, nil
}

func (authSrv *authService) ResetPasswordRequest(username string) (bool, error) {
	//check if the username exists and then send vertification code
	var (
		userId int
	)
	row := authSrv.db.QueryRow("SELECT id FROM users where username = $1 OR email_address = $1 ", username)
	if err := row.Scan(&userId); err != nil {
		log.Println(err)
		return false, ErrInvalidUsername
	}
	userDetails := authSrv.userService.Get(userId)
	if !userDetails.Active {
		return false, ErrAccountNotActive
	}
	tx, err := authSrv.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return false, err
	}
	// Delete any requested password requests
	if _, err = tx.Exec("DELETE FROM reset_password_requests WHERE user_id = $1", userId); err != nil {
		log.Println(err)
		return false, err
	}
	// Generate some random code to be sent to the user for reseting of the password
	randomCodes := utilities.GenerateRandomDigits(6)
	if _, err = tx.Exec("INSERT INTO reset_password_requests(user_id, code, created, expiry_time) values($1, $2, NOW(), $3)", userId, randomCodes, time.Now().Add(30*time.Minute)); err != nil {
		log.Println(err)
		return false, err
	}
	if err = authSrv.emailService.SendPasswordResetRequest(randomCodes, *userDetails); err != nil {
		log.Println("Email Error", err)
		return false, ErrSendingMail
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// VerifyAndSetNewPassword functions to verify and reset password
func (authSrv *authService) VerifyAndSetNewPassword(code string, password string) (bool, error) {
	if !utilities.StrongPasswordCheck(password) {
		return false, ErrStrongPassword
	}
	var userId int
	tx, err := authSrv.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return false, err
	}
	row := tx.QueryRow("SELECT user_id FROM reset_password_requests WHERE code = $1 AND expiry_time >= NOW()", code)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Invalid Code")
		return false, ErrInvalidCode
	}
	// update password and delete all refresh tokens
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return false, ErrPasswordUpdate
	}
	if _, err = tx.Exec("UPDATE users SET password = $2 WHERE id = $1", userId, passwordHash); err != nil {
		log.Println(err)
		return false, ErrPasswordUpdate
	}
	if _, err = tx.Exec("DELETE FROM user_refresh_tokens WHERE user_id = $1", userId); err != nil {
		log.Println(err)
		return false, ErrPasswordUpdate
	}
	_, err = tx.Exec("DELETE FROM reset_password_requests WHERE user_id = $1", userId)
	if err != nil {
		log.Println(err)
		return false, ErrPasswordUpdate
	}
	if err = tx.Commit(); err != nil {
		return false, ErrInvalidPassword
	}
	return true, nil
}

func (authSrv *authService) twoFactorRequest(userDetails models.User, ipAddress string, userAgent string) (*models.AuthenticationResponse, error) {
	tx, err := authSrv.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	authResult := &models.AuthenticationResponse{}
	// Expire after 5minutes
	expires := time.Duration(300 * time.Second)
	randomCodes := utilities.GenerateRandomDigits(6)
	requestId := utilities.GenerateOpaqueToken(60)
	queryString :=
		`INSERT 
		    INTO two_factor_requests 
            	(user_id, request_id, ip_address, code, user_agent, created_at, send_method, expiry_time)
	        VALUES
	        	($1, $2 ,$3 ,$4, $5, NOW(),'EMAIL', $6)
	    `
	if _, err = tx.Exec(queryString, userDetails.ID, requestId, ipAddress, randomCodes, userAgent, time.Now().Add(expires)); err != nil {
		log.Println(err)
		return nil, ErrTwoFactorRequest
	}
	if err = authSrv.emailService.SendTwoFactorRequest(randomCodes, userDetails); err != nil {
		log.Println("Sending Email error", err)
		return nil, ErrSendingMail
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	authResult.TwoFactorEnabled = true
	authResult.Token = requestId
	authResult.TwoFactorMethod = userDetails.TwoFactorMethod
	return authResult, nil
}

func (authSrv *authService) generateTokenDetails(userDetails models.User, ipAddress string, userAgent string) (*models.AuthenticationResponse, error) {
	authResult := &models.AuthenticationResponse{}
	tokenExpiry := time.Duration(authSrv.tokenTime)
	// Generates JWT Token and Refresh token that expires after xminutes
	jwtToken, err := utilities.GenerateJwtToken(userDetails.ID, userDetails.Roles, tokenExpiry)
	if err != nil {
		log.Println(err)
		return nil, ErrAccessToken
	}
	refreshToken := utilities.GenerateOpaqueToken(45)
	queryString :=
		`INSERT 
		    INTO user_refresh_tokens
            	(user_id, token, created, ip_address, user_agent, expiry_time)
	        VALUES
	        	($1, $2 ,NOW() ,$3 ,$4, $5)
	    `
	// Generate a jwt and refresh token
	if _, err = authSrv.db.Exec(queryString, userDetails.ID, refreshToken, ipAddress, userAgent, time.Now().Add(tokenExpiry)); err != nil {
		log.Println(err)
		return nil, ErrTokenGeneration
	}
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = userDetails.Roles
	authResult.Expires = int(tokenExpiry.Seconds())
	authResult.TwoFactorEnabled = userDetails.TwoFactorEnabled
	return authResult, nil
}

// Validate the two factor authentication request and complete the authentication request
func (authSrv *authService) ValidateTwoFactor(code, requestId string, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var userId int
	row := authSrv.db.QueryRow("SELECT user_id FROM two_factor_requests WHERE code = $1 AND request_id = $2 AND expiry_time > NOW() ", code, requestId)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Invalid Code")
		return nil, ErrTwoFactorCode
	}
	if _, err := authSrv.db.Exec("DELETE FROM two_factor_requests WHERE code = $1 AND request_id = $2", code, requestId); err != nil {
		log.Println(err)
		return nil, ErrTwoFactorCode
	}
	userDetails := authSrv.userService.Get(userId)
	return authSrv.generateTokenDetails(*userDetails, ipAddress, userAgent)

}

// Delete expired tokens
func (authSrv *authService) DeleteExpiredTokens(days int) error {
	// Deletes User Refresh tokens
	result, err := authSrv.db.Exec("DELETE FROM user_refresh_tokens WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	log.Println("DELETED number of rows for user expired tokens :", count)
	// Deletes Two factor requests
	result, err = authSrv.db.Exec("DELETE FROM two_factor_requests WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	if err != nil {
		return err
	}
	count, _ = result.RowsAffected()
	log.Println("DELETED number of rows for two_factor_requests tokens :", count)
	// Delete Reset Password Requests
	result, err = authSrv.db.Exec("DELETE FROM reset_password_requests WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ = result.RowsAffected()
	log.Println("DELETED number of rows for reset_password_requests tokens :", count)
	return err
}

// Verify the passcode
func (authSrv *authService) VerifyPassCode(userId int, passCode string) bool {
	userDetails := authSrv.userService.Get(userId)
	return totp.Validate(passCode, userDetails.TOTPSecret)
}

// Validates the TOTP before the user finally logs in
func (authSrv *authService) VerifyTOTP(userId int, passCode, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	userDetails := authSrv.userService.Get(userId)
	if !authSrv.VerifyPassCode(userId, passCode) {
		return nil, ErrPassCode
	}
	return authSrv.generateTokenDetails(*userDetails, ipAddress, userAgent)
}
