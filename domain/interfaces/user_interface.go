package interfaces

import (
	"context"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/google/uuid"
)

type IUserService interface {
	UpdateUser(ctx context.Context, req dto.UpdateUserRequest) error
	UpdatePassword(ctx context.Context, req dto.UpdatePasswordRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type IUserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	UpdatePassword(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
