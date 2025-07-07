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

func (ah *authHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if valErr := ah.val.Validate(req); valErr != nil {
		return valErr
	}

	if err := ah.as.Register(ctx, req); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}

	return response.Created(c)
}

func (ah *authHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if valErr := ah.val.Validate(req); valErr != nil {
		return valErr
	}

	resp, err := ah.as.Login(ctx, req)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")

	default:
	}
	return response.Success(c, "users", resp)
}

func (ah *authHandler) TokenAuthentication(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	resp, err := ah.as.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	
	return response.Success(c, "users", resp)
}
