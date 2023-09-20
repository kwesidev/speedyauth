package services

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/kwesidev/authserver/internal/models"
	"github.com/kwesidev/authserver/internal/utilities"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db           *sql.DB
	userService  *UserService
	emailService *EmailService
	tokenTime    time.Duration
}

func NewAuthService(db *sql.DB) *AuthService {
	tokenTime, _ := time.ParseDuration(os.Getenv("TOKEN_EXPIRY_TIME"))
	return &AuthService{
		db:           db,
		userService:  NewUserService(db),
		emailService: NewEmailService(true),
		tokenTime:    tokenTime,
	}

}

// Login function to authenticate user
func (authSrv *AuthService) Login(username, password, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId       int
		passwordHash string
		err          error
	)
	row := authSrv.db.QueryRow("SELECT id, password FROM users WHERE username = $1  LIMIT 1 ", username)
	row.Scan(&userId, &passwordHash)
	// Check if username is valid
	if userId == 0 {
		return nil, ErrorInvalidUsername
	}
	userDetails := authSrv.userService.Get(userId)
	if userDetails.Active == false {
		return nil, ErrorAccountNotActive
	}
	// Validates password
	if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrorInvalidPassword
	}
	// Check if two authentication is required
	if userDetails.TwoFactorEnabled {
		if userDetails.TwoFactorType != "TOTP" {
			return authSrv.twoFactorRequest(*userDetails, ipAddress, userAgent)
		}
		// Otherwise its TOTP then
		authResult := &models.AuthenticationResponse{}
		// Generate a short token which expires after 5minutes
		shortToken, _ := utilities.GenerateJwtToken(userId, userDetails.Roles, (time.Second * 300))
		authResult.TwoFactorEnabled = true
		authResult.Token = shortToken
		authResult.TwoFactorType = userDetails.TwoFactorType
		return authResult, nil
	}
	// Get user roles
	roles, err := authSrv.userService.GetRoles(userId)
	userDetails.Roles = roles
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return authSrv.generateTokenDetails(*userDetails, ipAddress, userAgent)
}

