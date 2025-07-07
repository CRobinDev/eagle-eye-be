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

func (uh *userHandler) UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5 * time.Second)
	defer cancel()

	var req dto.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if valErr := uh.val.Validate(req); valErr != nil {
		return valErr
	}

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.UserID = userID

	if err := uh.us.UpdateUser(ctx, req); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.Success(c, "users", nil)
}

func (uh *userHandler) UpdatePassword(c *fiber.Ctx) error {
	var req dto.UpdatePasswordRequest
	ctx, cancel := context.WithTimeout(c.UserContext(), 5 * time.Second)
	defer cancel()

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if valErr := uh.val.Validate(req); valErr != nil {
		return valErr
	}

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	req.UserID = userID
	if err := uh.us.UpdatePassword(ctx, req); err != nil {
		return err
	}

	select {
	case <- ctx.Done():
		return errors.New("timeout")
	default:
	}

	return response.Success(c, "users", nil)
}

func (uh *userHandler) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5 * time.Second)
	defer cancel()

	userID, err := jwt.GetUser(c)
	if err != nil {
		return err
	}

	if err := uh.us.DeleteUser(ctx, userID); err != nil {
		return err
	}

	select {
	case <- ctx.Done():
		return errors.New("timeout")
	default:
	}
	
	return response.NoContent(c)
}