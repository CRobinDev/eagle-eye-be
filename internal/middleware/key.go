package middleware

import (
	"os"
	"strings"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ValidateKey(cr interfaces.ICustomerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get(env.GetEnv().HttpHeader)
		if apiKey == "" {
			return errorz.ErrMissingAPIKey
		}

		splitter := strings.Split(apiKey, "_")
		if len(splitter) != 2 {
			return errorz.ErrInvalidAPIKey
		}

		prefix := splitter[0]
		key := splitter[1]

		if len(key) != 64 {
			return errorz.ErrUnauthorized
		}

		hashedApiKey := utils.HashAPIKey(prefix, key)

		customer, err := cr.UpdateUsage(c.UserContext(), prefix)
		if err != nil {
			return errorz.ErrFailedToUpdateUsage
		}

		if customer.ApiKey != hashedApiKey || customer.Prefix != prefix {
			return errorz.ErrMismatchAPIKey
		}

		var cust entities.Customer
		if  cust == customer{
			logrus.Error("customer is empty.")
			os.Exit(1)
		}
		c.Locals("customer_id", customer.ID)
		return c.Next()
	}
}
