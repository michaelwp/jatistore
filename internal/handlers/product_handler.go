package handlers

import (
	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product with the provided data
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Product data"
// @Success 201 {object} models.APIResponse{data=models.Product}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req models.CreateProductRequest

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
			Error:   "Product name is required",
		})
	}

	if req.SKU == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product SKU is required",
		})
	}

	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product price cannot be negative",
		})
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProductByID retrieves a product by its ID
// @Summary Get product by ID
// @Description Get a product by its unique identifier
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.APIResponse{data=models.Product}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product ID is required",
		})
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    product,
	})
}

// GetAllProducts retrieves all products
// @Summary Get all products
// @Description Get a list of all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Product}
// @Failure 500 {object} models.APIResponse
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    products,
	})
}

// UpdateProduct updates an existing product
// @Summary Update a product
// @Description Update a product with the provided data
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.UpdateProductRequest true "Updated product data"
// @Success 200 {object} models.APIResponse{data=models.Product}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product ID is required",
		})
	}

	var req models.UpdateProductRequest

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
			Error:   "Product name is required",
		})
	}

	if req.SKU == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product SKU is required",
		})
	}

	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product price cannot be negative",
		})
	}

	product, err := h.productService.UpdateProduct(id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct deletes a product
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product ID is required",
		})
	}

	err := h.productService.DeleteProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}
