package router

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"runtime/debug"
	"time"

	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/errorz"
	ut "github.com/CRobinDev/karsa/pkg/utils"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/sirupsen/logrus"
)

func NewFiber(logger *logrus.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:          "Eagle Eye Backend",
		JSONEncoder:      json.Marshal,
		JSONDecoder:      json.Unmarshal,
		ErrorHandler:     newErrorHandler(),
		BodyLimit:        10 * 1024 * 1024,
		DisableKeepalive: true,
		StrictRouting:    true,
		CaseSensitive:    true,
		UnescapePath:     true,
	})

	app.Use(middleware.Cors())
	app.Use(middleware.Trace())
	app.Use(healthcheck.New())
	app.Use(middleware.RateLimiter(time.Second, 1000, logger))
	app.Use(middleware.Logger(logger))
	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: stackTraceLogger,
	}))
	app.Use(compress.New(compress.Config{
		Level: gzip.BestSpeed,
	}))

	return app
}

func newErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {

		var apiError *errorz.Errors
		if errors.As(err, &apiError) {
			response := fiber.Map{
				"message": apiError.Error(),
			}

			if apiError.TraceID != nil {
				response["trace_id"] = apiError.TraceID.String()
			}

			return c.Status(apiError.Code).JSON(response)
		}

		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			return c.Status(fiberError.Code).JSON(fiber.Map{
				"message": utils.StatusMessage(fiberError.Code),
				"error":   err,
			})
		}

		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			validationDetails := fiber.Map{}
			for field, msg := range validationError {
				validationDetails[field] = msg
			}

			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": validationDetails,
			})
		}

		log.Errorf("Unhandled Error : %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": utils.StatusMessage(fiber.StatusInternalServerError),
			"error":   err.Error(),
		})
	}
}

func stackTraceLogger(c *fiber.Ctx, p interface{}) {
	traceID := ut.GetTraceID(c.UserContext())
	log.Tracef("[RECOVER] traceID=%v, panic=%v\n%s", traceID, p, debug.Stack())
}
