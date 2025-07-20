package handlers

import (
	"net/http"

	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	errCustomerNotFound = "customer not found"
)

type CustomerHandler struct {
	customerService *services.CustomerService
}

func NewCustomerHandler(customerService *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer with the provided information
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body models.CreateCustomerRequest true "Customer information"
// @Success 201 {object} models.APIResponse{data=models.Customer}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req models.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Customer name is required",
		})
	}

	if req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Customer email is required",
		})
	}

	customer, err := h.customerService.CreateCustomer(&req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Customer created successfully",
		Data:    customer,
	})
}

// GetCustomer godoc
// @Summary Get a customer by ID
// @Description Get customer details by customer ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.APIResponse{data=models.Customer}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid customer ID",
		})
	}

	customer, err := h.customerService.GetCustomer(id)
	if err != nil {
		if err.Error() == errCustomerNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Customer not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    customer,
	})
}

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Get a list of all customers
// @Tags customers
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Customer}
// @Failure 500 {object} models.APIResponse
// @Router /customers [get]
func (h *CustomerHandler) GetAllCustomers(c *fiber.Ctx) error {
	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    customers,
	})
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update customer information by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body models.UpdateCustomerRequest true "Updated customer information"
// @Success 200 {object} models.APIResponse{data=models.Customer}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid customer ID",
		})
	}

	var req models.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Customer name is required",
		})
	}

	if req.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Customer email is required",
		})
	}

	customer, err := h.customerService.UpdateCustomer(id, &req)
	if err != nil {
		if err.Error() == errCustomerNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Customer not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Customer updated successfully",
		Data:    customer,
	})
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer by ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid customer ID",
		})
	}

	err = h.customerService.DeleteCustomer(id)
	if err != nil {
		if err.Error() == errCustomerNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Customer not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Customer deleted successfully",
	})
}

// SearchCustomers godoc
// @Summary Search customers
// @Description Search customers by name, email, or phone
// @Tags customers
// @Produce json
// @Param q query string false "Search query"
// @Success 200 {object} models.APIResponse{data=[]models.Customer}
// @Failure 500 {object} models.APIResponse
// @Router /customers/search [get]
func (h *CustomerHandler) SearchCustomers(c *fiber.Ctx) error {
	query := c.Query("q", "")

	customers, err := h.customerService.SearchCustomers(query)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    customers,
	})
}
