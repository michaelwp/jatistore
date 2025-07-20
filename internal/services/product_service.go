package services

import (
	"fmt"

	"jatistore/internal/models"
	"jatistore/internal/repository"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	// Parse category ID
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Check if SKU already exists (only if SKU is provided)
	if req.SKU != "" {
		existingProduct, _ := s.productRepo.GetBySKU(req.SKU)
		if existingProduct != nil {
			return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
		}
	}

	product := &models.Product{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		BarcodeNumber: req.BarcodeNumber,
		CategoryID:    categoryID,
		Price:         req.Price,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Get the created product with category information
	createdProduct, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created product: %w", err)
	}

	return createdProduct, nil
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetAllProducts() ([]*models.Product, error) {
	products, err := s.productRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (s *ProductService) UpdateProduct(id string, req *models.UpdateProductRequest) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Parse category ID
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Get existing product
	existingProduct, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing product: %w", err)
	}

	// Check if SKU is being changed and if it already exists (only if SKU is provided)
	if req.SKU != "" && existingProduct.SKU != req.SKU {
		productWithSKU, _ := s.productRepo.GetBySKU(req.SKU)
		if productWithSKU != nil && productWithSKU.ID != productID {
			return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
		}
	}

	// Update product fields
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.SKU = req.SKU
	existingProduct.BarcodeNumber = req.BarcodeNumber
	existingProduct.CategoryID = categoryID
	existingProduct.Price = req.Price

	if err := s.productRepo.Update(existingProduct); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Get the updated product with category information
	updatedProduct, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated product: %w", err)
	}

	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	if err := s.productRepo.Delete(productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
