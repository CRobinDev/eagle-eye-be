package dto

import (
	"time"

	"github.com/google/uuid"
)

type GenerateAPIKeyRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	CustomerTier string    `json:"customer_tier" validate:"required,oneof=free basic premium"`
	Prefix       string    `json:"prefix" validate:"required,min=8,max=50"`
}

type GenerateAPIKeyResponse struct {
	APIKey    string    `json:"api_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetCustomerAPIStatusRequest struct {
	CustomerID uuid.UUID `json:"customer_id"`
}

type GetCustomerAPIStatusResponse struct {
	IsCustomer bool      `json:"is_customer"`
	ExpiresAt  time.Time `json:"expires_at"`
}
