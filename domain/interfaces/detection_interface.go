package interfaces

import (
	"context"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
)

type IDetectionService interface {
	DetectDeepFakeImage(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error)
	DetectDeepFakeAudio(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error)
	GetDetectionData(ctx context.Context, req dto.GetDetectionRequest) (dto.GetDetectionResponse, error)
	GetDeepFakeDetected(ctx context.Context, req dto.GetDetectionRequest) (dto.GetDetectionResponse, error)
	BlockDetection(ctx context.Context, req dto.DeleteDetectionRequest) error
	// DetectDeepFakeGrpc(ctx context.Context, req proto.DetectionRequest) (dto.DetectionResponse, error)
	// DetectDeepFakeVoice(ctx context.Context)
}

type IDetectionRepository interface {
	CreateDetection(ctx context.Context, detection *entities.Detection) error
	GetDetection(ctx context.Context, filter dto.GetDetectionFilter) ([]entities.Detection, error)
	GetDeepFakeDetected(ctx context.Context, filter dto.GetDetectionFilter) ([]entities.Detection, error)
	DeleteDetection(ctx context.Context, id uint16) error
}
