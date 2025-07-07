package handler

import (
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type paymentHandler struct {
	ps  interfaces.IPaymentService
	val validator.Validator
}

func NewPaymentHandler(ps interfaces.IPaymentService, val validator.Validator) *paymentHandler {
	return &paymentHandler{
		ps:  ps,
		val: val,
	}
}

func (ph *paymentHandler) SetEndpoint(router fiber.Router) {
	v1 := router.Group("payments")
	v1.Post("/create-payment", middleware.Authenticate(), ph.CreatePayment)
	v1.Post("/update-status", ph.UpdatePaymentStatus)
	v1.Get("/latest-status", middleware.Authenticate(), ph.LatestPaymentStatus)
}
