package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Detection struct {
	ID         uint64       `db:"id"`
	IPAddress  string       `db:"ip_address"`
	CustomerID uuid.UUID    `db:"customer_id"`
	Path       string       `db:"path"`
	Method     string       `db:"method"`
	StatusCode uint16       `db:"status_code"`
	Type       string       `db:"type"`
	Confidence float32      `db:"confidence"`
	IsBanned   bool         `db:"is_banned"`
	IsDeepFake bool         `db:"is_deepfake"`
	CreatedAt  time.Time    `db:"created_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
	TotalCount uint64       `db:"total_count"`

	Customer Customer `db:"-"`
}
