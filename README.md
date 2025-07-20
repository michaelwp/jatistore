# JatiStore - Inventory Management System

A modern, robust inventory management system built with Go, PostgreSQL, and Fiber web framework. JatiStore provides comprehensive product, category, and inventory management with full transaction tracking and audit capabilities.

## 🚀 Features

- **Product Management**: Complete CRUD operations for products with category organization
- **Category Management**: Hierarchical product categorization system
- **Inventory Management**: Real-time stock level tracking across multiple locations
- **Transaction Tracking**: Complete audit trail of all inventory movements (in, out, adjustments)
- **RESTful API**: Clean, intuitive API endpoints with comprehensive documentation
- **PostgreSQL Database**: Robust, scalable database with proper indexing and constraints
- **Swagger UI**: Interactive API documentation with live testing capabilities
- **Environment Configuration**: Flexible configuration management for different environments

## 📁 Project Structure

```
jatistore/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── env.example             # Environment variables example
├── README.md               # This file
├── Makefile                # Build and management commands
├── docs/                   # Swagger API documentation (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── bin/                    # Compiled binary output
└── internal/               # Internal application code
    ├── config/             # Configuration management
    ├── database/           # Database connection and setup
    │   └── migrations/     # Database migration files
    ├── models/             # Data models and structures
    ├── repository/         # Data access layer
    ├── services/           # Business logic layer
    ├── handlers/           # HTTP request handlers
    ├── middleware/         # HTTP middleware
    └── router/             # Route definitions
```

## ⚙️ Prerequisites

- **Go 1.24** or higher
- **PostgreSQL 12** or higher
- **Git**

## 🛠️ Makefile Commands

The project includes a comprehensive `Makefile` for common development tasks:

| Command         | Description                                                      |
|-----------------|------------------------------------------------------------------|
| `make build`    | Build the application binary into the `bin/` directory           |
| `make run`      | Run the application using `go run main.go`                       |
| `make swag`     | Generate Swagger API documentation into the `docs/` directory    |
| `make tidy`     | Clean up and verify Go module dependencies                       |
| `make clean`    | Remove the `bin/` and `docs/` directories                        |
| `make lint`     | Run golangci-lint for code quality checks                        |
| `make pre-commit` | Run pre-commit checks including linting                        |
| `make install-hooks` | Install git hooks for automated checks                      |
| `make migrate-up`   | Apply all database migrations (requires `migrate` tool)      |
| `make migrate-down` | Roll back the last database migration (requires `migrate`)   |

> **Note:**
> - `make migrate-up` and `make migrate-down` require the [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool to be installed.
> - `make swag` requires the [swag CLI](https://github.com/swaggo/swag) to be installed.

## 🚀 Quick Start

### 1. Clone and Setup
```bash
git clone <repository-url>
cd jatistore
go mod tidy
```

### 2. Database Setup
```bash
# Create PostgreSQL database
createdb jatistore

# Or using psql
psql -U postgres -c "CREATE DATABASE jatistore;"
```

### 3. Environment Configuration
```bash
cp env.example .env
# Edit .env file with your database credentials
```

Example `.env` configuration:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=y_username
DB_PASSWORD=y_password
DB_NAME=y_database
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

### 4. Generate API Documentation
```bash
make swag
```

### 5. Run the Application
```bash
make run
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Swagger UI
After running the app and generating docs, visit:
```
http://localhost:8080/swagger/index.html
```

Here you can view and interact with the complete API documentation.

## 🔌 API Endpoints

### Health Check
- `GET /health` - Check if the API is running

### Categories
- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/:id` - Get category by ID
- `POST /api/v1/categories` - Create a new category
- `PUT /api/v1/categories/:id` - Update a category
- `DELETE /api/v1/categories/:id` - Delete a category

### Products
- `GET /api/v1/products` - Get all products
- `GET /api/v1/products/:id` - Get product by ID
- `POST /api/v1/products` - Create a new product
- `PUT /api/v1/products/:id` - Update a product
- `DELETE /api/v1/products/:id` - Delete a product

### Inventory
- `GET /api/v1/inventory` - Get all inventory records
- `GET /api/v1/inventory/:id` - Get inventory by ID
- `POST /api/v1/inventory` - Create a new inventory record
- `PUT /api/v1/inventory/:id` - Update an inventory record
- `DELETE /api/v1/inventory/:id` - Delete an inventory record
- `POST /api/v1/inventory/adjust` - Adjust stock levels and record transactions

## 💡 API Examples

### Create a Category
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and accessories"
  }'
```

### Create a Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone model with advanced features",
    "sku": "IPHONE-15-128GB",
    "category_id": "category-uuid-here",
    "price": 999.99
  }'
```

### Create Inventory Record
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-here",
    "quantity": 50,
    "location": "Warehouse A"
  }'
