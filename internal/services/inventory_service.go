package services

import (
	"fmt"

	"jatistore/internal/models"
	"jatistore/internal/repository"

	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

func NewInventoryService(inventoryRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) CreateInventory(req *models.CreateInventoryRequest) (*models.Inventory, error) {
	inventory := &models.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Location:  req.Location,
	}

	if err := s.inventoryRepo.Create(inventory); err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	// Get the created inventory with product information
	createdInventory, err := s.inventoryRepo.GetByID(inventory.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created inventory: %w", err)
	}

	return createdInventory, nil
}

func (s *InventoryService) GetInventoryByID(id string) (*models.Inventory, error) {
	inventoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid inventory ID: %w", err)
	}

	inventory, err := s.inventoryRepo.GetByID(inventoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return inventory, nil
}

func (s *InventoryService) GetAllInventory() ([]*models.Inventory, error) {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get inventories: %w", err)
	}

	return inventories, nil
}

func (s *InventoryService) UpdateInventory(id string, req *models.UpdateInventoryRequest) (*models.Inventory, error) {
	inventoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid inventory ID: %w", err)
	}

	// Get existing inventory
	existingInventory, err := s.inventoryRepo.GetByID(inventoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing inventory: %w", err)
	}

	// Update inventory fields
	existingInventory.Quantity = req.Quantity
	existingInventory.Location = req.Location

	if err := s.inventoryRepo.Update(existingInventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Get the updated inventory with product information
	updatedInventory, err := s.inventoryRepo.GetByID(inventoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated inventory: %w", err)
	}

	return updatedInventory, nil
}

func (s *InventoryService) DeleteInventory(id string) error {
	inventoryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid inventory ID: %w", err)
	}

	if err := s.inventoryRepo.Delete(inventoryID); err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	return nil
}

func (s *InventoryService) AdjustStock(req *models.AdjustStockRequest) (*models.InventoryTransaction, error) {
	// Use product ID as a string (no UUID parsing)
	productID := req.ProductID

	// Get current inventory for the product
	inventories, err := s.inventoryRepo.GetByProductIDString(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product inventory: %w", err)
	}

	if len(inventories) == 0 {
		return nil, fmt.Errorf("no inventory found for product")
	}

	// For simplicity, we'll adjust the first inventory record
	// In a real application, you might want to specify which location to adjust
	inventory := inventories[0]

	// Calculate new quantity based on transaction type
	var newQuantity int
	switch req.Type {
	case "in":
		newQuantity = inventory.Quantity + req.Quantity
	case "out":
		newQuantity = inventory.Quantity - req.Quantity
		if newQuantity < 0 {
			return nil, fmt.Errorf("insufficient stock: current quantity is %d, trying to remove %d", inventory.Quantity, req.Quantity)
		}
	case "adjustment":
		newQuantity = req.Quantity
		if newQuantity < 0 {
			return nil, fmt.Errorf("quantity cannot be negative")
		}
	default:
		return nil, fmt.Errorf("invalid transaction type: %s", req.Type)
	}

	// Update inventory quantity
	inventory.Quantity = newQuantity
	if err := s.inventoryRepo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		ProductID: productID, // string
		Type:      req.Type,
		Quantity:  req.Quantity,
		Reason:    req.Reason,
		Reference: req.Reference,
	}

	if err := s.inventoryRepo.CreateTransactionString(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Get the created transaction with product information
	createdTransaction, err := s.inventoryRepo.GetTransactionsByProductIDString(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created transaction: %w", err)
	}

	if len(createdTransaction) > 0 {
		return createdTransaction[0], nil
	}

	return transaction, nil
}

func (s *InventoryService) GetInventoryByProductID(productID string) ([]*models.Inventory, error) {
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	inventories, err := s.inventoryRepo.GetByProductID(parsedProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product inventory: %w", err)
	}

	return inventories, nil
}

func (s *InventoryService) GetTransactionsByProductID(productID string) ([]*models.InventoryTransaction, error) {
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	transactions, err := s.inventoryRepo.GetTransactionsByProductID(parsedProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product transactions: %w", err)
	}

	return transactions, nil
}
