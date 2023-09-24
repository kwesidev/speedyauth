package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kwesidev/speedyauth/internal/models"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// List a bunch of users
func (usrSrv *UserService) List(offset int, limit int) ([]models.User, error) {
	users := []models.User{}
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

	return users, nil
}

// Get user details based on ID
func (usrSrv *UserService) Get(userId int) *models.User {
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
			users.two_factor_type,
			users.totp_secret ,
			users.totp_url
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
		&userDetails.TwoFactorType, &userDetails.TOTPSecret, &userDetails.TOTPURL,
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
func (usrSrv *UserService) GetByUsername(username string) *models.User {
	var userId int
	row := usrSrv.db.QueryRow("SELECT id FROM users WHERE username = $1 OR email_address = $2  LIMIT 1 ", username)
	row.Scan(&userId)
	return usrSrv.Get(userId)
}

// Register a new user
func (usrSrv *UserService) Register(userRegistrationRequest models.UserRegistrationRequest) (bool, error) {
	// Salt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userRegistrationRequest.Password), 10)
	if err != nil {
		return false, err
	}
	tx, err := usrSrv.db.Begin()
	defer tx.Rollback()
	queryString := `
    	INSERT INTO users
		    (username, password, first_name,last_name, email_address, phone_number, active, two_factor_enabled ,two_factor_type, totp_secret, totp_url)
		VALUES
			($1, $2, $3, $4, $5, $6, true, false, 'NONE','','') 
		RETURNING id ;`

	row := tx.QueryRow(queryString, userRegistrationRequest.Username, string(passwordHash),
		userRegistrationRequest.FirstName, userRegistrationRequest.LastName,
		userRegistrationRequest.EmailAddress, userRegistrationRequest.CellNumber)
	var newUserId int
	if err = row.Scan(&newUserId); err != nil {
		log.Println(err)
		return false, err
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
		return false, nil
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// GetRoles gets a list of user roles
func (usrSrv *UserService) GetRoles(userId int) ([]string, error) {
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
	return roles, nil

}

// Update User
func (usrSrv *UserService) Update(userId int, userUpdateRequest models.UserUpdateRequest) error {
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
		query += fmt.Sprintf("cell_number =$%d,", argCount)
		args = append(args, userUpdateRequest.CellNumber)
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
func (usrSrv *UserService) DeleteToken(userId int, refreshToken string) (bool, error) {
	if _, err := usrSrv.db.Exec("DELETE FROM user_refresh_tokens WHERE token = $1 AND user_id = $2", refreshToken, userId); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

// EnableTOTP
func (usrSrv *UserService) EnableTwoFactorTOTP(userId int) (*models.EnableTOTPResponse, error) {
	userDetails := usrSrv.Get(userId)
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
			two_factor_enabled = true , two_factor_type = 'TOTP',
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
func (usrSrv *UserService) EnableTwoFactor(userId int, typeCode string) error {
	queryString :=
		`UPDATE 
	        users 
		SET 
			two_factor_enabled = true , two_factor_type = $1,
			totp_secret = '', totp_url = '' , totp_created = NULL
	    WHERE 
		    id = $2
	    `
	if _, err := usrSrv.db.Exec(queryString, typeCode, userId); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
