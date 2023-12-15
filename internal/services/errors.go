package services

import "errors"

var (
	ErrUserNameExists   = errors.New("The username exists")
	ErrSendingMail      = errors.New("Failed sending Email")
	ErrAccountNotActive = errors.New("Account is not Active")
	ErrTokenGeneration  = errors.New("Failed to generate Token")
	ErrInvalidToken     = errors.New("Token is Invalid")
	ErrAccessToken      = errors.New("Failed to Access Token")
	ErrInvalidUsername  = errors.New("Invalid Username")
	ErrInvalidPassword  = errors.New("Invalid Password")
	ErrRegistration     = errors.New("Failed to register ")
	ErrPasswordUpdate   = errors.New("Failed to update password")
	ErrTwoFactorCode    = errors.New("Failed to Verify Two Factor Code")
	ErrTwoFactorRequest = errors.New("Failed to Send Two Factor Request")
	ErrInvalidCode      = errors.New("Code is invalid")
	ErrServer           = errors.New("Server Error, Try again later")
	ErrPassCode         = errors.New("Invalid Passcode")
	ErrStrongPassword   = errors.New("Password must be at least 8 characters and must contain special characters")
	ErrTOTPExists       = errors.New("TOTP Already Enabled ")
)
