package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/app/detection/mapper"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/gemini"
	"github.com/CRobinDev/karsa/pkg/log"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type detectionService struct {
	dr     interfaces.IDetectionRepository
	gemini gemini.IGemini
	logger *logrus.Logger
}

func NewDetectionService(dr interfaces.IDetectionRepository, gemini gemini.IGemini, logger *logrus.Logger) interfaces.IDetectionService {
	return &detectionService{
		dr:     dr,
		gemini: gemini,
		logger: logger,
	}
}

func (ds *detectionService) DetectDeepFakeImage(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error) {
	traceID := utils.GetTraceID(ctx)
	body, writer, fileBytes, err := ds.createFormFile(traceID, req.File)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	resp, err := ds.createHttpRequest(traceID, body, writer, "image")
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	if resp.Prediction == "" {
		return dto.DetectionResponse{}, fmt.Errorf("detectResp struct empty. Failed to unmarshal : %v", fiber.StatusInternalServerError)
	}

	ds.logger.WithFields(log.WithTraceID(traceID, nil)).Infof("[DetectionService][DetectDeepFake] response from model: %v", resp)
	detection := entities.Detection{
		IPAddress:  req.IP,
		CustomerID: req.CustomerID,
		Path:       req.Path,
		Method:     req.Method,
		Type:       "image",
		StatusCode: uint16(200),
	}

	geminiReq := dto.GeminiAnalyzeRequest{
		File:        fileBytes,
		ContentType: req.File.Header.Get("Content-Type"),
		Predict:     strings.ToLower(resp.Prediction),
		Confidence:  resp.Confidence,
	}

	geminiResp, err := ds.gemini.DetectDeepFakeImage(ctx, geminiReq)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	ds.logger.WithFields(log.WithTraceID(traceID, nil)).Infof("[DetectionService][DetectDeepFake] response from gemini: %v", geminiResp)

	isDeepFake, prediction, confidence := ds.compareResult(geminiResp, resp)

	detection.IsDeepFake = isDeepFake
	detection.Confidence = confidence

	if err := ds.dr.CreateDetection(ctx, &detection); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create detection history")
		return dto.DetectionResponse{}, errorz.ErrFailedToCreateDetection.WithTraceID(traceID)
	}

	return dto.DetectionResponse{
		Filename:   resp.Filename,
		Prediction: prediction,
		Confidence: confidence,
	}, nil
}

func (ds *detectionService) DetectDeepFakeAudio(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error) {
	traceID := utils.GetTraceID(ctx)
	body, writer, fileBytes, err := ds.createFormFile(traceID, req.File)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	resp, err := ds.createHttpRequest(traceID, body, writer, "audio")
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	if resp.Prediction == "" {
		return dto.DetectionResponse{}, fmt.Errorf("detectResp struct empty. Failed to unmarshal : %v", fiber.StatusInternalServerError)
	}

	geminiReq := dto.GeminiAnalyzeRequest{
		File:        fileBytes,
		ContentType: req.File.Header.Get("Content-Type"),
		Predict:     strings.ToLower(resp.Prediction),
		Confidence:  resp.Confidence,
	}

	detection := entities.Detection{
		IPAddress:  req.IP,
		CustomerID: req.CustomerID,
		Path:       req.Path,
		Method:     req.Method,
		Type:       "audio",
		StatusCode: uint16(200),
	}

	geminiResp, err := ds.gemini.DetectDeepFakeAudio(ctx, geminiReq)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	isDeepFake, prediction, confidence := ds.compareResult(geminiResp, resp)

	detection.IsDeepFake = isDeepFake
	detection.Confidence = confidence

	if err := ds.dr.CreateDetection(ctx, &detection); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create detection history")
		return dto.DetectionResponse{}, errorz.ErrFailedToCreateDetection.WithTraceID(traceID)
	}

	return dto.DetectionResponse{
		Filename:   resp.Filename,
		Prediction: prediction,
		Confidence: confidence,
	}, nil
}

func (ds *detectionService) GetDetectionData(ctx context.Context, req dto.GetDetectionRequest) (dto.GetDetectionResponse, error) {
	traceID := utils.GetTraceID(ctx)
	detection := dto.GetDetectionFilter{
		Offset: req.Limit * (req.CurrentPage - 1),
		Limit:  req.Limit,
	}

	if req.UserID != uuid.Nil {
		detection.CustomerID = req.UserID
	}

	detections, err := ds.dr.GetDetection(ctx, detection)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][GetDetectionData] failed to get detection data")
		return dto.GetDetectionResponse{}, errorz.ErrSaveDetection.WithTraceID(traceID)
	}

	if len(detections) == 0 {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][GetDetectionData] no detection found")
		return dto.GetDetectionResponse{
			CurrentPage: req.CurrentPage,
			Items:       []dto.DetectionDataResponse{},
			Limit:       req.Limit,
		}, nil
	}

	return mapper.ToDetectionPaginationResponse(detections, req.CurrentPage, req.Limit), nil
}

