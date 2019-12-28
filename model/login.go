package model

// LoginRequest api login request structure
type LoginRequest struct {
	ID       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse api login success response
type LoginResponse struct {
	PassphraseID int64  `json:"passphrase_id"`
	UserID       int64  `json:"user_id"`
	Passphrase   string `json:"passphrase"`
}

// TokenRequest api token request structure
type TokenRequest struct {
	Passphrase string `json:"passphrase" validate:"required"`
}
