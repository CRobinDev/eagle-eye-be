package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type DetectionRequest struct {
	File       *multipart.FileHeader `form:"file" validate:"required"`
	IP         string
	Path       string
	Method     string
	CustomerID uuid.UUID
}

type DetectionResponse struct {
	Filename   string `json:"filename"`
	Prediction string `json:"prediction"`
	Confidence float32 `json:"confidence"`
}

type GetDetectionFilter struct {
	CustomerID uuid.UUID
	Limit      uint64
	Offset     uint64
}

type DeleteDetectionRequest struct {
	DetectionID uint16 `json:"id" validate:"required"`
	UserID      uuid.UUID
}

type UndeleteDetectionRequest struct {
	DetectionID uint16 `json:"id" validate:"required"`
	UserID      uuid.UUID
}

type GetDetectionDetailsRequest struct {
	DetectionID uint16 `json:"id" validate:"required"`
	UserID      uuid.UUID
}

type GetDetectionRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	CurrentPage uint64    `json:"current_page" validate:"required"`
	Limit       uint64    `json:"limit" validate:"required"`
}

type GetDetectionResponse struct {
	CurrentPage uint64                  `json:"current_page"`
	TotalItems  []DetectionDataResponse `json:"total_items"`
	Limit       uint64                  `json:"limit"`
}

type DetectionDataResponse struct {
	ID          uint64    `json:"id"`
	IPAddress   string    `json:"ip_address"`
	Path        string    `json:"path"`
	Method      string    `json:"method"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"accessed_at"`
}

type GeminiAnalyzeRequest struct {
	File        []byte `json:"file"`
	ContentType string `json:"content_type"`
}

type GeminiAnalyzeResponse struct {
	Predict    string  `json:"prediction"`
	Confidence float32 `json:"confidence"`
}

type GetCustomerUsageRequest struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Mode       string    `json:"mode" validate:"required,oneof=hourly daily weekly monthly"`
	Date       string    `json:"date,omitempty"`   // untuk hourly
	Days       int       `json:"days,omitempty"`   // untuk daily
	Weeks      int       `json:"weeks,omitempty"`  // untuk weekly
	Months     int       `json:"months,omitempty"` // untuk monthly
}

type CustomerUsageResponse struct {
	Mode    string          `json:"mode" validate:"oneof=hourly daily weekly monthly"`
	Details []CustomerUsage `json:"details"`
}

type CustomerUsage struct {
	Time  string `json:"time"`
	Usage uint32 `json:"usage"`
}
