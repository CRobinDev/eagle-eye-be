package handler

import (
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	us  interfaces.IUserService
	val validator.Validator
}

func NewUserHandler(us interfaces.IUserService, val validator.Validator) *userHandler {
	return &userHandler{
		us:  us,
		val: val,
	}
}

func (uh *userHandler) SetEndpoint(router fiber.Router) {
	v1 := router.Group("/users")
	v1.Patch("/update-info", middleware.Authenticate(), uh.UpdateUser)
	v1.Patch("/update-password", middleware.Authenticate(), uh.UpdatePassword)
	v1.Delete("/delete-account", middleware.Authenticate(), uh.DeleteUser)
}
