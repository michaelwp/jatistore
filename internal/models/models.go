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

// Customer represents a customer in the POS system
type Customer struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Order represents a sales order in the POS system
type Order struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	OrderNumber    string      `json:"order_number" db:"order_number"`
	CustomerID     *uuid.UUID  `json:"customer_id,omitempty" db:"customer_id"`
	Status         string      `json:"status" db:"status"` // "pending", "completed", "cancelled"
	Subtotal       float64     `json:"subtotal" db:"subtotal"`
	TaxAmount      float64     `json:"tax_amount" db:"tax_amount"`
	DiscountAmount float64     `json:"discount_amount" db:"discount_amount"`
	TotalAmount    float64     `json:"total_amount" db:"total_amount"`
	PaymentStatus  string      `json:"payment_status" db:"payment_status"` // "pending", "paid", "refunded"
	Notes          string      `json:"notes" db:"notes"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
	Customer       *Customer   `json:"customer,omitempty"`
	Items          []OrderItem `json:"items,omitempty"`
	Payments       []Payment   `json:"payments,omitempty"`
}

// OrderItem represents an item in a sales order
type OrderItem struct {
	ID         uuid.UUID `json:"id" db:"id"`
	OrderID    uuid.UUID `json:"order_id" db:"order_id"`
	ProductID  uuid.UUID `json:"product_id" db:"product_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	UnitPrice  float64   `json:"unit_price" db:"unit_price"`
	Discount   float64   `json:"discount" db:"discount"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Product    *Product  `json:"product,omitempty"`
}

// Payment represents a payment for an order
type Payment struct {
	ID            uuid.UUID `json:"id" db:"id"`
	OrderID       uuid.UUID `json:"order_id" db:"order_id"`
	Amount        float64   `json:"amount" db:"amount"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"` // "cash", "card", "transfer", "digital_wallet"
	Reference     string    `json:"reference" db:"reference"`
	Status        string    `json:"status" db:"status"` // "pending", "completed", "failed", "refunded"
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Receipt represents a sales receipt
type Receipt struct {
	ID            uuid.UUID `json:"id" db:"id"`
	OrderID       uuid.UUID `json:"order_id" db:"order_id"`
	ReceiptNumber string    `json:"receipt_number" db:"receipt_number"`
	TotalAmount   float64   `json:"total_amount" db:"total_amount"`
	TaxAmount     float64   `json:"tax_amount" db:"tax_amount"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	Order         *Order    `json:"order,omitempty"`
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

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// UpdateCustomerRequest represents the request to update a customer
type UpdateCustomerRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	CustomerID     *string            `json:"customer_id"`
	Items          []OrderItemRequest `json:"items" validate:"required,min=1"`
	TaxAmount      float64            `json:"tax_amount"`
	DiscountAmount float64            `json:"discount_amount"`
	Notes          string             `json:"notes"`
}

// OrderItemRequest represents an item in order creation request
type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	Discount  float64   `json:"discount"`
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,min=0"`
	PaymentMethod string    `json:"payment_method" validate:"required,oneof=cash card transfer digital_wallet"`
	Reference     string    `json:"reference"`
}

// SalesReport represents sales report data
type SalesReport struct {
	TotalSales   float64        `json:"total_sales"`
	TotalOrders  int            `json:"total_orders"`
	AverageOrder float64        `json:"average_order"`
	TopProducts  []ProductSales `json:"top_products"`
	SalesByDate  []DailySales   `json:"sales_by_date"`
}

// ProductSales represents product sales data
type ProductSales struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Revenue     float64   `json:"revenue"`
}

// DailySales represents daily sales data
type DailySales struct {
	Date   string  `json:"date"`
	Sales  float64 `json:"sales"`
	Orders int     `json:"orders"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
