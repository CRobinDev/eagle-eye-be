package dto

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UpdateUserRequest struct {
	UserID     uuid.UUID `json:"user_id,omitempty"`
	Username   string    `json:"username" validate:"omitempty,min=3,max=50"`
	IsCustomer bool
}

type UpdatePasswordRequest struct {
	UserID          uuid.UUID `json:"user_id,omitempty"`
	NewPassword     string    `json:"new_password" validate:"required,min=8,max=100"`
	ConfirmPassword string    `json:"confirm_password" validate:"required,min=8,max=100"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	IsCustomer bool      `json:"is_customer"`
}
