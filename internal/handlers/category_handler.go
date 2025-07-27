package handlers

import (
	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Create a new category with the provided data
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param category body models.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.APIResponse{data=models.Category}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req models.CreateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Category name is required",
		})
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    category,
	})
}

// GetCategoryByID retrieves a category by its ID
// @Summary Get category by ID
// @Description Get a category by its unique identifier
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Category ID"
// @Success 200 {object} models.APIResponse{data=models.Category}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Category ID is required",
		})
	}

	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    category,
	})
}

// GetAllCategories retrieves all categories
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Category}
// @Failure 500 {object} models.APIResponse
// @Router /categories [get]
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"  // <-- required header
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    categories,
	})
}

// UpdateCategory updates an existing category
// @Summary Update a category
// @Description Update a category with the provided data
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Category ID"
// @Param category body models.UpdateCategoryRequest true "Updated category data"
// @Success 200 {object} models.APIResponse{data=models.Category}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Category ID is required",
		})
	}

	var req models.UpdateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Category name is required",
		})
	}

	category, err := h.categoryService.UpdateCategory(id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:    category,
	})
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Category ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Category ID is required",
		})
	}

	err := h.categoryService.DeleteCategory(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}
