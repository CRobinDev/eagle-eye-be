package handler

import (
	"context"
	"errors"
	"time"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/pkg/jwt"
	"github.com/CRobinDev/karsa/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (ph *paymentHandler) CreatePayment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if valErr := ph.val.Validate(req); valErr != nil {
		return valErr
	}

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.UserID = userID

	resp, err := ph.ps.CreatePayment(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "payments", resp)
}

func (ph *paymentHandler) UpdatePaymentStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.UpdatePaymentStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if valErr := ph.val.Validate(req); valErr != nil {
		return valErr
	}

	resp, err := ph.ps.UpdatePaymentStatus(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "payments", resp)
}

func (ph *paymentHandler) LatestPaymentStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.GetPaymentStatusRequest
	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.UserID = userID

	resp, err := ph.ps.LatestPaymentStatus(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "payments", resp)

}
