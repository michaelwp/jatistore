package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type CustomerRepository struct {
	db *database.DB
}

func NewCustomerRepository(db *database.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	customer.ID = uuid.New()
	customer.CreatedAt = now
	customer.UpdatedAt = now

	_, err := r.db.Exec(query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

func (r *CustomerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {
	query := `SELECT * FROM customers WHERE id = $1`

	var customer models.Customer
	err := r.db.QueryRow(query, id).Scan(
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
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

func (r *CustomerRepository) GetByEmail(email string) (*models.Customer, error) {
	query := `SELECT * FROM customers WHERE email = $1`

	var customer models.Customer
	err := r.db.QueryRow(query, email).Scan(
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
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

func (r *CustomerRepository) GetAll() ([]models.Customer, error) {
	query := `SELECT * FROM customers ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	query := `
		UPDATE customers 
		SET name = $1, email = $2, phone = $3, address = $4, updated_at = $5
		WHERE id = $6
	`

	customer.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.UpdatedAt,
		customer.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

func (r *CustomerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

func (r *CustomerRepository) Search(query string) ([]models.Customer, error) {
	sqlQuery := `
		SELECT * FROM customers 
		WHERE name ILIKE $1 OR email ILIKE $1 OR phone ILIKE $1
		ORDER BY created_at DESC
	`

	searchTerm := "%" + query + "%"
	rows, err := r.db.Query(sqlQuery, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
