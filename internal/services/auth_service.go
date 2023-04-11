package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kwesidev/authserver/internal/models"
	"github.com/kwesidev/authserver/internal/utilities"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db          *sql.DB
	userService *UserService
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		db:          db,
		userService: NewUserService(db),
	}
}

// Login function to authenticate user
func (this *AuthService) Login(username, password, ipAddress, userAgent string) (*models.AuthenticationResponse, error) {
	var (
		userId       int
		passwordHash string
		err          error
	)
	authResult := &models.AuthenticationResponse{}
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
	// Get user roles
	roles, err := this.userService.GetRoles(userId)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Failed to get roles")
	}
	// Generates JWT Token and Refresh token that expires after 30minutes
	jwtToken, err := utilities.GenerateJwtToken(userId, roles, time.Duration(30*time.Minute))
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error Generating Access Token")
	}
	refreshToken := utilities.GenerateOpaqueToken(45)
	queryString :=
		`INSERT 
		    INTO 
	 	        user_refresh_tokens
                    (user_id, token, created, ip_address, user_agent, expiry_time)
	            VALUES
	                ($1, $2 ,NOW() ,$3 ,$4, $5)
	    `
	// Generate a jwt and refresh token
	_, err = this.db.Exec(queryString, userId, refreshToken, ipAddress, userAgent, time.Now().Add(30*time.Minute))
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error Generating Refresh Token")
	}
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = roles
	return authResult, nil
}

// DeleteToken function to delete refresh Token
func (this *AuthService) DeleteToken(refreshToken string) (bool, error) {
	_, err := this.db.Exec("DELETE FROM user_refresh_tokens WHERE refresh_token = $1", refreshToken)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
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
	jwtToken, err := utilities.GenerateJwtToken(userId, roles, time.Duration(30*time.Minute))
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
	_, err = this.db.Exec("DELETE FROM user_refresh_tokens WHERE user_id = $1 AND token = $2", userId, oldRefreshToken)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}
	queryString :=
		`INSERT 
		    INTO 
	 	        user_refresh_tokens
                    (user_id, token, created, ip_address, user_agent, expiry_time)
	            VALUES
	                ($1, $2 ,NOW() ,$3 ,$4, $5)
	    `
	// Generate a jwt and refresh token
	_, err = this.db.Exec(queryString, userId, refreshToken, ipAddress, userAgent, time.Now().Add(30*time.Minute))
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return nil, err
	}
	tx.Commit()
	authResult.RefreshToken = refreshToken
	authResult.Token = jwtToken
	authResult.Roles = roles
	return authResult, nil
}
