# 📦 Inventory Management System (IMS) - Backend API

**Professional-grade Golang backend with 25 REST APIs**

## 🎯 Overview

Complete inventory management system built with Go, PostgreSQL, and GORM. Ready for production deployment.

## ⚡ Features

- ✅ Product catalog management
- ✅ Multi-warehouse inventory tracking
- ✅ Stock movement auditing
- ✅ Supplier management
- ✅ Purchase order processing
- ✅ Sales order management
- ✅ Real-time stock reports
- ✅ Comprehensive audit logs
- ✅ Low-stock alerts

## 🏗️ Architecture

```
Golang/
├── cmd/
│   └── root.go
├── internal/
│   ├── models.go          # All database models
│   ├── db.go              # Database & audit functions
│   ├── products/          # Product handlers
│   ├── warehouses/        # Warehouse handlers
│   ├── inventory/         # Inventory handlers
│   ├── suppliers/         # Supplier handlers
│   ├── orders/            # Order handlers (PO & Sales)
│   └── reports/           # Reports & audit handlers
├── migrations/
│   └── 002_ims_schema.sql
├── main.go                # Server & routes
└── .env
```

## 📡 API Endpoints (25 Total)

### 1️⃣ Product Management (6 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/products` | Create new product |
| GET | `/products` | List all products |
| GET | `/products/{id}` | Get product by ID |
| PUT | `/products/{id}` | Update product |
| DELETE | `/products/{id}` | Delete product |
| GET | `/products/search?q=keyword` | Search products |

**Example Request (Create Product):**
```json
POST /products
{
  "name": "Widget Pro",
  "sku": "WGT-PRO-001",
  "description": "Premium widget",
  "category": "Electronics",
  "price": 299.99,
  "cost": 150.00,
  "unit": "piece"
}
```

---

### 2️⃣ Warehouse Management (3 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/warehouses` | Create warehouse |
| GET | `/warehouses` | List all warehouses |
| GET | `/warehouses/{id}` | Get warehouse by ID |

**Example Request:**
```json
POST /warehouses
{
  "name": "North Warehouse",
  "location": "Seattle, WA",
  "capacity": 15000
}
```

---

### 3️⃣ Inventory / Stock (5 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/inventory` | Get current stock levels |
| GET | `/inventory/{productId}` | Get stock for specific product |
| POST | `/inventory/adjust` | Manual stock adjustment |
| GET | `/inventory/low-stock` | Get low-stock alerts |
| GET | `/inventory/movements` | View stock movement history |

**Example Request (Adjust Stock):**
```json
POST /inventory/adjust
{
  "product_id": 1,
  "warehouse_id": 1,
  "quantity": 100,
  "reason": "Initial stock"
}
```

---

### 4️⃣ Supplier Management (3 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/suppliers` | Create supplier |
| GET | `/suppliers` | List all suppliers |
| PUT | `/suppliers/{id}` | Update supplier |

**Example Request:**
```json
POST /suppliers
{
  "name": "TechParts Inc",
  "contact_name": "John Doe",
  "email": "john@techparts.com",
  "phone": "+1-555-0123",
  "address": "123 Supply St, NY",
  "rating": 4.5
}
```

---

### 5️⃣ Purchase Orders (3 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/purchase-orders` | Create purchase order |
| GET | `/purchase-orders` | List purchase orders |
| PUT | `/purchase-orders/{id}/receive` | Mark PO as received |

**Example Request (Create PO):**
```json
POST /purchase-orders
{
  "supplier_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 500,
      "unit_price": 15.00
    }
  ]
}
```

**Example Request (Receive PO):**
```json
PUT /purchase-orders/1/receive
{
  "warehouse_id": 1
}
```

---

### 6️⃣ Sales Orders (3 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders` | Create sales order |
| GET | `/orders` | List sales orders |
| PUT | `/orders/{id}/status` | Update order status |

**Example Request (Create Order):**
```json
POST /orders
{
  "customer_name": "ABC Corp",
  "customer_email": "orders@abc.com",
  "items": [
    {
      "product_id": 1,
      "warehouse_id": 1,
      "quantity": 10
    }
  ]
}
```

**Example Request (Update Status):**
```json
PUT /orders/1/status
{
  "status": "shipped"
}
```

**Order Statuses:** `pending`, `processing`, `shipped`, `delivered`, `cancelled`

---

### 7️⃣ Reports & Audit (2 APIs)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reports/stock-summary` | Get stock value report |
| GET | `/audit-logs` | View audit trail |

**Stock Summary Response:**
```json
{
  "status": "success",
  "data": {
    "items": [...],
    "total_value": 125000.50,
    "total_items": 2500
  }
}
```

---

## 🗄️ Database Schema

**Tables:**
- `products` - Product catalog
- `warehouses` - Storage locations
- `inventories` - Current stock levels
- `stock_movements` - Stock transaction history
- `suppliers` - Supplier information
- `purchase_orders` - Purchase orders
- `po_items` - PO line items
- `orders` - Sales orders
- `order_items` - Order line items
- `audit_logs` - System audit trail

---

## 🚀 Quick Start

### 1. Setup Environment

Create `.env` file:
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ims_db
DB_USER=postgres
DB_PASSWORD=yourpassword
```

### 2. Initialize Database

```bash
psql -U postgres -f migrations/002_ims_schema.sql
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Server

```bash
go run main.go
```

Server starts on: **http://localhost:8080**

---

## 🧪 Test the API

### Health Check
```bash
curl http://localhost:8080/health
```

### Create Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "sku": "TST-001",
    "price": 99.99,
    "cost": 50.00
  }'
```

### Get Low Stock
```bash
curl http://localhost:8080/inventory/low-stock
```

---

## 📊 Key Business Logic

### Stock Movement Types
- **IN** - Stock added (from purchase orders)
- **OUT** - Stock removed (from sales orders)
- **ADJUST** - Manual adjustments

### Automatic Stock Updates
- ✅ Creating sales order → reduces inventory
- ✅ Receiving purchase order → increases inventory
- ✅ All movements logged in `stock_movements`
- ✅ All actions recorded in `audit_logs`

### Inventory Checks
- ❌ Cannot create order if insufficient stock
- ⚠️ Low-stock alerts when `quantity <= min_stock`

---

## 💼 Production Ready Features

✅ GORM ORM with auto-migrations  
✅ Foreign key constraints  
✅ Indexed queries for performance  
✅ Audit logging for compliance  
✅ Error handling & validation  
✅ RESTful API design  
✅ Transaction history  

---

## 🔧 Tech Stack

- **Language:** Go 1.21+
- **Framework:** net/http (stdlib)
- **Database:** PostgreSQL 14+
- **ORM:** GORM v2
- **Config:** godotenv

---

## 📈 Next Steps (Optional Enhancements)

- [ ] JWT authentication
- [ ] Pagination for list endpoints
- [ ] CSV/Excel export for reports
- [ ] Email notifications (low stock, orders)
- [ ] Redis caching
- [ ] Docker containerization
- [ ] API documentation (Swagger)
- [ ] Unit tests

---

## 📝 License

MIT License - Free to use commercially

---

## 💡 Use Cases

Perfect for:
- **Retail businesses** - Multi-location inventory
- **Distributors** - Supplier & order management
- **Manufacturers** - Raw material tracking
- **E-commerce** - Order fulfillment
- **Small warehouses** - Stock control

---

**Built with ❤️ using Golang**
