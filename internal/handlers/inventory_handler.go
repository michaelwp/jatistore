package handlers

import (
	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// CreateInventory creates a new inventory record
// @Summary Create a new inventory record
// @Description Create a new inventory record with the provided data
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param inventory body models.CreateInventoryRequest true "Inventory data"
// @Success 201 {object} models.APIResponse{data=models.Inventory}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /inventory [post]
func (h *InventoryHandler) CreateInventory(c *fiber.Ctx) error {
	var req models.CreateInventoryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.ProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product ID is required",
		})
	}

	if req.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Quantity cannot be negative",
		})
	}

	if req.Location == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Location is required",
		})
	}

	inventory, err := h.inventoryService.CreateInventory(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Inventory created successfully",
		Data:    inventory,
	})
}

// GetInventoryByID retrieves an inventory record by its ID
// @Summary Get inventory by ID
// @Description Get an inventory record by its unique identifier
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Inventory ID"
// @Success 200 {object} models.APIResponse{data=models.Inventory}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /inventory/{id} [get]
func (h *InventoryHandler) GetInventoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Inventory ID is required",
		})
	}

	inventory, err := h.inventoryService.GetInventoryByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    inventory,
	})
}

// GetAllInventory retrieves all inventory records
// @Summary Get all inventory records
// @Description Get a list of all inventory records
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.APIResponse{data=[]models.Inventory}
// @Failure 500 {object} models.APIResponse
// @Router /inventory [get]
func (h *InventoryHandler) GetAllInventory(c *fiber.Ctx) error {
	inventories, err := h.inventoryService.GetAllInventory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    inventories,
	})
}

// UpdateInventory updates an existing inventory record
// @Summary Update an inventory record
// @Description Update an inventory record with the provided data
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Inventory ID"
// @Param inventory body models.UpdateInventoryRequest true "Updated inventory data"
// @Success 200 {object} models.APIResponse{data=models.Inventory}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /inventory/{id} [put]
func (h *InventoryHandler) UpdateInventory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Inventory ID is required",
		})
	}

	var req models.UpdateInventoryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Quantity cannot be negative",
		})
	}

	if req.Location == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Location is required",
		})
	}

	inventory, err := h.inventoryService.UpdateInventory(id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Inventory updated successfully",
		Data:    inventory,
	})
}

// DeleteInventory deletes an inventory record
// @Summary Delete an inventory record
// @Description Delete an inventory record by its ID
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Inventory ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /inventory/{id} [delete]
func (h *InventoryHandler) DeleteInventory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Inventory ID is required",
		})
	}

	err := h.inventoryService.DeleteInventory(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Inventory deleted successfully",
	})
}

// AdjustStock adjusts inventory stock levels
// @Summary Adjust inventory stock
// @Description Adjust inventory stock levels (in/out/adjustment)
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param adjustment body models.AdjustStockRequest true "Stock adjustment data"
// @Success 200 {object} models.APIResponse{data=models.InventoryTransaction}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /inventory/adjust [post]
func (h *InventoryHandler) AdjustStock(c *fiber.Ctx) error {
	var req models.AdjustStockRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.ProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Product ID is required",
		})
	}

	if req.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Transaction type is required",
		})
	}

	if req.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Reason is required",
		})
	}

	transaction, err := h.inventoryService.AdjustStock(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Stock adjusted successfully",
		Data:    transaction,
	})
}
