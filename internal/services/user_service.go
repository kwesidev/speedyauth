package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/kwesidev/authserver/internal/models"
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
func (this *UserService) List(offset int, limit int) ([]models.User, error) {
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
	rows, err := this.db.Query(queryString, offset, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.ID, &user.UUID, &user.Username, &user.FirstName,
			&user.LastName, &user.CellNumber,
			&user.EmailAddress, &user.Active, &user.TwoFactorEnabled)
		roles, _ := this.GetRoles(user.ID)
		user.Roles = roles
		users = append(users, user)
	}

	return users, nil
}

// Get user details based on ID
func (this *UserService) Get(userId int) *models.User {
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
			users.two_factor_enabled
		FROM 
			users 
		WHERE 
			users.id = $1      
		LIMIT 1
        `
	row := this.db.QueryRow(queryString, userId)
	// Inject the data into the struct
	err := row.Scan(&userDetails.ID, &userDetails.UUID, &userDetails.Username, &userDetails.FirstName,
		&userDetails.LastName, &userDetails.EmailAddress, &userDetails.CellNumber, &userDetails.Active, &userDetails.TwoFactorEnabled)
	roles, _ := this.GetRoles(userDetails.ID)
	userDetails.Roles = roles
	if err != nil {
		log.Println(err)
		return nil
	}
	return userDetails
}

// Register a new user
func (this *UserService) Register(userRegistrationRequest models.UserRegistrationRequest) (bool, error) {
	// Salt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userRegistrationRequest.Password), 10)
	if err != nil {
		return false, err
	}
	tx, err := this.db.Begin()
	defer tx.Rollback()
	queryString := `
    	INSERT INTO users
		    (username, password, first_name,last_name, email_address, phone_number, active, two_factor_enabled)
		VALUES
			($1, $2, $3, $4, $5, $6, true, false) 
		RETURNING id ;`

	row := tx.QueryRow(queryString, userRegistrationRequest.Username, string(passwordHash),
		userRegistrationRequest.FirstName, userRegistrationRequest.LastName,
		userRegistrationRequest.EmailAddress, userRegistrationRequest.CellNumber)
	var newUserId int
	err = row.Scan(&newUserId)
	if err != nil {
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
	_, err = tx.Exec(queryString, newUserId, "USER")
	if err != nil {
		log.Println(err)
		return false, nil
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

// GetRoles gets a list of user roles
func (this *UserService) GetRoles(userId int) ([]string, error) {
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
	rows, err := this.db.Query(queryString, userId)
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
func (this *UserService) Update(userId int, userUpdateRequest models.UserUpdateRequest) error {
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
	_, err := this.db.Exec(query, args...)
	if err != nil {
		log.Println("Updating user failed ", err)
		return err
	}
	return nil

}

// DeleteToken function to delete refresh Token
func (this *UserService) DeleteToken(userId int, refreshToken string) (bool, error) {
	_, err := this.db.Exec("DELETE FROM user_refresh_tokens WHERE token = $1 AND user_id = $2", refreshToken, userId)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
