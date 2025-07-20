package repository

import (
	"database/sql"
	"fmt"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
)

type CategoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `
		INSERT INTO categories (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	category.ID = uuid.New()
	category.CreatedAt = now
	category.UpdatedAt = now

	_, err := r.db.Exec(query,
		category.ID,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	category := &models.Category{}

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) GetAll() ([]*models.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := `
		UPDATE categories 
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	category.UpdatedAt = time.Now()

	result, err := r.db.Exec(query,
		category.Name,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
