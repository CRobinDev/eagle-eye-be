package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `db:"id"`
	Username     string     `db:"username"`
	Email        string     `db:"email"`
	Password     string     `db:"password_hash"`
	Photo        string     `db:"photo"`
	Role         UserRole   `db:"role"`
	IsVerified   bool       `db:"is_verified"`
	IsCustomer   bool       `db:"is_customer"`
	AuthProvider string     `db:"auth_provider"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)
