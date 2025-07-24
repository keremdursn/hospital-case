package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/pkg/errs"
)

// GlobalErrorHandler tüm hataları merkezi olarak yakalar ve standart response döner
func GlobalErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		// Eğer AppError ise, ona göre response dön
		if appErr, ok := err.(*errs.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    appErr.Code,
					"message": appErr.Message,
				},
			})
		}

		// Diğer hatalar için generic 500 response
		return c.Status(500).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "ERR_INTERNAL",
				"message": "Internal server error",
			},
		})
	}
}
