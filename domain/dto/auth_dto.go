package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
}
