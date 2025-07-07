package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetPaymentStatusRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type UpdatePaymentStatusRequest struct {
	Amount            string `json:"gross_amount" validate:"required"`
	OrderID           string `json:"order_id" validate:"required"`
	TransactionStatus string `json:"transaction_status" validate:"required"`
	FraudStatus       string `json:"fraud_status" validate:"required"`
	SignatureKey      string `json:"signature_key" validate:"required"`
	Code              string `json:"status_code" validate:"required"`
}

type PaymentRequest struct {
	UserID        uuid.UUID `json:"user_id"`
	TierOrder     string    `json:"tier_order" validate:"required"`
	OrderID       string    `json:"order_id"`
	Amount        int64     `json:"amount" validate:"required"`
	CustomerName  string    `json:"customer_name,omitempty"`
	CustomerEmail string    `json:"customer_email,omitempty"`
}

type PaymentResponse struct {
	SnapURL string `json:"snap_url,omitempty"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type PaymentEmailNotification struct {
	Fullname   string
	Email      string
	Tier       string
	Date       time.Time
	OrderID    string
	ExpiryDate time.Time
	Path       string
	Subject    string
}
