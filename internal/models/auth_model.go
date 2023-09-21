package models

type AuthenticationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
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
	TwoFactorType    string   `json:"twoFactorType,omitempty"`
}

type VerifyTwoFactorRequest struct {
	Type  string `json:"type" validate:"required"`
	Token string `json:"token" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

type GeneralErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Status       int    `json:"status"`
}

type PasswordResetRequest struct {
	Username string `json:"username" validate:"required"`
}

type VerifyChangePasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}
