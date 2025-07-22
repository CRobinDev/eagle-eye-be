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

func (ch *customerHandler) GenerateAPIKey(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.GenerateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.UserID = userID
	resp, err := ch.cs.GenerateAPIKey(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "customers", resp)
}

func (ch *customerHandler) UpdateAPIKey(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.GenerateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.UserID = userID
	resp, err := ch.cs.UpdateAPIKey(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "customers", resp)

}

func (ch *customerHandler) GetCustomerAPIKeyStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.GetCustomerRequest

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.CustomerID = userID
	resp, err := ch.cs.GetCustomerAPIStatus(ctx, req)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "customers", resp)

}

func (ch *customerHandler) GetCustomerTotalAPICalls(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	var req dto.GetCustomerRequest

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.CustomerID = userID
	resp, err := ch.cs.GetCustomerTotalCalls(ctx, req)
	if err != nil {
		return err
	}
	
	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "customers", resp)

}

func (ch *customerHandler) GetCustomerTier(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5 * time.Second)
	defer cancel()
	
	var req dto.GetCustomerRequest

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.CustomerID = userID

	resp, err := ch.cs.GetCustomerTier(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Success(c, "customers", resp)
}