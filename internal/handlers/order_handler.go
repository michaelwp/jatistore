package handlers

import (
	"net/http"

	"jatistore/internal/models"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	errOrderNotFound = "order not found"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new sales order with items
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param order body models.CreateOrderRequest true "Order information"
// @Success 201 {object} models.APIResponse{data=models.Order}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if len(req.Items) == 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Order must have at least one item",
		})
	}

	for i, item := range req.Items {
		if item.Quantity <= 0 {
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error:   "Item quantity must be greater than 0",
			})
		}
		if item.Discount < 0 {
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error:   "Item discount cannot be negative",
			})
		}
		req.Items[i] = item
	}

	if req.TaxAmount < 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Tax amount cannot be negative",
		})
	}

	if req.DiscountAmount < 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Discount amount cannot be negative",
		})
	}

	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
	})
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get order details by order ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Success 200 {object} models.APIResponse{data=models.Order}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid order ID",
		})
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		if err.Error() == errOrderNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Order not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    order,
	})
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get a list of all orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.APIResponse{data=[]models.Order}
// @Failure 500 {object} models.APIResponse
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    orders,
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Param status body map[string]string true "Order status"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid order ID",
		})
	}

	var req map[string]string
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	status, exists := req["status"]
	if !exists {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Status is required",
		})
	}

	err = h.orderService.UpdateOrderStatus(id, status)
	if err != nil {
		if err.Error() == errOrderNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Order not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Order status updated successfully",
	})
}

// ProcessPayment godoc
// @Summary Process payment for an order
// @Description Process a payment for an order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Param payment body models.CreatePaymentRequest true "Payment information"
// @Success 200 {object} models.APIResponse{data=models.Payment}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /orders/{id}/payments [post]
func (h *OrderHandler) ProcessPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid order ID",
		})
	}

	var req models.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.Amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Payment amount must be greater than 0",
		})
	}

	validPaymentMethods := map[string]bool{
		"cash":           true,
		"card":           true,
		"transfer":       true,
		"digital_wallet": true,
	}

	if !validPaymentMethods[req.PaymentMethod] {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid payment method",
		})
	}

	payment, err := h.orderService.ProcessPayment(id, &req)
	if err != nil {
		if err.Error() == errOrderNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Order not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Payment processed successfully",
		Data:    payment,
	})
}

// GenerateReceipt godoc
// @Summary Generate receipt for an order
// @Description Generate a receipt for a paid order
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Success 200 {object} models.APIResponse{data=models.Receipt}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /orders/{id}/receipt [post]
func (h *OrderHandler) GenerateReceipt(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid order ID",
		})
	}

	receipt, err := h.orderService.GenerateReceipt(id)
	if err != nil {
		if err.Error() == errOrderNotFound {
			return c.Status(http.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Error:   "Order not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Receipt generated successfully",
		Data:    receipt,
	})
}

// GetOrdersByCustomer godoc
// @Summary Get orders by customer
// @Description Get all orders for a specific customer
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param customerId path string true "Customer ID"
// @Success 200 {object} models.APIResponse{data=[]models.Order}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /customers/{customerId}/orders [get]
func (h *OrderHandler) GetOrdersByCustomer(c *fiber.Ctx) error {
	customerID, err := uuid.Parse(c.Params("customerId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid customer ID",
		})
	}

	orders, err := h.orderService.GetOrdersByCustomer(customerID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    orders,
	})
}
