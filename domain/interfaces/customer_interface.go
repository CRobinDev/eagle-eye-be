package interfaces

import (
	"context"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/google/uuid"
)

type ICustomerService interface {
	GenerateAPIKey(ctx context.Context, req dto.GenerateAPIKeyRequest) (dto.GenerateAPIKeyResponse, error)
	UpdateAPIKey(ctx context.Context, req dto.GenerateAPIKeyRequest) (dto.GenerateAPIKeyResponse, error)
}

type ICustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *entities.Customer) error
	UpdateAPIKey(ctx context.Context, customer *entities.Customer) (entities.Customer, error)
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (entities.Customer, error)
	UpdateUsage(ctx context.Context, prefix string) (entities.Customer, error)
}
