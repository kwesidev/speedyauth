package models

type AuthenticationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PasswordLessAuthRequest struct {
	Username   string `json:"username"`
	SendMethod string `json:"sendMethod"`
}
type PasswordLessAuthResponse struct {
	RequestId  string `json:"requestId" validate:"required"`
	SendMethod string `json:"sendMethod" validate:"required"`
}
type CompletePasswordLessRequest struct {
	RequestId string `json:"requestId" validate:"required"`
	Code      string `json:"code" validate:"required"`
}
type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthenticationResponse struct {
	Token            string   `json:"token"`
	RefreshToken     string   `json:"refreshToken,omitempty"`
	Roles            []string `json:"roles,omitempty"`
	Expires          int      `json:"expiresIn,omitempty"`
	TwoFactorEnabled bool     `json:"twoFactorEnabled"`
	TwoFactorMethod  string   `json:"twoFactorMethod,omitempty"`
}

type VerifyTwoFactorRequest struct {
	Method string `json:"method" validate:"required"`
	Token  string `json:"token" validate:"required"`
	Code   string `json:"code" validate:"required"`
}

type GeneralErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Status       int    `json:"status"`
}
type SuccessResponse struct {
	Success bool `json:"success"`
}

type PasswordResetRequest struct {
	Username string `json:"username" validate:"required"`
}

type VerifyChangePasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}
