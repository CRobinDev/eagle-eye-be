package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type DetectionRequest struct {
	File       *multipart.FileHeader `form:"file" validate:"required"`
	Email      string                `form:"email" validate:"required"`
	IP         string
	Path       string
	Method     string
	CustomerID uuid.UUID
}

type DetectionResponse struct {
	Filename   string `json:"filename"`
	Prediction string `json:"prediction" validate:"oneof=Real Fake"`
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
	ID         uint64    `json:"id"`
	IPAddress  string    `json:"ip_address"`
	Path       string    `json:"path"`
	Method     string    `json:"method"`
	StatusCode uint16    `json:"status_code"`
	CreatedAt  time.Time `json:"accessed_at"`
}

type GeminiAnalyzeRequest struct {
	File        []byte `json:"file"`
	ContentType string `json:"content_type"`
}

type GeminiAnalyzeResponse struct {
	Predict    string  `json:"prediction"`
	Confidence float32 `json:"confidence"`
}
