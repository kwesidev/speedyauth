package services

import "errors"

var (
	ErrUserNameExists     = errors.New("The username exists")
	ErrSendingMail        = errors.New("Failed sending Email")
	ErrorAccountNotActive = errors.New("Account is not Active")
	ErrorTokenGeneration  = errors.New("Failed to generate Token")
	ErrorInvalidToken     = errors.New("Token is Invalid")
	ErrorAccessToken      = errors.New("Failed to Access Token")
	ErrorInvalidUsername  = errors.New("Invalid Username")
	ErrorInvalidPassword  = errors.New("Invalid Password")
	ErrorRegistration     = errors.New("Failed to register ")
	ErrorPasswordUpdate   = errors.New("Failed to update password")
	ErrorTwoFactorCode    = errors.New("Failed to Verify Two Factor Code")
	ErrorTwoFactorRequest = errors.New("Failed to Send Two Factor Request")
	ErrorInvalidCode      = errors.New("Code is invalid")
	ErrorServer           = errors.New("Server Error, Try again later")
)
