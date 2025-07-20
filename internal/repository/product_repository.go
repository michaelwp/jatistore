package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type ProductRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, sku, category_id, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	product.ID = uuid.New()
	product.CreatedAt = now
	product.UpdatedAt = now

	_, err := r.db.Exec(query,
		product.ID,
		product.Name,
		product.Description,
		product.SKU,
		product.CategoryID,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	product := &models.Product{}
	var category models.Category

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.CategoryID,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	product.Category = &category
	return product, nil
}

func (r *ProductRepository) GetAll() ([]*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var category models.Category

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.SKU,
			&product.CategoryID,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		product.Category = &category
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) Update(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, description = $2, sku = $3, category_id = $4, price = $5, updated_at = $6
		WHERE id = $7
	`

	product.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		product.Name,
		product.Description,
		product.SKU,
		product.CategoryID,
		product.Price,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) GetBySKU(sku string) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.sku, p.category_id, p.price, p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.sku = $1
	`

	product := &models.Product{}
	var category models.Category

	err := r.db.QueryRow(query, sku).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.CategoryID,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	product.Category = &category
	return product, nil
}
