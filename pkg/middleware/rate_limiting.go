package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/keremdursn/hospital-case/pkg/errs"
	"github.com/keremdursn/hospital-case/pkg/metrics"
)

// AuthRateLimiter Auth endpointleri için rate limiter
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,               // 5 istek
		Expiration: 1 * time.Minute, // 1 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // IP bazlı limit
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("auth", c.IP()).Inc()
			return errs.SendErrorResponse(c, errs.ErrTooManyRequests)
		},
	})
}

// LoginRateLimiter Login endpointi için daha sıkı limit
func LoginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,               // 3 istek
		Expiration: 5 * time.Minute, // 5 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":login" // IP + endpoint bazlı
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("login", c.IP()).Inc()
			return errs.SendErrorResponse(c, errs.ErrTooManyRequests)
		},
	})
}

// GeneralRateLimiter Genel endpointler için rate limiter
func GeneralRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // 100 istek
		Expiration: 1 * time.Minute, // 1 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("general", c.IP()).Inc()
			return errs.SendErrorResponse(c, errs.ErrTooManyRequests)
		},
	})
}

// AdminRateLimiter Admin endpointleri için özel limit
func AdminRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,              // 20 istek
		Expiration: 1 * time.Minute, // 1 dakika
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + ":admin"
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Rate limit metriklerini artır
			metrics.RateLimitExceededCounter.WithLabelValues("admin", c.IP()).Inc()
			return errs.SendErrorResponse(c, errs.ErrTooManyRequests)
		},
	})
}
