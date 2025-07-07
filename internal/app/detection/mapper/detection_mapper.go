package mapper

import (
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
)

func ToDetectionResponse(detections []entities.Detection, currentPage, limit uint64) dto.GetDetectionResponse {
	var items []dto.DetectionDataResponse
	for _, det := range detections {
		items = append(items, dto.DetectionDataResponse{
			ID:         det.ID,
			IPAddress:  det.IPAddress,
			Path:       det.Path,
			Method:     det.Method,
			StatusCode: det.StatusCode,
			CreatedAt:  det.CreatedAt,
		})
	}

	return dto.GetDetectionResponse{
		CurrentPage: currentPage,
		TotalItems:  items,
		Limit:       limit,
	}
}
