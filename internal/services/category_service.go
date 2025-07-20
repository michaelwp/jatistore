package services

import (
	"fmt"

	"jatistore/internal/models"
	"jatistore/internal/repository"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(req *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Get the created category
	createdCategory, err := s.categoryRepo.GetByID(category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created category: %w", err)
	}

	return createdCategory, nil
}

func (s *CategoryService) GetCategoryByID(id string) (*models.Category, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

func (s *CategoryService) GetAllCategories() ([]*models.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (s *CategoryService) UpdateCategory(id string, req *models.UpdateCategoryRequest) (*models.Category, error) {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Get existing category
	existingCategory, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing category: %w", err)
	}

	// Update category fields
	existingCategory.Name = req.Name
	existingCategory.Description = req.Description

	if err := s.categoryRepo.Update(existingCategory); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Get the updated category
	updatedCategory, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated category: %w", err)
	}

	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	if err := s.categoryRepo.Delete(categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
