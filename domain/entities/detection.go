package entities

import (
	"time"

	"github.com/google/uuid"
)

type Detection struct {
	ID         uint64    `db:"id"`
	Email      string    `db:"email"`
	IPAddress  string    `db:"ip_address"`
	CustomerID uuid.UUID `db:"customer_id"`
	Path       string    `db:"path"`
	Method     string    `db:"method"`
	StatusCode uint16    `db:"status_code"`
	IsDeepFake bool      `db:"is_deepfake"`
	// FileType   string    `db:"file_type"` RENCANA
	CreatedAt  time.Time `db:"created_at"`
	DeletedAt  time.Time `db:"deleted_at"`

	Customer Customer `db:"-"`
}
