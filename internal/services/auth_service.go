package services

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kwesidev/authserver/internal/models"
	"github.com/kwesidev/authserver/internal/utilities"
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
func (this *AuthService) Login(username, password, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId       int
		passwordHash string
		err          error
	)
	row := this.db.QueryRow("SELECT id, password FROM users WHERE username = $1  LIMIT 1 ", username)
	row.Scan(&userId, &passwordHash)
	// Check if username is valid
	if userId == 0 {
		return nil, errors.New("Invalid Username")
	}
	userDetails := this.userService.Get(userId)
	if userDetails.Active == false {
		return nil, errors.New("Your account is not active , contact support ")
	}
	// Validates password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid Password")
	}
	// Check if two authentication is required
	if userDetails.TwoFactorEnabled {
		return this.twoFactorRequest(*userDetails, ipAddress, userAgent)
	}
	// Get user roles
	roles, err := this.userService.GetRoles(userId)
	userDetails.Roles = roles
	if err != nil {
		log.Println(err)
		return nil, errors.New("Failed to get roles")
	}
	return this.generateTokenDetails(*userDetails, ipAddress, userAgent)
}

// Refresh Token generates a new refresh token that will be used to get a new access token and a refresh token
func (this *AuthService) GenerateRefreshToken(oldRefreshToken, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId int
	)
	authResult := &models.AuthenticationResponse{}
	row := this.db.QueryRow("SELECT user_id FROM user_refresh_tokens WHERE token = $1 AND expiry_time > NOW() ", oldRefreshToken)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Refresh Token is not there")
		return nil, errors.New("Refresh Token is Invalid")
	}
	// Check if account is active before refreshing token
	userDetails := this.userService.Get(userId)
	if userDetails.Active == false {
		return nil, errors.New("Your account is not active , contact support ")
	}
	roles, _ := this.userService.GetRoles(userId)
	tokenExpiry := time.Duration(this.tokenTime)

	jwtToken, err := utilities.GenerateJwtToken(userId, roles, time.Duration(tokenExpiry))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Delete the old token and generate new access token and refresh token
	refreshToken := utilities.GenerateOpaqueToken(45)
	tx, err := this.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = tx.Exec("DELETE FROM user_refresh_tokens WHERE user_id = $1 AND token = $2", userId, oldRefreshToken)
	if err != nil {
		log.Println(err)
		tx.Rollback()
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
	_, err = tx.Exec(queryString, userId, refreshToken, ipAddress, userAgent, time.Now().Add(tokenExpiry))
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return nil, err
	}
	tx.Commit()
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = roles
	authResult.Expires = int(tokenExpiry.Seconds())
	return authResult, nil
}

func (this *AuthService) ResetPasswordRequest(username string) (bool, error) {
	//check if the username exists and then send vertification code
	var (
		userId                      int
		passwordResetTemplateBuffer bytes.Buffer
	)
	row := this.db.QueryRow("SELECT id FROM users where username = $1 OR email_address = $1 ", username)
	err := row.Scan(&userId)
	if err != nil {
		log.Println(err)
		return false, errors.New("Username or Email Address does not exists")
	}
	userDetails := this.userService.Get(userId)
	if userDetails.Active == false {
		log.Println(err)
		return false, errors.New("User is not active so password cannot be reset")
	}
	tx, err := this.db.Begin()
	if err != nil {
		log.Println(err)
		return false, err
	}
	// Delete any requested password requests
	_, err = tx.Exec("DELETE FROM reset_password_requests WHERE user_id = $1", userId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false, err
	}
	// Generate some random code to be sent to the user for reseting of the password
	rand.Seed(time.Now().UnixNano())
	randNumbers := make([]string, 6)
	for i := range randNumbers {
		randNumbers[i] = strconv.Itoa(rand.Intn(9))
	}
	randomNumberString := strings.Join(randNumbers, "")
	_, err = tx.Exec("INSERT INTO reset_password_requests(user_id, code, created, expiry_time) values($1, $2, NOW(), $3)", userId, randomNumberString, time.Now().Add(30*time.Minute))
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return false, err
	}
	// Get email template from directory and assign random code to it
	emailTemplateFile, err := template.ParseFiles("static/email_templates/PasswordRequest.html")
	if err != nil {
		log.Println("Email Error", err)
		return false, err
	}
	tmpl := template.Must(emailTemplateFile, err)
	emailTemplateData := struct {
		FullName   string
		RandomCode string
	}{}
	emailTemplateData.RandomCode = randomNumberString
	emailTemplateData.FullName = userDetails.FirstName + " " + userDetails.LastName
	tmpl.Execute(&passwordResetTemplateBuffer, emailTemplateData)
	recipient := []string{userDetails.EmailAddress}
	err = this.emailService.SendEmail(recipient, "Password Reset Request", passwordResetTemplateBuffer.String())
	if err != nil {
		log.Println("Email Error", err)
		tx.Rollback()
		return false, err
	}
	tx.Commit()
	return true, nil
}

