package handler

import (
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	as  interfaces.IAuthService
	val validator.Validator
}

func NewAuthHandler(as interfaces.IAuthService, val validator.Validator) *authHandler {
	return &authHandler{
		as:  as,
		val: val,
	}
}

func (ah *authHandler) SetEndpoint(router fiber.Router) {
	v1 := router.Group("/auth")
	v1.Post("/login", ah.Login)
	v1.Post("/register", ah.Register)
	v1.Get("/me", middleware.Authenticate(), ah.TokenAuthentication)
}
