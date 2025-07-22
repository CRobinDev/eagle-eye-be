package middleware

import (
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/gofiber/fiber/v2"
)

func ValidateIP(ds interfaces.IDetectionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		isBanned, err := ds.GetDetectionByIP(c.UserContext(), dto.ValidateIPRequest{
			IP: ip,
		})

		if err != nil {
			return errorz.ErrFailedToGetDetection
		}

		if isBanned {
			return errorz.ErrBanned
		}

		return c.Next()
	}
}
