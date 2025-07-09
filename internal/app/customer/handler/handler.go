package handler

import (
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type customerHandler struct {
	cs  interfaces.ICustomerService
	val validator.Validator
}

func NewCustomerHandler(cs interfaces.ICustomerService, val validator.Validator) *customerHandler {
	return &customerHandler{
		cs:  cs,
		val: val,
	}
}

func (cs *customerHandler) SetEndpoint(router fiber.Router) {
	v1 := router.Group("/customers")
	v1.Post("/generate-key", middleware.Authenticate(), cs.GenerateAPIKey)
	v1.Patch("/update-key", middleware.Authenticate(), cs.UpdateAPIKey)
	v1.Get("/key-status", middleware.Authenticate(), cs.GetCustomerAPIKeyStatus)
}
