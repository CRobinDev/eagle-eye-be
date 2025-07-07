package interfaces

import (
	"context"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/google/uuid"
)

type IPaymentService interface {
	CreatePayment(ctx context.Context, req dto.PaymentRequest) (dto.PaymentResponse, error)
	UpdatePaymentStatus(ctx context.Context, req dto.UpdatePaymentStatusRequest) (dto.PaymentResponse, error)
	LatestPaymentStatus(ctx context.Context, req dto.GetPaymentStatusRequest) (dto.PaymentResponse, error)
}

type IPaymentRepository interface {
	CreatePayment(ctx context.Context, payment *entities.Payment) error
	UpdatePaymentStatus(ctx context.Context, payment *entities.Payment) (error)
	GetStatusByOrderID(ctx context.Context, orderID string) (entities.Payment, error)
	GetStatusByUserID(ctx context.Context, userID uuid.UUID) (entities.Payment, error)
}
