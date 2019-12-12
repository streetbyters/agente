package model

// LoginRequest api login request structure
type LoginRequest struct {
	ID		 string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