// VerifyAndSetNewPassword functions to verify and reset password
func (this *AuthService) VerifyAndSetNewPassword(code string, password string) (bool, error) {
	// Check and see if code exists
	const errorMessage string = "Failed to Update password"
	var userId int
	tx, err := this.db.Begin()
	if err != nil {
		log.Println(err)
		return false, err
	}
	row := tx.QueryRow("SELECT user_id FROM reset_password_requests WHERE code = $1 AND expiry_time >= NOW()", code)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Invalid Code")
		return false, errors.New("Code is invalid")
	}
	// update password and delete all refresh tokens
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return false, errors.New(errorMessage)
	}
	_, err = tx.Exec("UPDATE users SET password = $2 WHERE id = $1", userId, passwordHash)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false, errors.New(errorMessage)
	}
	_, err = tx.Exec("DELETE FROM user_refresh_tokens WHERE user_id = $1", userId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false, errors.New(errorMessage)
	}
	_, err = tx.Exec("DELETE FROM reset_password_requests WHERE user_id = $1", userId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false, errors.New(errorMessage)
	}
	tx.Commit()
	return true, nil
}

func (this *AuthService) twoFactorRequest(userDetails models.User, ipAddress string, userAgent string) (*models.AuthenticationResponse, error) {
	var twoFactorRequestTemplateBuffer bytes.Buffer
	tx, err := this.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	authResult := &models.AuthenticationResponse{}
	rand.Seed(time.Now().UnixNano())
	randNumbers := make([]string, 6)
	for i := range randNumbers {
		randNumbers[i] = strconv.Itoa(rand.Intn(9))
	}
	// Expire after 5minutes
	expires := time.Duration(300 * time.Second)
	randomCode := strings.Join(randNumbers, "")
	requestId := utilities.GenerateOpaqueToken(60)
	queryString :=
		`INSERT 
		    INTO two_factor_requests 
            	(user_id, request_id, ip_address, code, user_agent, created_at, expiry_time)
	        VALUES
	        	($1, $2 ,$3 ,$4, $5, NOW(), $6)
	    `
	_, err = tx.Exec(queryString, userDetails.ID, requestId, ipAddress, randomCode, userAgent, time.Now().Add(expires))
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, errors.New("Error Generating Two factor request")
	}
	// Get email template from directory and assign random code to it
	emailTemplateFile, err := template.ParseFiles("static/email_templates/TwoFactorLogin.html")
	if err != nil {
		tx.Rollback()
		log.Println("Email Error", err)
		return nil, err
	}
	tmpl := template.Must(emailTemplateFile, err)
	emailTemplateData := struct {
		FullName   string
		RandomCode string
	}{}
	emailTemplateData.RandomCode = randomCode
	emailTemplateData.FullName = userDetails.FirstName + " " + userDetails.LastName
	tmpl.Execute(&twoFactorRequestTemplateBuffer, emailTemplateData)
	recipient := []string{userDetails.EmailAddress}
	err = this.emailService.SendEmail(recipient, "Two-factor login", twoFactorRequestTemplateBuffer.String())
	if err != nil {
		log.Println("Email Error", err)
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	authResult.TwoFactorEnabled = true
	authResult.Token = requestId
	return authResult, nil
}

func (this *AuthService) generateTokenDetails(userDetails models.User, ipAddress string, userAgent string) (*models.AuthenticationResponse, error) {
	authResult := &models.AuthenticationResponse{}
	tokenExpiry := time.Duration(this.tokenTime)
	// Generates JWT Token and Refresh token that expires after xminutes
	jwtToken, err := utilities.GenerateJwtToken(userDetails.ID, userDetails.Roles, tokenExpiry)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error Generating Access Token")
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
	_, err = this.db.Exec(queryString, userDetails.ID, refreshToken, ipAddress, userAgent, time.Now().Add(tokenExpiry))
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error Generating Refresh Token")
	}
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = userDetails.Roles
	authResult.Expires = int(tokenExpiry.Seconds())
	authResult.TwoFactorEnabled = userDetails.TwoFactorEnabled
	return authResult, nil
}

// Validate the two factor authentication request and complete the authentication request
func (this *AuthService) ValidateTwoFactor(code, requestId string, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var userId int
	row := this.db.QueryRow("SELECT user_id FROM two_factor_requests WHERE code = $1 AND request_id = $2 AND expiry_time > NOW() ", code, requestId)
	row.Scan(&userId)
	if userId == 0 {
		log.Println("Invalid Code")
		return nil, errors.New("Failed to validate")
	}
	if _, err := this.db.Exec("DELETE FROM two_factor_requests WHERE code = $1 AND request_id = $2", code, requestId); err != nil {
		log.Println(err)
		return nil, errors.New("Failed to validate")
	}
	userDetails := this.userService.Get(userId)
	return this.generateTokenDetails(*userDetails, ipAddress, userAgent)

}

// Delete expired tokens
func (this *AuthService) DeleteExpiredTokens(days int) error {
	// Deletes User Refresh tokens
	result, err := this.db.Exec("DELETE FROM user_refresh_tokens WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ := result.RowsAffected()
	log.Println("DELETED number of rows for user expired tokens :", count)
	// Deletes Two factor requests
	result, err = this.db.Exec("DELETE FROM two_factor_requests WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ = result.RowsAffected()
	log.Println("DELETED number of rows for two_factor_requests tokens :", count)
	// Delete Reset Password Requests
	result, err = this.db.Exec("DELETE FROM reset_password_requests WHERE (DATE_PART('day', AGE(NOW()::date ,expiry_time::date))) >= $1", days)
	count, _ = result.RowsAffected()
	log.Println("DELETED number of rows for reset_password_requests tokens :", count)
	return err
}
