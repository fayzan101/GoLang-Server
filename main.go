/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/internal"
	"myapp/internal/inventory"
	"myapp/internal/orders"
	"myapp/internal/products"
	"myapp/internal/reports"
	"myapp/internal/suppliers"
	"myapp/internal/warehouses"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	internal.InitDB(connStr)

	// Product Management APIs (6)
	http.HandleFunc("/products", handleProducts)
	http.HandleFunc("/products/", handleProductsWithID)
	http.HandleFunc("/products/search", products.SearchProducts)

	// Warehouse Management APIs (3)
	http.HandleFunc("/warehouses", handleWarehouses)
	http.HandleFunc("/warehouses/", handleWarehousesWithID)

	// Inventory APIs (5)
	http.HandleFunc("/inventory", handleInventory)
	http.HandleFunc("/inventory/", handleInventoryWithID)
	http.HandleFunc("/inventory/adjust", inventory.AdjustInventory)
	http.HandleFunc("/inventory/low-stock", inventory.GetLowStock)
	http.HandleFunc("/inventory/movements", inventory.GetStockMovements)

	// Supplier APIs (3)
	http.HandleFunc("/suppliers", handleSuppliers)
	http.HandleFunc("/suppliers/", handleSuppliersWithID)

	// Purchase Order APIs (3)
	http.HandleFunc("/purchase-orders", handlePurchaseOrders)
	http.HandleFunc("/purchase-orders/", handlePurchaseOrdersWithID)

	// Sales Order APIs (3)
	http.HandleFunc("/orders", handleOrders)
	http.HandleFunc("/orders/", handleOrdersWithID)

	// Reports & Audit APIs (2)
	http.HandleFunc("/reports/stock-summary", reports.GetStockSummary)
	http.HandleFunc("/audit-logs", reports.GetAuditLogs)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "IMS Backend API v1.0",
		})
	})

	log.Println("🚀 Inventory Management System API started on :3000")
	log.Println("📦 Total APIs: 25")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Route handlers with method routing

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		products.CreateProduct(w, r)
	case http.MethodGet:
		products.ListProducts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleProductsWithID(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/search") {
		products.SearchProducts(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		products.GetProduct(w, r)
	case http.MethodPut:
		products.UpdateProduct(w, r)
	case http.MethodDelete:
		products.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWarehouses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		warehouses.CreateWarehouse(w, r)
	case http.MethodGet:
		warehouses.ListWarehouses(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleWarehousesWithID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		warehouses.GetWarehouse(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		inventory.GetInventory(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleInventoryWithID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		inventory.GetProductInventory(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSuppliers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		suppliers.CreateSupplier(w, r)
	case http.MethodGet:
		suppliers.ListSuppliers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSuppliersWithID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		suppliers.UpdateSupplier(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePurchaseOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		orders.CreatePurchaseOrder(w, r)
	case http.MethodGet:
		orders.ListPurchaseOrders(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePurchaseOrdersWithID(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/receive") && r.Method == http.MethodPut {
		orders.ReceivePurchaseOrder(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		orders.CreateOrder(w, r)
	case http.MethodGet:
		orders.ListOrders(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleOrdersWithID(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/status") && r.Method == http.MethodPut {
		orders.UpdateOrderStatus(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
