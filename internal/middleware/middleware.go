package middleware

import (
	"log"

	"jatistore/internal/models"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles application errors
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Log the error
	log.Printf("Error: %v", err)

	// Default error response
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Return JSON error response
	return c.Status(code).JSON(models.APIResponse{
		Success: false,
		Error:   message,
	})
}
