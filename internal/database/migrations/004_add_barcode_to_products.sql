-- Migration: Add barcode_number to products and make sku nullable
ALTER TABLE products ADD COLUMN IF NOT EXISTS barcode_number VARCHAR(100);
ALTER TABLE products ALTER COLUMN sku DROP NOT NULL; 