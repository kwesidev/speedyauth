package models

type User struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	EmailAddress string   `json:"emailAddress"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	CellNumber   string   `json:"cellNumber"`
	Roles        []string `json:"roles"`
	Active       bool     `json:"active"`
}

type UserRegistrationRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	EmailAddress string `json:"emailAddress" validate:"required"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	CellNumber   string `json:"cellNumber" validate:"required"`
}
