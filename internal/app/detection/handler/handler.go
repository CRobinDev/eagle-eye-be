package handler

import (
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/middleware"
	"github.com/CRobinDev/karsa/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type detectionHandler struct {
	ds  interfaces.IDetectionService
	cs  interfaces.ICustomerService
	val validator.Validator
}

func NewDetectionHandler(ds interfaces.IDetectionService, cs interfaces.ICustomerService, val validator.Validator) *detectionHandler {
	return &detectionHandler{
		ds:  ds,
		cs:  cs,
		val: val,
	}
}

func (dh *detectionHandler) SetEndpoint(router fiber.Router) {
	v1 := router.Group("/detections")
	v1.Get("/get-detections", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.GetDetection)
	v1.Get("/get-detected", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.GetDetectedDeepFake)
	v1.Get("/get-details/:id", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.GetDetectionByID)
	v1.Get("/get-customer-usage", middleware.Authenticate(), middleware.RequiredOneOfRoles(), dh.CustomerUsage)
	v1.Post("/detect-image", middleware.ValidateIP(dh.ds), middleware.ValidateKey(dh.cs), dh.DetectDeepFakeImage)
	v1.Post("/detect-audio", middleware.ValidateIP(dh.ds), middleware.ValidateKey(dh.cs), dh.DetectDeepFakeAudio)
	v1.Patch("/unban-ip", middleware.Authenticate(), dh.UnbanIP)
	v1.Delete("/ban-ip/:ip", middleware.Authenticate(), dh.BanIP)
}
