package main

// @title JatiStore Inventory API
// @version 1.0
// @description RESTful API for inventory management using Go, Fiber, and PostgreSQL.
// @contact.name API Support
// @contact.email support@jatistore.local
// @host localhost:8080
// @BasePath /api/v1

import (
	"log"
	"os"

	"jatistore/internal/config"
	"jatistore/internal/database"
	"jatistore/internal/handlers"
	"jatistore/internal/middleware"
	"jatistore/internal/repository"
	"jatistore/internal/router"
	"jatistore/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create database tables
	if err := db.CreateTables(); err != nil {
		log.Fatal("Failed to create database tables:", err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)

	// Initialize services
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	inventoryService := services.NewInventoryService(inventoryRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Create handlers instance
	handlers := router.NewHandlers(productHandler, categoryHandler, inventoryHandler)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes
	router.SetupRoutes(app, handlers)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
