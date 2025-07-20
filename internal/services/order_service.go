package services

import (
	"fmt"

	"jatistore/internal/models"
	"jatistore/internal/repository"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	productRepo  *repository.ProductRepository
	customerRepo *repository.CustomerRepository
	paymentRepo  *repository.PaymentRepository
	receiptRepo  *repository.ReceiptRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	customerRepo *repository.CustomerRepository,
	paymentRepo *repository.PaymentRepository,
	receiptRepo *repository.ReceiptRepository,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
		paymentRepo:  paymentRepo,
		receiptRepo:  receiptRepo,
	}
}

func (s *OrderService) CreateOrder(req *models.CreateOrderRequest) (*models.Order, error) {
	// Validate customer if provided
	var customerID *uuid.UUID
	if req.CustomerID != nil {
		customerUUID, err := uuid.Parse(*req.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("invalid customer ID: %w", err)
		}

		_, err = s.customerRepo.GetByID(customerUUID)
		if err != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
		customerID = &customerUUID
	}

	// Process order items and calculate totals
	var orderItems []models.OrderItem
	var subtotal float64

	for _, itemReq := range req.Items {
		// Get product details
		product, err := s.productRepo.GetByID(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		// Calculate item total
		itemTotal := (product.Price * float64(itemReq.Quantity)) - itemReq.Discount
		if itemTotal < 0 {
			itemTotal = 0
		}

		orderItem := models.OrderItem{
			ProductID:  itemReq.ProductID,
			Quantity:   itemReq.Quantity,
			UnitPrice:  product.Price,
			Discount:   itemReq.Discount,
			TotalPrice: itemTotal,
		}

		orderItems = append(orderItems, orderItem)
		subtotal += itemTotal
	}

	// Calculate total amount
	totalAmount := subtotal + req.TaxAmount - req.DiscountAmount
	if totalAmount < 0 {
		totalAmount = 0
	}

	order := &models.Order{
		CustomerID:     customerID,
		Status:         "pending",
		Subtotal:       subtotal,
		TaxAmount:      req.TaxAmount,
		DiscountAmount: req.DiscountAmount,
		TotalAmount:    totalAmount,
		PaymentStatus:  "pending",
		Notes:          req.Notes,
		Items:          orderItems,
	}

	err := s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrder(id uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(id uuid.UUID, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	err := s.orderRepo.UpdateStatus(id, status)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *OrderService) ProcessPayment(orderID uuid.UUID, req *models.CreatePaymentRequest) (*models.Payment, error) {
	// Validate order exists
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check if payment amount is valid
	if req.Amount <= 0 {
		return nil, fmt.Errorf("payment amount must be greater than 0")
	}

	// Get total paid so far
	totalPaid, err := s.paymentRepo.GetTotalPaidByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total paid: %w", err)
	}

	// Check if payment exceeds order total
	if totalPaid+req.Amount > order.TotalAmount {
		return nil, fmt.Errorf("payment amount exceeds order total")
	}

	payment := &models.Payment{
		OrderID:       orderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Reference:     req.Reference,
		Status:        "completed",
	}

	err = s.paymentRepo.Create(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Update order payment status if fully paid
	if totalPaid+req.Amount >= order.TotalAmount {
		err = s.orderRepo.UpdatePaymentStatus(orderID, "paid")
		if err != nil {
			return nil, fmt.Errorf("failed to update payment status: %w", err)
		}
	}

	return payment, nil
}

func (s *OrderService) GenerateReceipt(orderID uuid.UUID) (*models.Receipt, error) {
	// Get order details
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check if order is paid
	if order.PaymentStatus != "paid" {
		return nil, fmt.Errorf("cannot generate receipt for unpaid order")
	}

	// Check if receipt already exists
	existingReceipt, err := s.receiptRepo.GetByOrderID(orderID)
	if err == nil && existingReceipt != nil {
		return existingReceipt, nil
	}

	receipt := &models.Receipt{
		OrderID:     orderID,
		TotalAmount: order.TotalAmount,
		TaxAmount:   order.TaxAmount,
	}

	err = s.receiptRepo.Create(receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to create receipt: %w", err)
	}

	return receipt, nil
}

func (s *OrderService) GetOrdersByCustomer(customerID uuid.UUID) ([]models.Order, error) {
	orders, err := s.orderRepo.GetByCustomerID(customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer orders: %w", err)
	}

	return orders, nil
}
