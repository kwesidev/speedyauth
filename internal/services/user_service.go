package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kwesidev/speedyauth/internal/models"
	"github.com/kwesidev/speedyauth/internal/utilities"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	List(offset int, limit int) ([]models.User, error)
	Get(userId int) *models.User
	GetByUsername(username string) *models.User
	GetRoles(userId int) ([]string, error)
	Update(userId int, userUpdateRequest models.UserUpdateRequest) error
	DeleteToken(userId int, refreshToken string) (bool, error)
	EnableTwoFactorTOTP(userId int) (*models.EnableTOTPResponse, error)
	EnableTwoFactor(userId int, methodCode string) error
	Register(userRegistrationRequest models.UserRegistrationRequest) (bool, error)
}

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &userService{
		db: db,
	}
}

// List a bunch of users
func (usrSrv *userService) List(offset int, limit int) ([]models.User, error) {
	users := make([]models.User, 0)
	// Get the list of users
	queryString :=
		`SELECT 
			users.id,
			users.uu_id,
			users.username,
			users.first_name,
			users.last_name,
			users.email_address,
			users.phone_number,
			users.active,
			users.two_factor_enabled
		FROM 
			users 
		OFFSET $1
		LIMIT $2     
        `
	rows, err := usrSrv.db.Query(queryString, offset, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.ID, &user.UUID, &user.Username, &user.FirstName,
			&user.LastName, &user.CellNumber,
			&user.EmailAddress, &user.Active, &user.TwoFactorEnabled)
		roles, _ := usrSrv.GetRoles(user.ID)
		user.Roles = roles
		users = append(users, user)
	}
	defer rows.Close()

	return users, nil
}

// Get user details based on ID
func (usrSrv *userService) Get(userId int) *models.User {
	userDetails := &models.User{}
	queryString :=
		`SELECT 
			users.id,
			users.uu_id,
			users.username,
			users.first_name,
			users.last_name,
			users.email_address,
			users.phone_number,
			users.active,
			users.two_factor_enabled,
			users.two_factor_method,
			users.totp_secret ,
			users.totp_url,
			users.meta_data
		FROM 
			users 
		WHERE 
			users.id = $1      
		LIMIT 1
        `
	row := usrSrv.db.QueryRow(queryString, userId)
	// Inject the data into the struct
	err := row.Scan(&userDetails.ID, &userDetails.UUID, &userDetails.Username, &userDetails.FirstName,
		&userDetails.LastName, &userDetails.EmailAddress, &userDetails.CellNumber, &userDetails.Active, &userDetails.TwoFactorEnabled,
		&userDetails.TwoFactorMethod, &userDetails.TOTPSecret, &userDetails.TOTPURL, &userDetails.Metadata,
	)
	roles, _ := usrSrv.GetRoles(userDetails.ID)
	userDetails.Roles = roles
	if err != nil {
		log.Println(err)
		return nil
	}
	return userDetails
}

// GetUsername gets the usersDetails by username
func (usrSrv *userService) GetByUsername(username string) *models.User {
	var userId int
	row := usrSrv.db.QueryRow("SELECT id FROM users WHERE username = $1 OR email_address = $1  LIMIT 1 ", username)
	row.Scan(&userId)
	return usrSrv.Get(userId)
}