```

### Adjust Stock (Record Transaction)
```bash
# Add stock (incoming shipment)
curl -X POST http://localhost:8080/api/v1/inventory/adjust \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-here",
    "quantity": 25,
    "type": "in",
    "reason": "New shipment received",
    "reference": "PO-2024-001"
  }'

# Remove stock (sale or loss)
curl -X POST http://localhost:8080/api/v1/inventory/adjust \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-here",
    "quantity": 5,
    "type": "out",
    "reason": "Customer order fulfilled",
    "reference": "SO-2024-005"
  }'

# Manual adjustment (stock count correction)
curl -X POST http://localhost:8080/api/v1/inventory/adjust \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-here",
    "quantity": 45,
    "type": "adjustment",
    "reason": "Physical count correction",
    "reference": "STOCK-COUNT-2024-01"
  }'
```

## 🗄️ Database Schema

The application automatically creates the following tables with proper relationships and constraints:

### Core Tables
- **categories**: Product categories with unique names
- **products**: Product information linked to categories
- **inventory**: Stock levels and locations (unique constraint on product_id + location)
- **inventory_transactions**: Complete audit trail of all stock movements

### Key Features
- **Foreign Key Constraints**: Proper referential integrity
- **Unique Constraints**: Prevent duplicate entries where appropriate
- **Check Constraints**: Ensure data validity (e.g., non-negative quantities)
- **Indexes**: Optimized for common query patterns
- **Cascade Deletes**: Automatic cleanup of related records

## 🔄 Inventory Transactions

The system automatically tracks all inventory movements through the `inventory_transactions` table:

### Transaction Types
- **`in`**: Stock added (shipments, returns, etc.)
- **`out`**: Stock removed (sales, damage, etc.)
- **`adjustment`**: Manual correction (stock counts, corrections)

### Transaction Fields
- **`product_id`**: Reference to the product
- **`type`**: Transaction type (in/out/adjustment)
- **`quantity`**: Amount of stock moved
- **`reason`**: Human-readable reason for the transaction
- **`reference`**: External reference (PO number, SO number, etc.)
- **`created_at`**: Timestamp of the transaction

### Benefits
- **Complete Audit Trail**: Track every inventory change
- **Compliance**: Meet regulatory and audit requirements
- **Troubleshooting**: Easily identify and investigate issues
- **Reporting**: Generate detailed inventory movement reports

## 🛠️ Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
make build
```

### Database Migrations
The application automatically creates tables on startup. For production environments, consider using a proper migration tool like `golang-migrate`.

### Code Structure
The application follows clean architecture principles:
- **Handlers**: HTTP request/response handling
- **Services**: Business logic and validation
- **Repository**: Data access and persistence
- **Models**: Data structures and validation rules

## ⚙️ Environment Variables

| Variable      | Description                  | Default      | Required |
|---------------|------------------------------|--------------|----------|
| `DB_HOST`     | Database host                | `localhost`  | Yes      |
| `DB_PORT`     | Database port                | `5432`       | Yes      |
| `DB_USER`     | Database user                | `y_username` | Yes      |
| `DB_PASSWORD` | Database password            | `y_password` | Yes      |
| `DB_NAME`     | Database name                | `y_database` | Yes      |
| `PORT`        | Server port                  | `8080`       | No       |
| `ENVIRONMENT` | Application environment      | `development`| No       |
| `LOG_LEVEL`   | Logging level                | `info`       | No       |

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines
- Follow Go coding standards and conventions
- Add tests for new functionality
- Update documentation as needed
- Ensure all tests pass before submitting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [API Documentation](http://localhost:8080/swagger/index.html)
2. Review the [Issues](../../issues) page
3. Create a new issue with detailed information

---

**Built with ❤️ using Go, PostgreSQL, and Fiber** 