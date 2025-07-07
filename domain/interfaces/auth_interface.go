package interfaces

import (
	"context"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/google/uuid"
)

type IAuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (dto.UserResponse, error)
	// VerificationEmail(ctx context.Context, token string) error
	// RefreshToken(ctx context.Context, refreshToken string) (dto.LoginResponse, error)
}

type IAuthRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (entities.User, error)
	// UpdateVerificationStatus(ctx context.Context, status bool) error
}