// Refresh Token generates a new refresh token that will be used to get a new access token and a refresh token
func (authSrv *AuthService) GenerateRefreshToken(oldRefreshToken, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId int
	)
	authResult := &models.AuthenticationResponse{}
	row := authSrv.db.QueryRow("SELECT user_id FROM user_refresh_tokens WHERE token = $1 AND expiry_time > NOW() ", oldRefreshToken)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Refresh Token is not there")
		return nil, ErrorInvalidToken
	}
	// Check if account is active before refreshing token
	userDetails := authSrv.userService.Get(userId)
	if userDetails.Active == false {
		return nil, ErrorAccountNotActive
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

func (authSrv *AuthService) ResetPasswordRequest(username string) (bool, error) {
	//check if the username exists and then send vertification code
	var (
		userId int
	)
	row := authSrv.db.QueryRow("SELECT id FROM users where username = $1 OR email_address = $1 ", username)
	if err := row.Scan(&userId); err != nil {
		log.Println(err)
		return false, ErrorInvalidUsername
	}
	userDetails := authSrv.userService.Get(userId)
	if userDetails.Active == false {
		return false, ErrorAccountNotActive
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
func (authSrv *AuthService) VerifyAndSetNewPassword(code string, password string) (bool, error) {
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
		return false, ErrorInvalidCode
	}
	// update password and delete all refresh tokens
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return false, ErrorPasswordUpdate
	}
	if _, err = tx.Exec("UPDATE users SET password = $2 WHERE id = $1", userId, passwordHash); err != nil {
		log.Println(err)
		return false, ErrorPasswordUpdate
	}
	if _, err = tx.Exec("DELETE FROM user_refresh_tokens WHERE user_id = $1", userId); err != nil {
		log.Println(err)
		return false, ErrorPasswordUpdate
	}
	_, err = tx.Exec("DELETE FROM reset_password_requests WHERE user_id = $1", userId)
	if err != nil {
		log.Println(err)
		return false, ErrorPasswordUpdate
	}
	if err = tx.Commit(); err != nil {
		return false, ErrorInvalidPassword
	}
	return true, nil
}

func (authSrv *AuthService) twoFactorRequest(userDetails models.User, ipAddress string, userAgent string) (*models.AuthenticationResponse, error) {
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
            	(user_id, request_id, ip_address, code, user_agent, created_at, send_type, expiry_time)
	        VALUES
	        	($1, $2 ,$3 ,$4, $5, NOW(),'EMAIL', $6)
	    `
	if _, err = tx.Exec(queryString, userDetails.ID, requestId, ipAddress, randomCodes, userAgent, time.Now().Add(expires)); err != nil {
		log.Println(err)
		return nil, ErrorTwoFactorRequest
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
	authResult.TwoFactorType = userDetails.TwoFactorType
	return authResult, nil
}

func (authSrv *AuthService) generateTokenDetails(userDetails models.User, ipAddress string, userAgent string) (*models.AuthenticationResponse, error) {
	authResult := &models.AuthenticationResponse{}
	tokenExpiry := time.Duration(authSrv.tokenTime)
	// Generates JWT Token and Refresh token that expires after xminutes
	jwtToken, err := utilities.GenerateJwtToken(userDetails.ID, userDetails.Roles, tokenExpiry)
	if err != nil {
		log.Println(err)
		return nil, ErrorAccessToken
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
		return nil, ErrorTokenGeneration
	}
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = userDetails.Roles
	authResult.Expires = int(tokenExpiry.Seconds())
	authResult.TwoFactorEnabled = userDetails.TwoFactorEnabled
	return authResult, nil
}

// Validate the two factor authentication request and complete the authentication request
func (authSrv *AuthService) ValidateTwoFactor(code, requestId string, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var userId int
	row := authSrv.db.QueryRow("SELECT user_id FROM two_factor_requests WHERE code = $1 AND request_id = $2 AND expiry_time > NOW() ", code, requestId)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Invalid Code")
		return nil, ErrorTwoFactorCode
	}
	if _, err := authSrv.db.Exec("DELETE FROM two_factor_requests WHERE code = $1 AND request_id = $2", code, requestId); err != nil {
		log.Println(err)
		return nil, ErrorTwoFactorCode
	}
	userDetails := authSrv.userService.Get(userId)
	return authSrv.generateTokenDetails(*userDetails, ipAddress, userAgent)

}

// Delete expired tokens
func (authSrv *AuthService) DeleteExpiredTokens(days int) error {
	// Deletes User Refresh tokens
	result, err := authSrv.db.Exec("DELETE FROM user_refresh_tokens WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ := result.RowsAffected()
	log.Println("DELETED number of rows for user expired tokens :", count)
	// Deletes Two factor requests
	result, err = authSrv.db.Exec("DELETE FROM two_factor_requests WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ = result.RowsAffected()
	log.Println("DELETED number of rows for two_factor_requests tokens :", count)
	// Delete Reset Password Requests
	result, err = authSrv.db.Exec("DELETE FROM reset_password_requests WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ = result.RowsAffected()
	log.Println("DELETED number of rows for reset_password_requests tokens :", count)
	return err
}

// Verify the passcode
func (authSrv *AuthService) VerifyPassCode(userId int, passCode string) bool {
	userDetails := authSrv.userService.Get(userId)
	if valid := totp.Validate(passCode, userDetails.TOTPSecret); !valid {
		return false
	}
	return true
}

// Validates the TOTP before the user finally logs in
func (authSrv *AuthService) VerifyTOTP(userId int, passCode, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	userDetails := authSrv.userService.Get(userId)
	if authSrv.VerifyPassCode(userId, passCode) {
		return nil, ErrorPassCode
	}
	return authSrv.generateTokenDetails(*userDetails, ipAddress, userAgent)
}
