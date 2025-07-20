# JatiStore - Point of Sales (POS) System

A modern, robust Point of Sales system built with Go, PostgreSQL, and Fiber web framework. JatiStore provides comprehensive inventory management, customer management, order processing, payment handling, and receipt generation with full transaction tracking and audit capabilities.

## 🚀 Features

- **Product Management**: Complete CRUD operations for products with category organization
- **Category Management**: Hierarchical product categorization system
- **Inventory Management**: Real-time stock level tracking across multiple locations
- **Customer Management**: Complete customer database with search capabilities
- **Order Processing**: Create and manage sales orders with multiple items
- **Payment Processing**: Support for multiple payment methods (cash, card, transfer, digital wallet)
- **Receipt Generation**: Automatic receipt generation for completed orders
- **Transaction Tracking**: Complete audit trail of all inventory movements and sales
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

### Customers
- `GET /api/v1/customers` - Get all customers
- `GET /api/v1/customers/search` - Search customers by name, email, or phone
- `GET /api/v1/customers/:id` - Get customer by ID
- `POST /api/v1/customers` - Create a new customer
- `PUT /api/v1/customers/:id` - Update a customer
- `DELETE /api/v1/customers/:id` - Delete a customer

### Orders
- `GET /api/v1/orders` - Get all orders
- `GET /api/v1/orders/:id` - Get order by ID
- `POST /api/v1/orders` - Create a new order
- `PUT /api/v1/orders/:id/status` - Update order status
- `POST /api/v1/orders/:id/payments` - Process payment for an order
- `POST /api/v1/orders/:id/receipt` - Generate receipt for an order
- `GET /api/v1/customers/:customerId/orders` - Get orders by customer

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

### Create a Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "address": "123 Main St, City, State 12345"
  }'
```

### Create an Order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-uuid-here",
    "items": [
      {
        "product_id": "product-uuid-here",
        "quantity": 2,
        "discount": 10.00
      }
    ],
    "tax_amount": 15.00,
    "discount_amount": 5.00,
    "notes": "Customer requested express delivery"
  }'
```

### Process Payment
```bash
curl -X POST http://localhost:8080/api/v1/orders/order-uuid-here/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 150.00,
    "payment_method": "card",
    "reference": "TXN-123456"
  }'
```

### Generate Receipt
```bash
curl -X POST http://localhost:8080/api/v1/orders/order-uuid-here/receipt \
  -H "Content-Type: application/json"
```
```

## 🗄️ Database Schema

The application automatically creates the following tables with proper relationships and constraints:

### Core Tables
- **categories**: Product categories with unique names
- **products**: Product information linked to categories
- **inventory**: Stock levels and locations (unique constraint on product_id + location)
- **inventory_transactions**: Complete audit trail of all stock movements
- **customers**: Customer information with unique email addresses
- **orders**: Sales orders with customer association and status tracking
- **order_items**: Individual items within orders with pricing and discounts
- **payments**: Payment records for orders with multiple payment method support
- **receipts**: Receipt records for completed orders

### Key Features
- **Foreign Key Constraints**: Proper referential integrity
- **Unique Constraints**: Prevent duplicate entries where appropriate
- **Check Constraints**: Ensure data validity (e.g., non-negative quantities)
- **Indexes**: Optimized for common query patterns
- **Cascade Deletes**: Automatic cleanup of related records
- **Automatic Numbering**: Order and receipt numbers generated automatically
- **Transaction Support**: Database transactions for data consistency

## 🔄 Inventory Transactions

The system automatically tracks all inventory movements through the `inventory_transactions` table:

### Transaction Types
- **`in`**: Stock added (shipments, returns, etc.)
- **`out`**: Stock removed (sales, damage, etc.)
- **`adjustment`**: Manual stock corrections (physical counts, etc.)

## 💳 Payment Processing

The POS system supports multiple payment methods and tracks payment status:

### Payment Methods
- **`cash`**: Cash payments
- **`card`**: Credit/debit card payments
- **`transfer`**: Bank transfer payments
- **`digital_wallet`**: Digital wallet payments (e.g., PayPal, Apple Pay)

### Payment Status
- **`pending`**: Payment initiated but not completed
- **`completed`**: Payment successfully processed
- **`failed`**: Payment processing failed
- **`refunded`**: Payment has been refunded

## 📋 Order Management

### Order Status
- **`pending`**: Order created but not yet processed
- **`completed`**: Order has been fulfilled
- **`cancelled`**: Order has been cancelled

### Order Features
- **Automatic Numbering**: Orders get unique numbers (ORD-1000, ORD-1001, etc.)
- **Customer Association**: Orders can be linked to customers (optional)
- **Item Management**: Multiple items per order with individual pricing
- **Discounts**: Item-level and order-level discounts
- **Tax Calculation**: Support for tax amounts
- **Payment Tracking**: Track payment status separately from order status
- **`adjustment`**: Manual correction (stock counts, corrections)

## 🔄 Complete POS Workflow

Here's a typical workflow for processing a sale in the POS system:

### 1. Setup (One-time)
```bash
# Create categories
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics", "description": "Electronic devices"}'

# Create products
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "sku": "IPHONE-15-128GB",
    "category_id": "category-uuid-here",
    "price": 999.99
  }'

# Add inventory
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-here",
    "quantity": 50,
    "location": "Main Store"
  }'
```

### 2. Customer Management
```bash
# Create customer (optional - can create anonymous orders)
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "address": "123 Main St"
  }'
