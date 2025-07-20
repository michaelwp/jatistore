package middleware

import (
	"strings"

	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	userService *services.UserService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService}
}

// Authenticate validates JWT token and sets user context
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error:   "Authorization header is required",
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error:   "Invalid authorization header format",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := m.userService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error:   "Invalid or expired token",
			})
		}

		// Get user from database to ensure user still exists and is active
		user, err := m.userService.GetUserByID(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error:   "User not found",
			})
		}

		if !user.IsActive {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error:   "Account is deactivated",
			})
		}

		// Set user in context
		c.Locals("user", user)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequireRole creates middleware that requires specific role(s)
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by Authenticate middleware)
		user := c.Locals("user").(*models.User)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error:   "Authentication required",
			})
		}

		// Check if user has required role
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Error:   "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// GetCurrentUser retrieves the current user from context
func GetCurrentUser(c *fiber.Ctx) *models.User {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	return user.(*models.User)
}

// GetCurrentUserID retrieves the current user ID from context
func GetCurrentUserID(c *fiber.Ctx) uuid.UUID {
	user := GetCurrentUser(c)
	if user == nil {
		return uuid.Nil
	}
	return user.ID
}

// GetCurrentUserClaims retrieves the current user claims from context
func GetCurrentUserClaims(c *fiber.Ctx) *models.Claims {
	claims := c.Locals("claims")
	if claims == nil {
		return nil
	}
	return claims.(*models.Claims)
}
