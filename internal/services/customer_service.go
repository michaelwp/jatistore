package services

import (
	"fmt"

	"jatistore/internal/models"
	"jatistore/internal/repository"

	"github.com/google/uuid"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

func (s *CustomerService) CreateCustomer(req *models.CreateCustomerRequest) (*models.Customer, error) {
	// Check if customer with email already exists
	if req.Email != "" {
		existingCustomer, err := s.customerRepo.GetByEmail(req.Email)
		if err == nil && existingCustomer != nil {
			return nil, fmt.Errorf("customer with email %s already exists", req.Email)
		}
	}

	customer := &models.Customer{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	err := s.customerRepo.Create(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerService) GetCustomer(id uuid.UUID) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerService) GetAllCustomers() ([]models.Customer, error) {
	customers, err := s.customerRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return customers, nil
}

func (s *CustomerService) UpdateCustomer(id uuid.UUID, req *models.UpdateCustomerRequest) (*models.Customer, error) {
	// Check if customer exists
	existingCustomer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Check if email is being changed and if it already exists
	if req.Email != existingCustomer.Email && req.Email != "" {
		emailCustomer, err := s.customerRepo.GetByEmail(req.Email)
		if err == nil && emailCustomer != nil {
			return nil, fmt.Errorf("customer with email %s already exists", req.Email)
		}
	}

	// Update customer fields
	existingCustomer.Name = req.Name
	existingCustomer.Email = req.Email
	existingCustomer.Phone = req.Phone
	existingCustomer.Address = req.Address

	err = s.customerRepo.Update(existingCustomer)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return existingCustomer, nil
}

func (s *CustomerService) DeleteCustomer(id uuid.UUID) error {
	err := s.customerRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

func (s *CustomerService) SearchCustomers(query string) ([]models.Customer, error) {
	if query == "" {
		return s.GetAllCustomers()
	}

	customers, err := s.customerRepo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	return customers, nil
}
