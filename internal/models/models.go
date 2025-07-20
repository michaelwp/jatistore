package models

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the inventory system
type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	SKU         string    `json:"sku" db:"sku"`
	CategoryID  uuid.UUID `json:"category_id" db:"category_id"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Category    *Category `json:"category,omitempty"`
}

// Category represents a product category
type Category struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Inventory represents inventory stock for a product
type Inventory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Location  string    `json:"location" db:"location"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Product   *Product  `json:"product,omitempty"`
}

// InventoryTransaction represents inventory movement
type InventoryTransaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID string    `json:"product_id" db:"product_id"` // changed from uuid.UUID to string
	Type      string    `json:"type" db:"type"`             // "in", "out", "adjustment"
	Quantity  int       `json:"quantity" db:"quantity"`
	Reason    string    `json:"reason" db:"reason"`
	Reference string    `json:"reference" db:"reference"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Product   *Product  `json:"product,omitempty"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	SKU         string  `json:"sku" validate:"required"`
	CategoryID  string  `json:"category_id" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	SKU         string  `json:"sku" validate:"required"`
	CategoryID  string  `json:"category_id" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// CreateInventoryRequest represents the request to create inventory
type CreateInventoryRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=0"`
	Location  string `json:"location" validate:"required"`
}

// UpdateInventoryRequest represents the request to update inventory
type UpdateInventoryRequest struct {
	Quantity int    `json:"quantity" validate:"required,min=0"`
	Location string `json:"location" validate:"required"`
}

// AdjustStockRequest represents the request to adjust stock
type AdjustStockRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=in out adjustment"`
	Reason    string `json:"reason" validate:"required"`
	Reference string `json:"reference"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
