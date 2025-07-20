package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type InventoryRepository struct {
	db *database.DB
}

func NewInventoryRepository(db *database.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	query := `
		INSERT INTO inventory (id, product_id, quantity, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	inventory.ID = uuid.New()
	inventory.CreatedAt = now
	inventory.UpdatedAt = now

	_, err := r.db.Exec(query,
		inventory.ID,
		inventory.ProductID,
		inventory.Quantity,
		inventory.Location,
		inventory.CreatedAt,
		inventory.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetByID(id uuid.UUID) (*models.Inventory, error) {
	query := `
		SELECT i.id, i.product_id, i.quantity, i.location, i.created_at, i.updated_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM inventory i
		LEFT JOIN products p ON i.product_id = p.id
		WHERE i.id = $1
	`

	inventory := &models.Inventory{}
	var product models.Product

	err := r.db.QueryRow(query, id).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.Location,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	inventory.Product = &product
	return inventory, nil
}

func (r *InventoryRepository) GetAll() ([]*models.Inventory, error) {
	query := `
		SELECT i.id, i.product_id, i.quantity, i.location, i.created_at, i.updated_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM inventory i
		LEFT JOIN products p ON i.product_id = p.id
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory: %w", err)
	}
	defer rows.Close()

	var inventories []*models.Inventory
	for rows.Next() {
		inventory := &models.Inventory{}
		var product models.Product

		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.Quantity,
			&inventory.Location,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
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
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}

		inventory.Product = &product
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

func (r *InventoryRepository) Update(inventory *models.Inventory) error {
	query := `
		UPDATE inventory 
		SET quantity = $1, location = $2, updated_at = $3
		WHERE id = $4
	`

	inventory.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		inventory.Quantity,
		inventory.Location,
		inventory.UpdatedAt,
		inventory.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory not found")
	}

	return nil
}

func (r *InventoryRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM inventory WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory not found")
	}

	return nil
}

func (r *InventoryRepository) GetByProductID(productID uuid.UUID) ([]*models.Inventory, error) {
	query := `
		SELECT i.id, i.product_id, i.quantity, i.location, i.created_at, i.updated_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM inventory i
		LEFT JOIN products p ON i.product_id = p.id
		WHERE i.product_id = $1
		ORDER BY i.location ASC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory by product ID: %w", err)
	}
	defer rows.Close()

	var inventories []*models.Inventory
	for rows.Next() {
		inventory := &models.Inventory{}
		var product models.Product

		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.Quantity,
			&inventory.Location,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
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
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}

		inventory.Product = &product
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

func (r *InventoryRepository) CreateTransaction(transaction *models.InventoryTransaction) error {
	query := `
		INSERT INTO inventory_transactions (id, product_id, type, quantity, reason, reference, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	transaction.ID = uuid.New()
	transaction.CreatedAt = now

	_, err := r.db.Exec(query,
		transaction.ID,
		transaction.ProductID,
		transaction.Type,
		transaction.Quantity,
		transaction.Reason,
		transaction.Reference,
		transaction.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create inventory transaction: %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetTransactionsByProductID(productID uuid.UUID) ([]*models.InventoryTransaction, error) {
	query := `
		SELECT it.id, it.product_id, it.type, it.quantity, it.reason, it.reference, it.created_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM inventory_transactions it
		LEFT JOIN products p ON it.product_id = p.id
		WHERE it.product_id = $1
		ORDER BY it.created_at DESC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.InventoryTransaction
	for rows.Next() {
		transaction := &models.InventoryTransaction{}
		var product models.Product

		err := rows.Scan(
			&transaction.ID,
			&transaction.ProductID,
			&transaction.Type,
			&transaction.Quantity,
			&transaction.Reason,
			&transaction.Reference,
			&transaction.CreatedAt,
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
			return nil, fmt.Errorf("failed to scan inventory transaction: %w", err)
		}

		transaction.Product = &product
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *InventoryRepository) GetByProductIDString(productID string) ([]*models.Inventory, error) {
	query := `
		SELECT i.id, i.product_id, i.quantity, i.location, i.created_at, i.updated_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM inventory i
		LEFT JOIN products p ON i.product_id = p.id
		WHERE i.product_id = $1
		ORDER BY i.location ASC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory by product ID: %w", err)
	}
	defer rows.Close()

	var inventories []*models.Inventory
	for rows.Next() {
		inventory := &models.Inventory{}
		var product models.Product

		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&inventory.Quantity,
			&inventory.Location,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
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
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}

		inventory.Product = &product
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

func (r *InventoryRepository) CreateTransactionString(transaction *models.InventoryTransaction) error {
	query := `
		INSERT INTO inventory_transactions (id, product_id, type, quantity, reason, reference, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	`

	_, err := r.db.Exec(query,
		transaction.ProductID,
		transaction.Type,
		transaction.Quantity,
		transaction.Reason,
		transaction.Reference,
	)

	if err != nil {
		return fmt.Errorf("failed to create inventory transaction: %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetTransactionsByProductIDString(productID string) ([]*models.InventoryTransaction, error) {
	query := `
		SELECT it.id, it.product_id, it.type, it.quantity, it.reason, it.reference, it.created_at,
		       p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at
		FROM inventory_transactions it
		LEFT JOIN products p ON it.product_id = p.id
		WHERE it.product_id = $1
		ORDER BY it.created_at DESC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*models.InventoryTransaction
	for rows.Next() {
		transaction := &models.InventoryTransaction{}
		var product models.Product

		err := rows.Scan(
			&transaction.ID,
			&transaction.ProductID,
			&transaction.Type,
			&transaction.Quantity,
			&transaction.Reason,
			&transaction.Reference,
			&transaction.CreatedAt,
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
			return nil, fmt.Errorf("failed to scan inventory transaction: %w", err)
		}

		transaction.Product = &product
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
