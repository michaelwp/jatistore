package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type ReceiptRepository struct {
	db *database.DB
}

func NewReceiptRepository(db *database.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) Create(receipt *models.Receipt) error {
	query := `
		INSERT INTO receipts (id, order_id, receipt_number, total_amount, tax_amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	receipt.ID = uuid.New()
	receipt.CreatedAt = time.Now()

	_, err := r.db.Exec(query,
		receipt.ID,
		receipt.OrderID,
		receipt.ReceiptNumber,
		receipt.TotalAmount,
		receipt.TaxAmount,
		receipt.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create receipt: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) GetByID(id uuid.UUID) (*models.Receipt, error) {
	query := `
		SELECT r.id, r.order_id, r.receipt_number, r.total_amount, r.tax_amount, r.created_at,
		       o.id, o.order_number, o.customer_id, o.status, o.subtotal, o.tax_amount, o.discount_amount, o.total_amount, o.payment_status, o.notes, o.created_at, o.updated_at
		FROM receipts r
		LEFT JOIN orders o ON r.order_id = o.id
		WHERE r.id = $1
	`

	var receipt models.Receipt
	var order models.Order

	err := r.db.QueryRow(query, id).Scan(
		&receipt.ID,
		&receipt.OrderID,
		&receipt.ReceiptNumber,
		&receipt.TotalAmount,
		&receipt.TaxAmount,
		&receipt.CreatedAt,
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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receipt not found")
		}
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}

	receipt.Order = &order
	return &receipt, nil
}

func (r *ReceiptRepository) GetByOrderID(orderID uuid.UUID) (*models.Receipt, error) {
	query := `
		SELECT r.id, r.order_id, r.receipt_number, r.total_amount, r.tax_amount, r.created_at,
		       o.id, o.order_number, o.customer_id, o.status, o.subtotal, o.tax_amount, o.discount_amount, o.total_amount, o.payment_status, o.notes, o.created_at, o.updated_at
		FROM receipts r
		LEFT JOIN orders o ON r.order_id = o.id
		WHERE r.order_id = $1
	`

	var receipt models.Receipt
	var order models.Order

	err := r.db.QueryRow(query, orderID).Scan(
		&receipt.ID,
		&receipt.OrderID,
		&receipt.ReceiptNumber,
		&receipt.TotalAmount,
		&receipt.TaxAmount,
		&receipt.CreatedAt,
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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receipt not found")
		}
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}

	receipt.Order = &order
	return &receipt, nil
}

func (r *ReceiptRepository) GetAll() ([]models.Receipt, error) {
	query := `
		SELECT r.id, r.order_id, r.receipt_number, r.total_amount, r.tax_amount, r.created_at,
		       o.id, o.order_number, o.customer_id, o.status, o.subtotal, o.tax_amount, o.discount_amount, o.total_amount, o.payment_status, o.notes, o.created_at, o.updated_at
		FROM receipts r
		LEFT JOIN orders o ON r.order_id = o.id
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query receipts: %w", err)
	}
	defer rows.Close()

	var receipts []models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		var order models.Order

		err := rows.Scan(
			&receipt.ID,
			&receipt.OrderID,
			&receipt.ReceiptNumber,
			&receipt.TotalAmount,
			&receipt.TaxAmount,
			&receipt.CreatedAt,
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
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan receipt: %w", err)
		}

		receipt.Order = &order
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
