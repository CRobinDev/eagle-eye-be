package handler

import (
	"context"
	"errors"
	"time"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/jwt"
	"github.com/CRobinDev/karsa/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (dh *detectionHandler) DetectDeepFakeImage(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	email := c.FormValue("email")

	var req dto.DetectionRequest
	req.File = file
	req.Email = email
	req.IP = c.IP()
	req.Path = c.Path()
	req.Method = c.Method()
	req.CustomerID = c.Locals("customer_id").(uuid.UUID)

	resp, err := dh.ds.DetectDeepFakeImage(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.Success(c, "detections", resp)
}

func (dh *detectionHandler) GetDetection(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	currentPage := c.QueryInt("current")

	limit := c.QueryInt("limit")

	var req dto.GetDetectionRequest
	req.CurrentPage = uint64(currentPage)
	req.Limit = uint64(limit)

	if valErr := dh.val.Validate(req); valErr != nil {
		return valErr
	}

	if c.Locals("role").(string) == "user" && c.Locals("is_customer").(bool) {
		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}
		req.UserID = userID
	}

	resp, err := dh.ds.GetDetectionData(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.Success(c, "detections", resp)
}

func (dh *detectionHandler) DetectDeepFakeAudio(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	email := c.FormValue("email") // abu: BINGUNG

	var req dto.DetectionRequest
	req.File = file
	req.Email = email
	req.IP = c.IP()
	req.Path = c.Path()
	req.Method = c.Method()
	req.CustomerID = c.Locals("customer_id").(uuid.UUID)

	resp, err := dh.ds.DetectDeepFakeAudio(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.Success(c, "detections", resp)
}

func (dh *detectionHandler) GetDetectedDeepFake(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	currentPage := c.QueryInt("current")

	limit := c.QueryInt("limit")

	var req dto.GetDetectionRequest
	req.CurrentPage = uint64(currentPage)
	req.Limit = uint64(limit)

	if valErr := dh.val.Validate(req); valErr != nil {
		return valErr
	}

	if c.Locals("role").(string) == "user" && c.Locals("is_customer").(bool) {
		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}
		req.UserID = userID
	}

	resp, err := dh.ds.GetDeepFakeDetected(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.Success(c, "detections", resp)
}

func (dh *detectionHandler) BanIP(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var req dto.DeleteDetectionRequest
	req.DetectionID = uint16(id)

	if c.Locals("role").(string) != "user" && !c.Locals("is_customer").(bool) {
		return errorz.ErrForbiddenRole
	}

	if valErr := dh.val.Validate(req); valErr != nil {
		return valErr
	}

	if err := dh.ds.BlockDetection(ctx, req); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.NoContent(c)
}
