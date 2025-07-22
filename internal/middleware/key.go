package middleware

import (
	"errors"
	"strings"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

func ValidateKey(cs interfaces.ICustomerService) fiber.Handler {
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

		customer, err := cs.UpdateUsage(c.UserContext(), dto.UpdateUsageRequest{
			Prefix: prefix,
		})

		if err != nil {
			var pqErr *pgconn.PgError
			if errors.As(err, &pqErr) {
				if pqErr.ConstraintName == "customers_check" {
					return errorz.ErrMonthlyLimitReached
				}
			}
			return errorz.ErrFailedToUpdateUsage
		}

		if customer.HashedKey != hashedApiKey || customer.Prefix != prefix {
			return errorz.ErrMismatchAPIKey
		}

		var cust dto.UpdateUsageResponse
		if cust == customer {
			logrus.Error("customer is empty.")
			return errorz.ErrUserNotFound
		}

		c.Locals("customer_id", customer.ID)
		return c.Next()
	}
}
