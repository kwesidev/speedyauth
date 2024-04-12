package services

import "errors"

var (
	ErrUserNameExists   = errors.New("the username exists")
	ErrSendingMail      = errors.New("failed sending email")
	ErrAccountNotActive = errors.New("account is not active")
	ErrTokenGeneration  = errors.New("failed to generate token")
	ErrInvalidToken     = errors.New("token is invalid")
	ErrAccessToken      = errors.New("failed to access token")
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrRegistration     = errors.New("failed to register ")
	ErrPasswordUpdate   = errors.New("failed to update password")
	ErrTwoFactorCode    = errors.New("failed to verify two factor code")
	ErrTwoFactorRequest = errors.New("failed to send two factor request")
	ErrInvalidCode      = errors.New("code is invalid")
	ErrServer           = errors.New("server error, try again later")
	ErrPassCode         = errors.New("invalid Passcode")
	ErrStrongPassword   = errors.New("password must be at least 8 characters and must contain special characters")
	ErrTOTPExists       = errors.New("totp already enabled ")
)
