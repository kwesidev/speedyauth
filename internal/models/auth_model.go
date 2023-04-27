package models

type AuthenticationRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
type AuthenticationResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	Roles        []string `json:"roles"`
	Expires      int      `json:"expires"`
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
