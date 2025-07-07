package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func Recover() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("[RECOVER] panic: %v\n%s", p, debug.Stack())
				err = fiber.NewError(fiber.StatusInternalServerError, "Panic causes Internal Server Error")
			}
		}()
		return ctx.Next()
	}
}
