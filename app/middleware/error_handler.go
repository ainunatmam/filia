package middleware

import (
	"mini-exchange/app/errors"
	"mini-exchange/app/libraries"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a centralized error handling middleware
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			// Log the error
			libraries.Logger.Error().
				Err(err).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Str("ip", c.IP()).
				Msg("Request error")

			// Check if it's our custom AppError
			if appErr, ok := err.(*errors.AppError); ok {
				return c.Status(appErr.Status).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    appErr.Code,
						"message": appErr.Message,
					},
				})
			}

			// Handle fiber errors
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    "FIBER_ERROR",
						"message": fiberErr.Message,
					},
				})
			}

			// Default internal server error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    errors.ErrCodeInternal,
					"message": "Internal server error",
				},
			})
		}

		return nil
	}
}
