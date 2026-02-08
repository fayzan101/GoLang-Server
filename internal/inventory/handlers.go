package inventory

import (
	"encoding/json"
	"myapp/internal"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetInventory - GET /inventory
func GetInventory(w http.ResponseWriter, r *http.Request) {
	var inventory []internal.Inventory
	query := internal.DB.Preload("Product").Preload("Warehouse")

	// Filter by warehouse
	warehouseID := r.URL.Query().Get("warehouse_id")
	if warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if err := query.Find(&inventory).Error; err != nil {
		http.Error(w, "Failed to fetch inventory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   inventory,
	})
}

// GetProductInventory - GET /inventory/{productId}
func GetProductInventory(w http.ResponseWriter, r *http.Request) {
	productID := extractID(r.URL.Path, "/inventory/")
	if productID == 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var inventory []internal.Inventory
	if err := internal.DB.Preload("Warehouse").Where("product_id = ?", productID).Find(&inventory).Error; err != nil {
		http.Error(w, "Failed to fetch inventory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   inventory,
	})
}

// AdjustInventory - POST /inventory/adjust
func AdjustInventory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID   uint   `json:"product_id"`
		WarehouseID uint   `json:"warehouse_id"`
		Quantity    int    `json:"quantity"`
		Reason      string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var inv internal.Inventory
	err := internal.DB.Where("product_id = ? AND warehouse_id = ?", req.ProductID, req.WarehouseID).First(&inv).Error

	if err != nil {
		// Create new inventory record
		inv = internal.Inventory{
			ProductID:   req.ProductID,
			WarehouseID: req.WarehouseID,
			Quantity:    req.Quantity,
		}
		internal.DB.Create(&inv)
	} else {
		// Update existing
		inv.Quantity += req.Quantity
		internal.DB.Save(&inv)
	}

	// Log stock movement
	movement := internal.StockMovement{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		Type:        "ADJUST",
		Quantity:    req.Quantity,
		Reason:      req.Reason,
		CreatedBy:   "system",
		CreatedAt:   time.Now(),
	}
	internal.DB.Create(&movement)

	internal.LogAudit("ADJUST", "Inventory", inv.ID, "system", req.Reason)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   inv,
	})
}

// GetLowStock - GET /inventory/low-stock
func GetLowStock(w http.ResponseWriter, r *http.Request) {
	var inventory []internal.Inventory
	if err := internal.DB.Preload("Product").Preload("Warehouse").
		Where("quantity <= min_stock").Find(&inventory).Error; err != nil {
		http.Error(w, "Failed to fetch low stock items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   inventory,
		"count":  len(inventory),
	})
}

// GetStockMovements - GET /inventory/movements
func GetStockMovements(w http.ResponseWriter, r *http.Request) {
	var movements []internal.StockMovement
	query := internal.DB.Preload("Product").Order("created_at DESC").Limit(100)

	// Filter by product
	productID := r.URL.Query().Get("product_id")
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	// Filter by type
	movementType := r.URL.Query().Get("type")
	if movementType != "" {
		query = query.Where("type = ?", movementType)
	}

	if err := query.Find(&movements).Error; err != nil {
		http.Error(w, "Failed to fetch stock movements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   movements,
	})
}

func extractID(path, prefix string) int {
	idStr := strings.TrimPrefix(path, prefix)
	if idx := strings.Index(idStr, "/"); idx != -1 {
		idStr = idStr[:idx]
	}
	id, _ := strconv.Atoi(idStr)
	return id
}
