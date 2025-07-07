package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()
		latency := stop.Sub(start)

		fields := logrus.Fields{
			"[ip]":      c.IP(),
			"[host]":    c.Hostname(),
			"[url]":     c.OriginalURL(),
			"[method]":  c.Method(),
			"[status]":  c.Response().StatusCode(),
			"[latency]": latency.String(),
			"[ua]":      c.Get("User-Agent"),
		}

		status := c.Response().StatusCode()
		switch {
		case status >= 500:
			logger.WithFields(fields).Error("[server error]")
		case status >= 400:
			logger.WithFields(fields).Warn("[client error]")
		default:
			logger.WithFields(fields).Info("[success]")
		}

		return err
	}
}
