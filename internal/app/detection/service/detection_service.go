package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

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
	// dsGrpc proto.ModelServiceClient
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

	resp, err := ds.createHttpRequest(traceID, body, writer)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	if resp.Prediction == "" {
		return dto.DetectionResponse{}, fmt.Errorf("detectResp struct empty. Failed to unmarshal : %v", fiber.StatusInternalServerError)
	}

	detection := entities.Detection{
		IPAddress:  req.IP,
		Email:      req.Email,
		CustomerID: req.CustomerID,
		Path:       req.Path,
		Method:     req.Method,
		StatusCode: uint16(200),
	}

	geminiReq := dto.GeminiAnalyzeRequest{
		File:        fileBytes,
		ContentType: req.File.Header.Get("Content-Type"),
	}

	geminiResp, err := ds.gemini.DetectDeepFakeImage(ctx, geminiReq)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	if geminiResp.Predict != "fake" || strings.ToLower(resp.Prediction) != "fake" {
		detection.IsDeepFake = false
	} else {
		detection.IsDeepFake = true
	}

	if err := ds.dr.CreateDetection(ctx, &detection); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create detection history")
		return dto.DetectionResponse{}, errorz.ErrFailedToCreateDetection.WithTraceID(traceID)
	}

	return resp, nil
}

func (ds *detectionService) DetectDeepFakeAudio(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error) {
	traceID := utils.GetTraceID(ctx)
	_, _, fileBytes, err := ds.createFormFile(traceID, req.File)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	// resp, err := ds.createHttpRequest(traceID, body, writer)
	// if err != nil {
	// 	return dto.DetectionResponse{}, err
	// }

	// if resp.Prediction == "" {
	// 	return dto.DetectionResponse{}, fmt.Errorf("detectResp struct empty. Failed to unmarshal : %v", fiber.StatusInternalServerError)
	// }

	geminiReq := dto.GeminiAnalyzeRequest{
		File:        fileBytes,
		ContentType: req.File.Header.Get("Content-Type"),
	}

	detection := entities.Detection{
		IPAddress:  req.IP,
		Email:      req.Email,
		CustomerID: req.CustomerID,
		Path:       req.Path,
		Method:     req.Method,
		StatusCode: uint16(200),
	}

	resp, err := ds.gemini.DetectDeepFakeAudio(ctx, geminiReq)
	if err != nil {
		return dto.DetectionResponse{}, err
	}

	if resp.Predict == "real" {
		detection.IsDeepFake = false
	} else {
		detection.IsDeepFake = true
	}

	if err := ds.dr.CreateDetection(ctx, &detection); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create detection history")
		return dto.DetectionResponse{}, errorz.ErrFailedToCreateDetection.WithTraceID(traceID)
	}

	return dto.DetectionResponse{
		Filename:   req.File.Filename,
		Prediction: resp.Predict,
	}, nil
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

func (ds *detectionService) createHttpRequest(traceID uuid.UUID, body *bytes.Buffer, writer *multipart.Writer) (dto.DetectionResponse, error) {
	request, err := http.NewRequest(http.MethodPost, env.GetEnv().DetectionUrl, body)
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

	return mapper.ToDetectionResponse(detections, req.CurrentPage, req.Limit), nil
}

func (ds *detectionService) BlockDetection(ctx context.Context, req dto.DeleteDetectionRequest) error {
	traceID := utils.GetTraceID(ctx)

	if err := ds.dr.DeleteDetection(ctx, req.DetectionID); err != nil {
		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DeleteDetection] failed to ban IP")
		return errorz.ErrBanIP.WithTraceID(traceID)
	}

	return nil
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

	return mapper.ToDetectionResponse(detections, req.CurrentPage, req.Limit), nil
}

// func (ds *detectionService) DetectDeepFake(ctx context.Context, req dto.DetectionRequest) (dto.DetectionResponse, error) {
// 	traceID := utils.GetTraceID(ctx)
// 	file, err := req.File.Open()
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to open image")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToOpenFile.WithTraceID(traceID)
// 	}
// 	defer file.Close()

// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to read image")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToReadFile.WithTraceID(traceID)
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	part, err := writer.CreateFormFile("file", req.File.Filename)
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create form file")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToCreateFormFile.WithTraceID(traceID)
// 	}

// 	_, err = part.Write(fileBytes)
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] to write bytes")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToWriteBytes.WithTraceID(traceID)
// 	}

// 	if err := writer.Close(); err != nil {
// 		return dto.DetectionResponse{}, fmt.Errorf("failed to close writer : %v", err)
// 	}

// 	request, err := http.NewRequest(http.MethodPost, env.GetEnv().DetectionUrl, body)
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create http request")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToCreateHTTPRequest.WithTraceID(traceID)
// 	}
// 	request.Header.Set("Content-Type", writer.FormDataContentType())

// 	client := &http.Client{}
// 	resp, err := client.Do(request)
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to make http client")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToInitiateHTTPClient.WithTraceID(traceID)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to analyze image")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToAnalyzeImage.WithTraceID(traceID)
// 	}

// 	bodyResp, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to read response body")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToReadFile.WithTraceID(traceID)
// 	}

// 	var detectResp dto.DetectionResponse
// 	if err := json.Unmarshal(bodyResp, &detectResp); err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to unmarshal json to struct")
// 		return dto.DetectionResponse{}, errorz.ErrUnmarshal.WithTraceID(traceID)
// 	}

// 	if detectResp.Prediction == "" {
// 		return dto.DetectionResponse{}, fmt.Errorf("detectResp struct empty. Failed to unmarshal : %v", fiber.StatusInternalServerError)
// 	}

// 	detection := entities.Detection{
// 		IPAddress:  req.IP,
// 		Email:      req.Email,
// 		CustomerID: req.CustomerID,
// 		Path:       req.Path,
// 		Method:     req.Method,
// 		StatusCode: uint16(200),
// 	}

// 	if detectResp.Prediction == "Real" {
// 		detection.IsDeepFake = false
// 	} else {
// 		detection.IsDeepFake = true
// 	}

// 	if err := ds.dr.CreateDetection(ctx, &detection); err != nil {
// 		ds.logger.WithFields(log.WithTraceID(traceID, err)).Error("[DetectionService][DetectDeepFake] failed to create detection history")
// 		return dto.DetectionResponse{}, errorz.ErrFailedToCreateDetection.WithTraceID(traceID)
// 	}

// 	return detectResp, nil
// }
