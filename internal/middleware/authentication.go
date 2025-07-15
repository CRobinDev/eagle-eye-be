package middleware

import (
	"strings"

	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if bearer == "" {
			log.Errorf("Invalid header, try again later")
			return errorz.ErrUnauthorized
		}

		tokenSlice := strings.Split(bearer, " ")
		if len(tokenSlice) != 2 {
			log.Errorf("Authorization header invalid")
			return errorz.ErrInvalidToken
		}

		jwtToken := tokenSlice[1]
		claims, err := jwt.DecodeToken(jwtToken)
		if err != nil {
			log.Errorf("decode jwt error : %v", err)
			return errorz.ErrFailedToDecodeJWT
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
