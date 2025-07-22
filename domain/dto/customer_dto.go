package dto

import (
	"time"

	"github.com/google/uuid"
)

type GenerateAPIKeyRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Prefix string    `json:"prefix" validate:"required,min=8,max=50"`
}

type GenerateAPIKeyResponse struct {
	APIKey    string    `json:"api_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetCustomerRequest struct {
	CustomerID uuid.UUID `json:"customer_id"`
}

type GetCustomerAPIStatusResponse struct {
	IsCustomer bool      `json:"is_customer"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type GetCustomerTotalCallsResponse struct {
	TotalCalls uint64    `json:"total_calls"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type GetCustomerTierResponse struct {
	Tier string `json:"tier"`
}

type UpdateUsageRequest struct {
	Prefix string
}

type UpdateUsageResponse struct {
	ID           uuid.UUID
	CustomerTier string
	HashedKey    string
	Prefix       string
	CurrentUsage uint64
	MonthlyLimit uint32
	LastUsed     time.Time
}