```

### 3. Create Order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-uuid-here",
    "items": [
      {
        "product_id": "product-uuid-here",
        "quantity": 1,
        "discount": 50.00
      }
    ],
    "tax_amount": 75.00,
    "discount_amount": 25.00,
    "notes": "Customer requested gift wrapping"
  }'
```

### 4. Process Payment
```bash
curl -X POST http://localhost:8080/api/v1/orders/order-uuid-here/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000.00,
    "payment_method": "card",
    "reference": "TXN-123456"
  }'
```

### 5. Generate Receipt
```bash
curl -X POST http://localhost:8080/api/v1/orders/order-uuid-here/receipt \
  -H "Content-Type: application/json"
```

### 6. Update Order Status
```bash
curl -X PUT http://localhost:8080/api/v1/orders/order-uuid-here/status \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## 🛠️ Troubleshooting

### Common Issues

#### Database Connection Issues
- **Error**: `failed to connect to database`
- **Solution**: Check your `.env` file and ensure PostgreSQL is running
- **Command**: `pg_isready -h localhost -p 5432`

#### Port Already in Use
- **Error**: `address already in use`
- **Solution**: Change the port in your `.env` file or kill the existing process
- **Command**: `lsof -ti:8080 | xargs kill -9`

#### Missing Dependencies
- **Error**: `go: module not found`
- **Solution**: Run `go mod tidy` to download dependencies

#### Swagger Documentation Issues
- **Error**: Swagger UI not loading
- **Solution**: Regenerate docs with `swag init -g main.go -o docs`

### Performance Tips

1. **Database Indexes**: The system automatically creates indexes for common queries
2. **Connection Pooling**: Configured for optimal performance with 25 connections
3. **Caching**: Consider adding Redis for session management in production
4. **Load Balancing**: Use multiple instances behind a load balancer for high traffic

## 🔒 Security Considerations

### Production Deployment

1. **Environment Variables**: Never commit `.env` files to version control
2. **Database Security**: Use strong passwords and restrict database access
3. **HTTPS**: Always use HTTPS in production
4. **Rate Limiting**: Implement rate limiting for API endpoints
5. **Authentication**: Add JWT or OAuth2 authentication for production use
6. **Input Validation**: All inputs are validated, but consider additional sanitization

### Data Backup

```bash
# Backup database
pg_dump jatistore > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore database
psql jatistore < backup_file.sql
```

## 📊 Monitoring and Logging

### Application Logs
The application logs important events to stdout. In production, consider:
- Structured logging with JSON format
- Log aggregation (ELK stack, Fluentd)
- Application performance monitoring (APM)

### Database Monitoring
- Monitor query performance with `EXPLAIN ANALYZE`
- Set up database connection monitoring
- Configure alerts for disk space and connection limits

## 🚀 Deployment

### Docker Deployment
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Environment Variables for Production
```env
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=jatistore
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info
```

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

## 📡 API Response Format

All API endpoints return a consistent response format:

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error description"
}
```

### HTTP Status Codes
- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request data
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

## 🔍 Advanced Queries

### Search Customers
```bash
# Search by name, email, or phone
curl "http://localhost:8080/api/v1/customers/search?q=john"
```

### Get Customer Orders
```bash
# Get all orders for a specific customer
curl "http://localhost:8080/api/v1/customers/customer-uuid-here/orders"
```

### Update Order Status
```bash
# Mark order as completed
curl -X PUT http://localhost:8080/api/v1/orders/order-uuid-here/status \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## 📈 Business Intelligence

### Key Metrics Available
- **Total Sales**: Sum of all completed orders
- **Order Count**: Number of orders by status
- **Customer Analytics**: Customer order history
- **Product Performance**: Most sold products
- **Payment Analytics**: Payment method distribution
- **Inventory Levels**: Current stock levels

### Sample Queries
```sql
-- Total sales today
SELECT SUM(total_amount) FROM orders 
WHERE DATE(created_at) = CURRENT_DATE 
AND payment_status = 'paid';

-- Top selling products
SELECT p.name, SUM(oi.quantity) as total_sold
FROM order_items oi
JOIN products p ON oi.product_id = p.id
GROUP BY p.id, p.name
ORDER BY total_sold DESC
LIMIT 10;

-- Customer order history
SELECT c.name, COUNT(o.id) as order_count, SUM(o.total_amount) as total_spent
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id
GROUP BY c.id, c.name
ORDER BY total_spent DESC;
```

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

## 🔄 Upgrade from Inventory to POS

This application has been upgraded from a simple inventory management system to a full-featured Point of Sales (POS) system. Here's what changed:

### What's New
- **Customer Management**: Complete customer database with search
- **Order Processing**: Sales order creation and management
- **Payment Processing**: Multiple payment methods support
- **Receipt Generation**: Automatic receipt creation
- **Enhanced API**: New endpoints for POS operations
- **Business Intelligence**: Built-in analytics capabilities

### Backward Compatibility
- All existing inventory management features remain unchanged
- Existing API endpoints continue to work as before
- Database schema includes all original tables plus new POS tables
- No data migration required

### Migration Path
1. **Existing Users**: Your current inventory data is preserved
2. **New Features**: Start using customer and order management
3. **Gradual Adoption**: Use POS features as needed
4. **Full Integration**: Eventually integrate inventory with sales

### Benefits of the Upgrade
- **Complete Business Solution**: From inventory to sales
- **Customer Relationship Management**: Track customer history
- **Sales Analytics**: Understand your business better
- **Professional Receipts**: Generate proper sales receipts
- **Payment Tracking**: Monitor cash flow and payments
- **Audit Trail**: Complete transaction history

---

**Built with ❤️ using Go, PostgreSQL, and Fiber**

*Upgraded from Inventory Management to Full POS System* 