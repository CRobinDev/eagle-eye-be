package middleware

import (
	"context"

	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func Trace() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := uuid.New()
		ctx := context.WithValue(c.UserContext(), utils.TraceCtxKey, traceID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
