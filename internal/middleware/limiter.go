package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type rateLimiter struct {
	visitors sync.Map // map[string]*visitor
	rate     time.Duration
	burst    int
}

type visitor struct {
	tokens     int
	lastRefill time.Time
	mutex      sync.Mutex
}

func newRateLimiter(rate time.Duration, burst int) *rateLimiter {
	rl := &rateLimiter{
		rate:  rate,
		burst: burst,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *rateLimiter) getVisitor(key string) *visitor {
	v, ok := rl.visitors.Load(key)
	if ok {
		return v.(*visitor)
	}
	visitor := &visitor{
		tokens:     rl.burst,
		lastRefill: time.Now(),
	}
	rl.visitors.Store(key, visitor)
	return visitor
}

func (rl *rateLimiter) allow(key string) (allowed bool, remaining int, reset time.Time) {
	v := rl.getVisitor(key)
	v.mutex.Lock()
	defer v.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(v.lastRefill)
	refill := int(elapsed / rl.rate)
	if refill > 0 {
		v.tokens += refill
		if v.tokens > rl.burst {
			v.tokens = rl.burst
		}
		v.lastRefill = v.lastRefill.Add(time.Duration(refill) * rl.rate)
	}
	allowed = v.tokens > 0
	if allowed {
		v.tokens--
	}
	remaining = v.tokens
	reset = v.lastRefill.Add(rl.rate * time.Duration(rl.burst-v.tokens))
	return
}

func (rl *rateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		rl.visitors.Range(func(key, value any) bool {
			v := value.(*visitor)
			v.mutex.Lock()
			idle := now.Sub(v.lastRefill)
			v.mutex.Unlock()
			if idle > time.Hour {
				rl.visitors.Delete(key)
			}
			return true
		})
	}
}

// Optionally, you can use a custom key extractor (e.g., by user ID, API key, etc.)
func extractKey(c *fiber.Ctx) string {
	// Default: use IP
	return c.IP()
}

func RateLimiter(rate time.Duration, burst int, logger *logrus.Logger) fiber.Handler {
	limiter := newRateLimiter(rate, burst)
	return func(c *fiber.Ctx) error {
		key := extractKey(c)
		allowed, remaining, reset := limiter.allow(key)
		c.Set("X-RateLimit-Limit", strconv.Itoa(burst))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
		if !allowed {
			//taruh logger
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests",
			})
		}
		return c.Next()
	}
}
