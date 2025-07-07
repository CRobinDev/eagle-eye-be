package entities

import "github.com/google/uuid"

type Payment struct {
	UserID  uuid.UUID `db:"user_id"`
	OrderID string    `db:"order_id"`
	Amount  float64   `db:"amount"`
	Tier    string    `db:"tier"`
	Status  string    `db:"status"`
	// PaymentType string  `db:"payment_type"`
	CreatedAt string  `db:"created_at"`
	UpdatedAt string  `db:"updated_at"`
	DeletedAt *string `db:"deleted_at"`

	Client User `db:"-"`
}
