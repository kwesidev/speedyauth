package services

import (
	"database/sql"
	"log"

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
			users.username,
			users.first_name,
			users.last_name,
			users.email_address,
			users.phone_number,
			users.active
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
		rows.Scan(&user.ID, &user.Username, &user.FirstName,
			&user.LastName, &user.CellNumber,
			&user.EmailAddress, &user.Active)
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
			users.username,
			users.first_name,
			users.last_name,
			users.email_address,
			users.phone_number,
			users.active
		FROM 
			users 
		WHERE 
			users.id = $1      
		LIMIT 1
        `
	row := this.db.QueryRow(queryString, userId)
	// Inject the data into the struct
	err := row.Scan(&userDetails.ID, &userDetails.Username, &userDetails.FirstName,
		&userDetails.LastName, &userDetails.EmailAddress, &userDetails.CellNumber, &userDetails.Active)
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
	queryString := `
    	INSERT INTO users
		    (username, password, first_name,last_name, email_address, phone_number, active)
		VALUES
			($1, $2, $3, $4, $5, $6, true) 
		RETURNING id ;`

	row := tx.QueryRow(queryString, userRegistrationRequest.Username, string(passwordHash),
		userRegistrationRequest.FirstName, userRegistrationRequest.LastName,
		userRegistrationRequest.EmailAddress, userRegistrationRequest.CellNumber)
	var newUserId int
	err = row.Scan(&newUserId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
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
		tx.Rollback()
		log.Println(err)
		return false, nil
	}
	tx.Commit()
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
