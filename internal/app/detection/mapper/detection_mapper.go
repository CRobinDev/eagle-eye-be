package mapper

import (
	"fmt"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
)

var (
	DetectionStatusMapper = map[bool]string{
		true:  "Ditolak",
		false: "Diterima",
	}

	DetectionTypeMapper = map[string]string{
		"audio":   "Suara",
		"image":   "Gambar",
		"unknown": "Tidak Diketahui",
	}
)

func ToDetectionResponse(detection entities.Detection) dto.DetectionDataResponse {
	resp := dto.DetectionDataResponse{
		ID:        detection.ID,
		IPAddress: detection.IPAddress,
		Path:      detection.Path,
		Method:    detection.Method,
		Status:    DetectionStatusMapper[detection.IsDeepFake],
		CreatedAt: detection.CreatedAt,
	}

	if detection.IsDeepFake {
		resp.Description = fmt.Sprintf("Ditolak pada API %s", DetectionTypeMapper[detection.Type])
	}

	return resp
}

func ToDetectionPaginationResponse(detections []entities.Detection, currentPage, limit uint64) dto.GetDetectionResponse {
	var items []dto.DetectionDataResponse
	for _, det := range detections {
		item := dto.DetectionDataResponse{
			ID:        det.ID,
			IPAddress: det.IPAddress,
			Path:      det.Path,
			Method:    det.Method,
			Status:    DetectionStatusMapper[det.IsDeepFake],
			CreatedAt: det.CreatedAt,
		}

		if det.IsDeepFake {
			item.Description = fmt.Sprintf("Ditolak pada API %s", DetectionTypeMapper[det.Type])
		}
		items = append(items, item)
	}

	return dto.GetDetectionResponse{
		CurrentPage: currentPage,
		TotalItems:  items,
		Limit:       limit,
	}
}
