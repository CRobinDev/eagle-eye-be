package jwt

import (
	"time"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	secretKey   []byte
	expiredTime time.Duration
)

func InitJWT() {
	secretKey = env.GetEnv().JwtSecretKey
	expiredTime = env.GetEnv().JwtExpiredTime
}

type Claims struct {
	Role        string    `json:"role"`
	UserID      uuid.UUID `json:"user_id"`
	DisplayName string    `json:"username"`
	jwt.RegisteredClaims
}

func CreateToken(user *entities.User) (string, error) {
	claims := Claims{
		Role:        string(user.Role),
		UserID:      user.ID,
		DisplayName: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredTime)),
		},
	}
	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return unsignedJWT.SignedString(secretKey)
}

func DecodeToken(tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Error(map[string]interface{}{
			"error": err.Error(),
		}, "[jwt.DecodeToken] failed to parse token")
		return claims, jwt.ErrTokenExpired
	}

	if !token.Valid {
		return claims, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func GetUser(c *fiber.Ctx) (uuid.UUID, error) {
	claims, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Can't retrieve claims")
	}

	userID := claims

	return userID, nil
}
