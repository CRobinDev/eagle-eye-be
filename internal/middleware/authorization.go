package middleware

import (
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/gofiber/fiber/v2"
)

func RequiredOneOfRoles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Can't retrieve claims")
		}
		if role != "admin" && role != "user" {
			return errorz.ErrForbiddenRole
		}
		return c.Next()
	}
}
