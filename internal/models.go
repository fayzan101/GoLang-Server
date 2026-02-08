package internal

import (
	"time"
)

// Product represents a product in the catalog
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	SKU         string    `gorm:"uniqueIndex;not null" json:"sku"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `gorm:"not null" json:"price"`
	Cost        float64   `json:"cost"`
	Unit        string    `json:"unit"` // e.g., "piece", "kg", "liter"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Warehouse represents a storage location
type Warehouse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Location  string    `json:"location"`
	Capacity  int       `json:"capacity"`
	ManagerID uint      `json:"manager_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Inventory represents current stock levels
type Inventory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	WarehouseID uint      `gorm:"not null;index" json:"warehouse_id"`
	Quantity    int       `gorm:"not null;default:0" json:"quantity"`
	MinStock    int       `gorm:"default:10" json:"min_stock"`
	MaxStock    int       `gorm:"default:1000" json:"max_stock"`
	UpdatedAt   time.Time `json:"updated_at"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
}

// StockMovement represents stock in/out transactions
type StockMovement struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	WarehouseID uint      `gorm:"not null;index" json:"warehouse_id"`
	Type        string    `gorm:"not null" json:"type"` // "IN", "OUT", "ADJUST"
	Quantity    int       `gorm:"not null" json:"quantity"`
	Reference   string    `json:"reference"` // Order ID, PO ID, etc.
	Reason      string    `json:"reason"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// Supplier represents a supplier
type Supplier struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	ContactName string    `json:"contact_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PurchaseOrder represents orders placed to suppliers
type PurchaseOrder struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	PONumber   string     `gorm:"uniqueIndex;not null" json:"po_number"`
	SupplierID uint       `gorm:"not null;index" json:"supplier_id"`
	Status     string     `gorm:"not null;default:'pending'" json:"status"` // pending, received, cancelled
	TotalCost  float64    `json:"total_cost"`
	OrderDate  time.Time  `json:"order_date"`
	ReceivedAt *time.Time `json:"received_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Supplier   Supplier   `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Items      []POItem   `gorm:"foreignKey:POID" json:"items,omitempty"`
}

// POItem represents items in a purchase order
type POItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	POID      uint    `gorm:"not null;index" json:"po_id"`
	ProductID uint    `gorm:"not null;index" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// Order represents sales orders
type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	OrderNumber   string      `gorm:"uniqueIndex;not null" json:"order_number"`
	CustomerName  string      `gorm:"not null" json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	Status        string      `gorm:"not null;default:'pending'" json:"status"` // pending, processing, shipped, delivered, cancelled
	TotalAmount   float64     `json:"total_amount"`
	OrderDate     time.Time   `json:"order_date"`
	ShippedAt     *time.Time  `json:"shipped_at,omitempty"`
	DeliveredAt   *time.Time  `json:"delivered_at,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Items         []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// OrderItem represents items in a sales order
type OrderItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	OrderID     uint    `gorm:"not null;index" json:"order_id"`
	ProductID   uint    `gorm:"not null;index" json:"product_id"`
	WarehouseID uint    `gorm:"not null" json:"warehouse_id"`
	Quantity    int     `gorm:"not null" json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Product     Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// AuditLog represents audit trail
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Action    string    `gorm:"not null" json:"action"`
	Entity    string    `gorm:"not null" json:"entity"`
	EntityID  uint      `json:"entity_id"`
	UserID    string    `json:"user_id"`
	Details   string    `gorm:"type:text" json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