func (ds *detectionService) GetDetectionByID(ctx context.Context, req dto.GetDetectionDetailsRequest) (dto.DetectionDataResponse, error) {
	detection, err := ds.dr.GetDetectionByID(ctx, req.DetectionID)
	if err != nil {
		return dto.DetectionDataResponse{}, err
	}

	return mapper.ToDetectionResponse(detection), nil
}

func (ds *detectionService) GetDetectionByIP(ctx context.Context, req dto.ValidateIPRequest) (bool, error) {
	detection, err := ds.dr.GetDetectionByIP(ctx, req.IP)
	if err != nil {
		return detection.IsBanned, err
	}

	return detection.IsBanned, nil
}

func (ds *detectionService) GetDeepFakeDetected(ctx context.Context, req dto.GetDetectionRequest) (dto.GetDetectionResponse, error) {
	traceID := utils.GetTraceID(ctx)
	detection := dto.GetDetectionFilter{
		Offset: req.Limit * (req.CurrentPage - 1),
		Limit:  req.Limit,
	}

	if req.UserID != uuid.Nil {
		detection.CustomerID = req.UserID
	}

	detections, err := ds.dr.GetDeepFakeDetected(ctx, detection)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][GetDetectionData] failed to get detection data")
		return dto.GetDetectionResponse{}, errorz.ErrSaveDetection.WithTraceID(traceID)
	}

	if len(detections) == 0 {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Warn("[DetectionService][GetDetectionData] no detection found")
		return dto.GetDetectionResponse{
			CurrentPage: req.CurrentPage,
			Items:       []dto.DetectionDataResponse{},
			Limit:       req.Limit,
		}, nil
	}

	return mapper.ToDetectionPaginationResponse(detections, req.CurrentPage, req.Limit), nil
}

func (ds *detectionService) BlockDetection(ctx context.Context, req dto.DeleteDetectionRequest) error {
	traceID := utils.GetTraceID(ctx)

	if err := ds.dr.DeleteDetection(ctx, req.IP); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DeleteDetection] failed to ban IP")
		return errorz.ErrBanIP.WithTraceID(traceID)
	}

	return nil
}

func (ds *detectionService) UnblockDetection(ctx context.Context, req dto.UndeleteDetectionRequest) error {
	traceID := utils.GetTraceID(ctx)

	if err := ds.dr.UndeleteDetection(ctx, req.IP); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DeleteDetection] failed to unban IP")
		return errorz.ErrBanIP.WithTraceID(traceID)
	}

	return nil
}

func (ds *detectionService) GetCustomerUsage(ctx context.Context, req dto.GetCustomerUsageRequest) (dto.CustomerUsageResponse, error) {
	emptyCustomerUsage := make([]dto.CustomerUsage, 0)
	loc := utils.GetTimeLocation()
	switch req.Mode {
	case "hourly":
		date, err := time.ParseInLocation("2006-01-02", req.Date, loc)
		if err != nil {
			return dto.CustomerUsageResponse{}, errors.New("invalid date format, use YYYY-MM-DD")
		}

		now := time.Now().In(loc)

		dateWithCurrentTime := time.Date(
			date.Year(), date.Month(), date.Day(),
			now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), loc,
		)

		resp, err := ds.dr.GetHourlyUsage(ctx, req.CustomerID, date)
		if err != nil {
			return dto.CustomerUsageResponse{}, err
		}

		var customerUsage []dto.CustomerUsage
		if resp != nil {
			customerUsage = mapper.ToCustomerUsageHourlyResponse(resp, dateWithCurrentTime)
		}

		return dto.CustomerUsageResponse{
			Mode:    req.Mode,
			Details: customerUsage,
		}, nil

	case "daily":
		if req.Days <= 0 {
			req.Days = 7
		}

		resp, err := ds.dr.GetDailyUsage(ctx, req.CustomerID, req.Days)
		if len(resp) == 0 {
			resp = emptyCustomerUsage
		} else if err != nil {
			return dto.CustomerUsageResponse{}, errors.New("failed to get daily usage")
		}

		return dto.CustomerUsageResponse{
			Mode:    req.Mode,
			Details: mapper.ToCustomerUsageDailyResponse(resp, req.Days),
		}, nil

	case "weekly":
		if req.Weeks <= 0 {
			req.Weeks = 4
		}
		resp, err := ds.dr.GetWeeklyUsage(ctx, req.CustomerID, req.Weeks)
		if len(resp) == 0 {
			resp = emptyCustomerUsage
		} else if err != nil {
			return dto.CustomerUsageResponse{}, errors.New("failed to get weekly usage")
		}

		return dto.CustomerUsageResponse{
			Mode:    req.Mode,
			Details: mapper.ToCustomerUsageWeeklyResponse(resp, req.Weeks),
		}, nil

	case "monthly":
		if req.Months <= 0 {
			req.Months = 6
		}
		resp, err := ds.dr.GetMonthlyUsage(ctx, req.CustomerID, req.Months)
		if len(resp) == 0 {
			resp = emptyCustomerUsage
		} else if err != nil {
			return dto.CustomerUsageResponse{}, errors.New("failed to get monthly usage")
		}

		return dto.CustomerUsageResponse{
			Mode:    req.Mode,
			Details: resp,
		}, nil

	default:
		return dto.CustomerUsageResponse{Details: emptyCustomerUsage}, errors.New("invalid mode")
	}
}

