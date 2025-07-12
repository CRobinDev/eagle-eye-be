package handler

import (
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type detectionHandler struct {
	ds  interfaces.IDetectionService
	cr  interfaces.ICustomerRepository
	val validator.Validator
}

func NewDetectionHandler(ds interfaces.IDetectionService, cr interfaces.ICustomerRepository, val validator.Validator) *detectionHandler {
	return &detectionHandler{
		ds:  ds,
		cr:  cr,
		val: val,
	}
}

func (dh *detectionHandler) SetEndpoint(router fiber.Router) {
	v1 := router.Group("/detections")
	v1.Get("/get-detections", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.GetDetection)
	v1.Get("/get-detected", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.GetDetectedDeepFake)
	v1.Get("/get-details/:id", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.GetDetectionByID)
	v1.Get("/get-customer-usage", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.CustomerUsage)
	v1.Post("/detect-image", middleware.ValidateKey(dh.cr), dh.DetectDeepFakeImage)
	v1.Post("/detect-audio", middleware.ValidateKey(dh.cr), dh.DetectDeepFakeAudio)
	v1.Patch("/unban-ip/:id", middleware.Authenticate(), dh.UnbanIP)
	v1.Delete("/ban-ip/:id", middleware.Authenticate(), dh.BanIP)
}
