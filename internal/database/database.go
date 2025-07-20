package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to database")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// CreateTables creates all necessary tables for the application
func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			sku VARCHAR(100) NOT NULL UNIQUE,
			category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
			price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS inventory (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
			location VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(product_id, location)
		)`,
		`CREATE TABLE IF NOT EXISTS inventory_transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL CHECK (type IN ('in', 'out', 'adjustment')),
			quantity INTEGER NOT NULL,
			reason VARCHAR(255) NOT NULL,
			reference VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)`,
		`CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_transactions_product_id ON inventory_transactions(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_transactions_created_at ON inventory_transactions(created_at)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}
