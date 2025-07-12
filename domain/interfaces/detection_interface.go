package interfaces

import (
	"context"
	"time"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/google/uuid"
)

type IDetectionService interface {
	DetectDeepFakeImage(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error)
	DetectDeepFakeAudio(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error)
	GetDetectionData(ctx context.Context, req dto.GetDetectionRequest) (dto.GetDetectionResponse, error)
	GetDeepFakeDetected(ctx context.Context, req dto.GetDetectionRequest) (dto.GetDetectionResponse, error)
	GetDetectionByID(ctx context.Context, req dto.GetDetectionDetailsRequest) (dto.DetectionDataResponse, error)
	BlockDetection(ctx context.Context, req dto.DeleteDetectionRequest) error
	UnblockDetection(ctx context.Context, req dto.UndeleteDetectionRequest) error
	GetCustomerUsage(ctx context.Context, req dto.GetCustomerUsageRequest) (dto.CustomerUsageResponse, error)
}

type IDetectionRepository interface {
	CreateDetection(ctx context.Context, detection *entities.Detection) error
	GetDetection(ctx context.Context, filter dto.GetDetectionFilter) ([]entities.Detection, error)
	GetDetectionByID(ctx context.Context, id uint16) (entities.Detection, error)
	GetDeepFakeDetected(ctx context.Context, filter dto.GetDetectionFilter) ([]entities.Detection, error)
	DeleteDetection(ctx context.Context, id uint16) error
	UndeleteDetection(ctx context.Context, id uint16) error
	GetHourlyUsage(ctx context.Context, customerID uuid.UUID, date time.Time) ([]dto.CustomerUsage, error)
	GetDailyUsage(ctx context.Context, customerID uuid.UUID, days int) ([]dto.CustomerUsage, error)
	GetWeeklyUsage(ctx context.Context, customerID uuid.UUID, weeks int) ([]dto.CustomerUsage, error)
	GetMonthlyUsage(ctx context.Context, customerID uuid.UUID, weeks int) ([]dto.CustomerUsage, error)
}