// Register a new user
func (usrSrv *userService) Register(userRegistrationRequest models.UserRegistrationRequest) (bool, error) {
	// Check passwor strength
	if !utilities.StrongPasswordCheck(userRegistrationRequest.Password) {
		return false, ErrStrongPassword
	}
	// Salt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userRegistrationRequest.Password), 10)
	if err != nil {
		return false, err
	}
	tx, err := usrSrv.db.Begin()
	if err != nil {
		log.Println("Failed to start transaction")
		return false, ErrRegistration
	}
	defer tx.Rollback()
	queryString := `
    	INSERT INTO users
		    (username, password, first_name,last_name, email_address, phone_number, active, two_factor_enabled ,two_factor_method, totp_secret, totp_url)
		VALUES
			($1, $2, $3, $4, $5, $6, true, false, 'NONE','','') 
		RETURNING id ;`

	row := tx.QueryRow(queryString, userRegistrationRequest.Username, string(passwordHash),
		userRegistrationRequest.FirstName, userRegistrationRequest.LastName,
		userRegistrationRequest.EmailAddress, userRegistrationRequest.CellNumber)
	var newUserId int
	if err = row.Scan(&newUserId); err != nil {
		log.Println(err)
		return false, ErrRegistration
	}
	queryString = `
	    INSERT 
		    INTO user_roles
			(user_id, role_id) 
	    VALUES
		    ($1, (SELECT id FROM roles WHERE type = $2))
	    `
	if _, err = tx.Exec(queryString, newUserId, "USER"); err != nil {
		log.Println(err)
		return false, ErrRegistration
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// GetRoles gets a list of user roles
func (usrSrv *userService) GetRoles(userId int) ([]string, error) {
	roles := []string{}
	// Get user roles
	queryString := `
		SELECT 
			roles.type AS role_name
		FROM 
			user_roles
		LEFT JOIN 
			roles ON user_roles.role_id = roles.id 
		WHERE 
			user_roles.user_id = $1
	    `
	rows, err := usrSrv.db.Query(queryString, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role string
		rows.Scan(&role)
		roles = append(roles, role)
	}
	defer rows.Close()
	return roles, nil

}

// Update User
func (usrSrv *userService) Update(userId int, userUpdateRequest models.UserUpdateRequest) error {
	query := "UPDATE users SET "
	var args []any
	argCount := 1
	// Update first Name
	if strings.Trim(userUpdateRequest.FirstName, "") != "" {
		query += fmt.Sprintf("first_name = $%d, ", argCount)
		args = append(args, userUpdateRequest.FirstName)
		argCount++
	}
	// Update last Name
	if strings.Trim(userUpdateRequest.LastName, "") != "" {
		query += fmt.Sprintf("last_name = $%d, ", argCount)
		args = append(args, userUpdateRequest.LastName)
		argCount++
	}
	// Update email Address
	if strings.Trim(userUpdateRequest.EmailAddress, "") != "" {
		query += fmt.Sprintf("email_address = $%d, ", argCount)
		args = append(args, userUpdateRequest.EmailAddress)
		argCount++
	}
	// Update cell number
	if strings.Trim(userUpdateRequest.CellNumber, "") != "" {
		query += fmt.Sprintf("cell_number =$%d, ", argCount)
		args = append(args, userUpdateRequest.CellNumber)
		argCount++
	}
	// Update meta data
	if userUpdateRequest.Metadata != nil {
		query += fmt.Sprintf("meta_data = $%d, ", argCount)
		args = append(args, userUpdateRequest.Metadata)
		argCount++
	}
	// Remove the trailing comma and space
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, userId)

	// Execute the query with the dynamic arguments
	if _, err := usrSrv.db.Exec(query, args...); err != nil {
		log.Println("Updating user failed ", err)
		return err
	}
	return nil

}

// DeleteToken function to delete refresh Token
func (usrSrv *userService) DeleteToken(userId int, refreshToken string) (bool, error) {
	if _, err := usrSrv.db.Exec("DELETE FROM user_refresh_tokens WHERE token = $1 AND user_id = $2", refreshToken, userId); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

// EnableTOTP
func (usrSrv *userService) EnableTwoFactorTOTP(userId int) (*models.EnableTOTPResponse, error) {
	userDetails := usrSrv.Get(userId)
	if userDetails.TwoFactorMethod == "TOTP" {
		return nil, ErrTOTPExists
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      os.Getenv("ISSUER_NAME"),
		AccountName: userDetails.Username,
		SecretSize:  50,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	enableTOTPResponse := &models.EnableTOTPResponse{}
	enableTOTPResponse.URL = key.URL()
	queryString :=
		`UPDATE 
	        users 
		SET 
			two_factor_enabled = true , two_factor_method = 'TOTP',
			totp_secret = $1, totp_url = $2 , totp_created = NOW()
	    WHERE 
		    id = $3
	    `
	if _, err := usrSrv.db.Exec(queryString, key.Secret(), key.URL(), userId); err != nil {
		log.Println(err)
		return nil, err
	}

	return enableTOTPResponse, nil
}

// EnableTwoFactor SMS OR EMAIL
func (usrSrv *userService) EnableTwoFactor(userId int, methodCode string) error {
	queryString :=
		`UPDATE 
	        users 
		SET 
			two_factor_enabled = true , two_factor_method = $1,
			totp_secret = '', totp_url = '' , totp_created = NULL
	    WHERE 
		    id = $2
	    `
	if _, err := usrSrv.db.Exec(queryString, methodCode, userId); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
