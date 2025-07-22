package mapper

import (
	"fmt"
	"math"
	"time"
	_ "time/tzdata"

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
		IsBanned:  detection.IsBanned,
	}

	if detection.IsDeepFake {
		resp.Description = fmt.Sprintf("Ditolak pada API %s", DetectionTypeMapper[detection.Type])
	}

	return resp
}

func ToDetectionPaginationResponse(detections []entities.Detection, currentPage, limit uint64) dto.GetDetectionResponse {
	var items []dto.DetectionDataResponse
	var totalCount uint64

	if len(detections) > 0 {
		totalCount = detections[0].TotalCount
	}

	for _, det := range detections {
		item := dto.DetectionDataResponse{
			ID:        det.ID,
			IPAddress: det.IPAddress,
			Path:      det.Path,
			Method:    det.Method,
			Status:    DetectionStatusMapper[det.IsDeepFake],
			CreatedAt: det.CreatedAt,
			IsBanned:  det.IsBanned,
		}

		if det.IsDeepFake {
			item.Description = fmt.Sprintf("Ditolak pada API %s", DetectionTypeMapper[det.Type])
		}
		items = append(items, item)
	}

	return dto.GetDetectionResponse{
		CurrentPage: currentPage,
		Items:       items,
		Limit:       limit,
		TotalPage:   math.Ceil(float64(totalCount) / float64(limit)),
	}
}

func ToCustomerUsageHourlyResponse(raw []dto.CustomerUsage, now time.Time) []dto.CustomerUsage {
	usageMap := make(map[string]uint32)
	for _, u := range raw {
		usageMap[u.Time] = u.Usage
	}

	currentHour := now.Hour()
	var result []dto.CustomerUsage

	for h := 0; h <= currentHour; h++ {
		label := fmt.Sprintf("%02d:00", h)
		usage := usageMap[label]
		result = append(result, dto.CustomerUsage{
			Time:  label,
			Usage: usage,
		})
	}

	return result
}

func ToCustomerUsageDailyResponse(raw []dto.CustomerUsage, days int) []dto.CustomerUsage {
	usageMap := make(map[string]int)
	for _, u := range raw {
		usageMap[u.Time] = int(u.Usage)
	}

	var result []dto.CustomerUsage

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		label := date.Format("2006-01-02")
		usage := usageMap[label]

		result = append(result, dto.CustomerUsage{
			Time:  label,
			Usage: uint32(usage),
		})
	}

	return result
}

func ToCustomerUsageWeeklyResponse(raw []dto.CustomerUsage, weeks int) []dto.CustomerUsage {
	usageMap := make(map[string]int)
	for _, u := range raw {
		usageMap[u.Time] = int(u.Usage)
	}

	var result []dto.CustomerUsage

	for i := (weeks * 7) - 2; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		label := date.Format("2006-01-02")
		usage := usageMap[label]

		result = append(result, dto.CustomerUsage{
			Time:  label,
			Usage: uint32(usage),
		})
	}

	return result
}
