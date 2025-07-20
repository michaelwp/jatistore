package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type OrderRepository struct {
	db *database.DB
}

func NewOrderRepository(db *database.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `
		INSERT INTO orders (id, order_number, customer_id, status, subtotal, tax_amount, discount_amount, total_amount, payment_status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	order.ID = uuid.New()
	order.CreatedAt = now
	order.UpdatedAt = now

	_, err = tx.Exec(orderQuery,
		order.ID,
		order.OrderNumber,
		order.CustomerID,
		order.Status,
		order.Subtotal,
		order.TaxAmount,
		order.DiscountAmount,
		order.TotalAmount,
		order.PaymentStatus,
		order.Notes,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Insert order items
	for i := range order.Items {
		item := &order.Items[i]
		itemQuery := `
			INSERT INTO order_items (id, order_id, product_id, quantity, unit_price, discount, total_price, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		item.ID = uuid.New()
		item.OrderID = order.ID
		item.CreatedAt = now

		_, err = tx.Exec(itemQuery,
			item.ID,
			item.OrderID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.Discount,
			item.TotalPrice,
			item.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*models.Order, error) {
	// Get order with customer
	orderQuery := `
		SELECT o.id, o.order_number, o.customer_id, o.status, o.subtotal, o.tax_amount, o.discount_amount, o.total_amount, o.payment_status, o.notes, o.created_at, o.updated_at,
		       c.id, c.name, c.email, c.phone, c.address, c.created_at, c.updated_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		WHERE o.id = $1
	`

	var order models.Order
	var customer models.Customer

	err := r.db.QueryRow(orderQuery, id).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.CustomerID,
		&order.Status,
		&order.Subtotal,
		&order.TaxAmount,
		&order.DiscountAmount,
		&order.TotalAmount,
		&order.PaymentStatus,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order.Customer = &customer

	// Get order items
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.discount, oi.total_price, oi.created_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM order_items oi
		LEFT JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var product models.Product

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Discount,
			&item.TotalPrice,
			&item.CreatedAt,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.SKU,
			&product.CategoryID,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		item.Product = &product
		items = append(items, item)
	}

	order.Items = items

	return &order, nil
}

func (r *OrderRepository) GetAll() ([]models.Order, error) {
	query := `
		SELECT o.id, o.order_number, o.customer_id, o.status, o.subtotal, o.tax_amount, o.discount_amount, o.total_amount, o.payment_status, o.notes, o.created_at, o.updated_at,
		       c.id, c.name, c.email, c.phone, c.address, c.created_at, c.updated_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var customer models.Customer

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.CustomerID,
			&order.Status,
			&order.Subtotal,
			&order.TaxAmount,
			&order.DiscountAmount,
			&order.TotalAmount,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		order.Customer = &customer
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(id uuid.UUID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (r *OrderRepository) UpdatePaymentStatus(id uuid.UUID, paymentStatus string) error {
	query := `UPDATE orders SET payment_status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.Exec(query, paymentStatus, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (r *OrderRepository) GetByCustomerID(customerID uuid.UUID) ([]models.Order, error) {
	query := `
		SELECT o.id, o.order_number, o.customer_id, o.status, o.subtotal, o.tax_amount, o.discount_amount, o.total_amount, o.payment_status, o.notes, o.created_at, o.updated_at,
		       c.id, c.name, c.email, c.phone, c.address, c.created_at, c.updated_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		WHERE o.customer_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query customer orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var customer models.Customer

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.CustomerID,
			&order.Status,
			&order.Subtotal,
			&order.TaxAmount,
			&order.DiscountAmount,
			&order.TotalAmount,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		order.Customer = &customer
		orders = append(orders, order)
	}

	return orders, nil
}
