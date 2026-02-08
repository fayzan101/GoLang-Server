# Inventory Management System (IMS)

A comprehensive REST API for enterprise-grade inventory management built with Go, GORM, and PostgreSQL. This project provides complete inventory tracking, warehouse management, purchase orders, sales orders, and audit logging capabilities.

## Features

- **Product Management** – Create, read, update, delete products with detailed specifications
- **Warehouse Management** – Multi-warehouse support with location tracking
- **Inventory Tracking** – Real-time stock level monitoring and stock movement history
- **Purchase Orders** – Supplier order management with item-level tracking
- **Sales Orders** – Customer order management with fulfillment tracking
- **Audit Logging** – Comprehensive audit trail for compliance and accountability
- **Reporting** – Generate inventory and sales reports
- **RESTful API** – Clean, standard REST endpoints for all operations
- **Database Migrations** – Automated schema management with SQL migrations
- **Environment Configuration** – Flexible configuration via .env file

## Prerequisites

- Go 1.25+ ([download](https://golang.org/dl/))
- PostgreSQL 12+ ([download](https://www.postgresql.org/download/))
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Golang
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   createdb ims_db
   createuser ims_user -P  # You'll be prompted for a password
   ```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=ims_user
DB_PASSWORD=your_password_here
DB_NAME=ims_db
SERVER_PORT=8080
```

## Running the Application

Start the server:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## Project Structure

```
.
├── main.go                  # Application entry point
├── go.mod                   # Go module definition
├── api/                     # API-related code
├── cmd/                     # Command-line interface (Cobra CLI)
├── configs/                 # Configuration management
├── internal/                # Private application code
│   ├── db.go               # Database initialization
│   ├── models.go           # Data models
│   ├── audit/              # Audit logging module
│   ├── inventory/          # Inventory management module
│   ├── orders/             # Sales orders module
│   ├── products/           # Product management module
│   ├── reports/            # Reporting module
│   ├── suppliers/          # Supplier management module
│   └── warehouses/         # Warehouse management module
├── migrations/              # SQL migration files
│   ├── 001_create_users.sql
│   └── 002_ims_schema.sql
├── pkg/                     # Public/reusable packages
│   ├── utils.go            # Utility functions
│   └── middleware/          # HTTP middleware
├── test/                    # Test files
├── web/                     # Web assets
└── docs/                    # Documentation
```

## API Endpoints

### Products
- `GET /products` – List all products
- `GET /products/{id}` – Get product details
- `POST /products` – Create a new product
- `PUT /products/{id}` – Update a product
- `DELETE /products/{id}` – Delete a product
- `GET /products/search` – Search products

### Warehouses
- `GET /warehouses` – List all warehouses
- `GET /warehouses/{id}` – Get warehouse details
- `POST /warehouses` – Create a warehouse
- `PUT /warehouses/{id}` – Update a warehouse
- `DELETE /warehouses/{id}` – Delete a warehouse

### Inventory
- `GET /inventory` – List inventory items
- `GET /inventory/{id}` – Get inventory item details
- `POST /inventory` – Add inventory
- `PUT /inventory/{id}` – Update inventory levels
- `DELETE /inventory/{id}` – Remove inventory

### Purchase Orders
- `GET /purchase-orders` – List purchase orders
- `POST /purchase-orders` – Create purchase order
- `GET /purchase-orders/{id}` – Get order details
- `PUT /purchase-orders/{id}` – Update order
- `DELETE /purchase-orders/{id}` – Cancel order

### Sales Orders
- `GET /orders` – List sales orders
- `POST /orders` – Create order
- `GET /orders/{id}` – Get order details
- `PUT /orders/{id}` – Update order
- `DELETE /orders/{id}` – Cancel order

### Audit Logs
- `GET /audit-logs` – List audit logs

For complete API documentation, see [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md) or import [IMS_Postman_Collection.json](IMS_Postman_Collection.json) into Postman.

## Database Schema

The application uses automated migrations via GORM. Database migrations are located in the `migrations/` directory:
- `001_create_users.sql` – User table and initial setup
- `002_ims_schema.sql` – Complete IMS schema

Migrations run automatically on application startup.

## Development

### Running Tests
```bash
go test ./...
```

Or using the provided test script:
```bash
bash scripts/test.sh
```

### Code Structure Best Practices

This project follows Go conventions:
- **Package organization** – Code is organized by feature domain in `internal/`
- **Exports** – Only exported types and functions (capitalized) are part of the public API
- **Interfaces** – Used for abstraction and testability
- **Error handling** – Explicit error handling throughout

### Adding New Features

1. Create a new package under `internal/` for your feature
2. Define models in the appropriate package or in `internal/models.go`
3. Create handlers in `{feature}/handlers.go`
4. Register routes in `main.go`
5. Add database migrations if needed
6. Update API documentation

## Logging

The application includes comprehensive logging:
- Database connection status
- Migration results
- Audit trail of all data modifications
- Error logging for debugging

Logs use Go's standard `log` package and are printed to stdout.

## Dependencies

Key dependencies:
- **GORM** – Go ORM for database operations
- **PostgreSQL Driver** – GORM PostgreSQL driver
- **Cobra** – CLI framework for command-line tools
- **godotenv** – Environment variable management

See `go.mod` for complete list of dependencies.

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

---

For more information, see [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md) or contact the development team.
