package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type PaymentRepository struct {
	db *database.DB
}

func NewPaymentRepository(db *database.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, amount, payment_method, reference, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	payment.ID = uuid.New()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	_, err := r.db.Exec(query,
		payment.ID,
		payment.OrderID,
		payment.Amount,
		payment.PaymentMethod,
		payment.Reference,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *PaymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	query := `SELECT * FROM payments WHERE id = $1`

	var payment models.Payment
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.Reference,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) GetByOrderID(orderID uuid.UUID) ([]models.Payment, error) {
	query := `SELECT * FROM payments WHERE order_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.Reference,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) UpdateStatus(id uuid.UUID, status string) error {
	query := `UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (r *PaymentRepository) GetTotalPaidByOrderID(orderID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE order_id = $1 AND status = 'completed'`

	var total float64
	err := r.db.QueryRow(query, orderID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total paid: %w", err)
	}

	return total, nil
}