func (ds *detectionService) compareResult(geminiPredict dto.GeminiAnalyzeResponse, modelPredict dto.DetectionResponse) (bool, string, float32) {
	var isDeepFake bool
	var prediction string

	modelWeight := 0.6
	geminiWeight := 0.4

	geminiNumeric := getNumericLabel(geminiPredict.Predict)
	modelNumeric := getNumericLabel(modelPredict.Prediction)

	modelScore := float64(modelNumeric) * modelWeight * float64(modelPredict.Confidence)
	geminiScore := float64(geminiNumeric) * geminiWeight * float64(geminiPredict.Confidence)

	finalScore := modelScore + geminiScore

	if finalScore >= 0.5 {
		prediction = "fake"
		isDeepFake = true
	} else {
		prediction = "real"
		isDeepFake = false
	}

	confidence := (geminiPredict.Confidence + modelPredict.Confidence) / 2

	return isDeepFake, prediction, confidence
}

func getNumericLabel(pred string) int {
	if strings.ToLower(pred) == "fake" {
		return 1
	}
	return 0
}

func (ds *detectionService) createFormFile(traceID uuid.UUID, image *multipart.FileHeader) (*bytes.Buffer, *multipart.Writer, []byte, error) {
	file, err := image.Open()
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to open image")
		return nil, nil, nil, errorz.ErrFailedToOpenFile.WithTraceID(traceID)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to read image")
		return nil, nil, nil, errorz.ErrFailedToReadFile.WithTraceID(traceID)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", image.Filename)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create form file")
		return nil, nil, nil, errorz.ErrFailedToCreateFormFile.WithTraceID(traceID)
	}

	_, err = part.Write(fileBytes)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] to write bytes")
		return nil, nil, nil, errorz.ErrFailedToWriteBytes.WithTraceID(traceID)
	}

	if err := writer.Close(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to close writer : %v", err)
	}

	return body, writer, fileBytes, nil
}

func (ds *detectionService) createHttpRequest(traceID uuid.UUID, body *bytes.Buffer, writer *multipart.Writer, types string) (dto.DetectionResponse, error) {
	var url string
	if types == "audio" {
		url = env.GetEnv().AudioModelUrl
	} else {
		url = env.GetEnv().ImageModelUrl
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create http request")
		return dto.DetectionResponse{}, errorz.ErrFailedToCreateHTTPRequest.WithTraceID(traceID)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to make http client")
		return dto.DetectionResponse{}, errorz.ErrFailedToInitiateHTTPClient.WithTraceID(traceID)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to analyze image")
		return dto.DetectionResponse{}, errorz.ErrFailedToAnalyzeImage.WithTraceID(traceID)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to read response body")
		return dto.DetectionResponse{}, errorz.ErrFailedToReadFile.WithTraceID(traceID)
	}

	var detectResp dto.DetectionResponse
	if err := json.Unmarshal(bodyResp, &detectResp); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to unmarshal json to struct")
		return dto.DetectionResponse{}, errorz.ErrUnmarshal.WithTraceID(traceID)
	}

	return detectResp, nil
}
