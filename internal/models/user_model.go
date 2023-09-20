package models

type User struct {
	ID               int      `json:"-"`
	UUID             string   `json:"id"`
	Username         string   `json:"username"`
	EmailAddress     string   `json:"emailAddress"`
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	CellNumber       string   `json:"cellNumber"`
	Roles            []string `json:"roles"`
	Active           bool     `json:"active"`
	TwoFactorEnabled bool     `json:"twoFactorEnabled"`
	TwoFactorType    string   `json:"twoFactorType"`
	TOTPSecret       string
	TOTPURL          string
}
type EnableTwoFactorRequest struct {
	Type string `json:"type"`
}
type EnableTOTPResponse struct {
	URL string `json:"url"`
}
type VerifyPassCodeRequest struct {
	Code string `json:"code"`
}
type UserRegistrationRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	EmailAddress string `json:"emailAddress" validate:"required"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	CellNumber   string `json:"cellNumber" validate:"required"`
}

type UserUpdateRequest struct {
	EmailAddress                 string `json:"emailAddress"`
	FirstName                    string `json:"firstName"`
	LastName                     string `json:"lastName"`
	CellNumber                   string `json:"cellNumber"`
	AllowTwoFactorAuthentication bool   `json:"allowTwoFactorAuthentication"`
}
