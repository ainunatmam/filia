package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type RateLimiterConfig struct {
	Max            int
	Expiration     time.Duration
	KeyGenerator   func(c *fiber.Ctx) string
	LimiterHandler func(c *fiber.Ctx) error
}

var rateLimiterConfig RateLimiterConfig

// InitRateLimiter initializes the rate limiter configuration with provided values
func InitRateLimiter(maxRequests int, expirationSeconds int) {
	rateLimiterConfig = RateLimiterConfig{
		Max:          maxRequests,
		Expiration:   time.Duration(expirationSeconds) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
	}
}

// RateLimiter returns the rate limiting middleware
func RateLimiter() fiber.Handler {
	// Initialize config if not already done
	if rateLimiterConfig.Max == 0 {
		InitRateLimiter(100, 60) // fallback to defaults
	}

	return limiter.New(limiter.Config{
		Max:        rateLimiterConfig.Max,
		Expiration: rateLimiterConfig.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP + user agent for more accurate rate limiting
			return c.IP() + ":" + c.Get("User-Agent")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "rate limit exceeded",
				"retry_after": rateLimiterConfig.Expiration.Seconds(),
			})
		},
		Storage: nil, // Use default memory storage
	})
}

// IPBanned checks if an IP is banned (stub - can be extended with Redis)
func IPBanned(ip string) bool {
	// TODO: Implement with Redis or database for persistence
	return false
}

// BanIP bans an IP (stub - can be extended with Redis)
func BanIP(ip string, duration time.Duration) {
	// TODO: Implement with Redis or database for persistence
}
