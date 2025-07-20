// Package main provides the entry point for the JatiStore POS API server.
// This application serves as a comprehensive Point of Sales system with
// inventory management, customer management, order processing, payments, and receipts.
package main

// @title JatiStore POS API
// @version 2.0
// @description RESTful API for Point of Sales (POS) system using Go, Fiber, and PostgreSQL. Includes inventory management, customer management, order processing, payments, and receipts.
// @contact.name API Support
// @contact.email support@jatistore.local
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"log"
	"os"
	"strings"

	"jatistore/internal/config"
	"jatistore/internal/database"
	"jatistore/internal/handlers"
	"jatistore/internal/middleware"
	"jatistore/internal/repository"
	"jatistore/internal/router"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	docs "jatistore/docs"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Dynamically set Swagger host
	setSwaggerHost(cfg)

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create database tables
	if err := db.CreateTables(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v", closeErr)
		}
		log.Fatal("Failed to create database tables:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	receiptRepo := repository.NewReceiptRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	inventoryService := services.NewInventoryService(inventoryRepo)
	customerService := services.NewCustomerService(customerRepo)
	orderService := services.NewOrderService(orderRepo, productRepo, customerRepo, paymentRepo, receiptRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Initialize authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Create handlers instance
	handlers := router.NewHandlers(authHandler, productHandler, categoryHandler, inventoryHandler, customerHandler, orderHandler)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes
	router.SetupRoutes(app, handlers, authMiddleware)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Printf("Server error: %v", err)
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v", closeErr)
		}
		os.Exit(1)
	}
}

func setSwaggerHost(cfg *config.Config) {
	host := ""
	if cfg.BaseURL != "" {
		// Remove protocol if present
		base := cfg.BaseURL
		if strings.HasPrefix(base, "http://") {
			base = strings.TrimPrefix(base, "http://")
		} else if strings.HasPrefix(base, "https://") {
			base = strings.TrimPrefix(base, "https://")
		}
		// If port is not included, append it
		if !strings.Contains(base, ":") && cfg.Port != "80" && cfg.Port != "443" {
			base = base + ":" + cfg.Port
		}
		host = base
	} else {
		host = "localhost"
	}
	docs.SwaggerInfo.Host = host
}
