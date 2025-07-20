package handlers

import (
	"jatistore/internal/middleware"
	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username must be between 3 and 50 characters",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Email is required",
		})
	}

	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Password is required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Password must be at least 6 characters",
		})
	}

	if req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Role is required",
		})
	}

	if req.Role != "admin" && req.Role != "user" && req.Role != "cashier" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Role must be admin, user, or cashier",
		})
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
	}

	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Password is required",
		})
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// GetProfile retrieves the current user's profile
// @Summary Get user profile
// @Description Get the current authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 401 {object} models.APIResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update the current authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
	}

	var req models.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username must be between 3 and 50 characters",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Email is required",
		})
	}

	if req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Role is required",
		})
	}

	if req.Role != "admin" && req.Role != "user" && req.Role != "cashier" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Role must be admin, user, or cashier",
		})
	}

	user, err := h.userService.UpdateUser(currentUser.ID, &req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    user,
	})
}

// ChangePassword changes the current user's password
// @Summary Change password
// @Description Change the current authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Password change data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
	}

	var req models.ChangePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.CurrentPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Current password is required",
		})
	}

	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "New password is required",
		})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "New password must be at least 6 characters",
		})
	}

	err := h.userService.ChangePassword(currentUser.ID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// GetAllUsers retrieves all users (admin only)
// @Summary Get all users
// @Description Get all users in the system (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.User}
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Router /auth/users [get]
func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve users",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    users,
	})
}

// GetUserByID retrieves a user by ID (admin only)
// @Summary Get user by ID
// @Description Get a specific user by ID (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /auth/users/{id} [get]
func (h *AuthHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateUser updates a user (admin only)
// @Summary Update user
// @Description Update a specific user (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /auth/users/{id} [put]
func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Username must be between 3 and 50 characters",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Email is required",
		})
	}

	if req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Role is required",
		})
	}

	if req.Role != "admin" && req.Role != "user" && req.Role != "cashier" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Role must be admin, user, or cashier",
		})
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "user not found" {
			status = fiber.StatusNotFound
		} else if err.Error() == "username already exists" || err.Error() == "email already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser deletes a user (admin only)
// @Summary Delete user
// @Description Delete a specific user (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /auth/users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "user not found" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
