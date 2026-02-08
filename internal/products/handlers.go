package products

import (
	"encoding/json"
	"myapp/internal"
	"net/http"
	"strconv"
	"strings"
)

// CreateProduct - POST /products
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product internal.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := internal.DB.Create(&product).Error; err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Audit log
	internal.LogAudit("CREATE", "Product", product.ID, "system", "Created new product")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   product,
	})
}

// ListProducts - GET /products
func ListProducts(w http.ResponseWriter, r *http.Request) {
	var products []internal.Product
	query := internal.DB

	// Filter by category if provided
	category := r.URL.Query().Get("category")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&products).Error; err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   products,
	})
}

// GetProduct - GET /products/{id}
func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/products/")
	if id == 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product internal.Product
	if err := internal.DB.First(&product, id).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   product,
	})
}

// UpdateProduct - PUT /products/{id}
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/products/")
	if id == 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product internal.Product
	if err := internal.DB.First(&product, id).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var updates internal.Product
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := internal.DB.Model(&product).Updates(updates).Error; err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("UPDATE", "Product", product.ID, "system", "Updated product")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   product,
	})
}

// DeleteProduct - DELETE /products/{id}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/products/")
	if id == 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := internal.DB.Delete(&internal.Product{}, id).Error; err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("DELETE", "Product", uint(id), "system", "Deleted product")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}

// SearchProducts - GET /products/search?q=keyword
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		http.Error(w, "Search keyword required", http.StatusBadRequest)
		return
	}

	var products []internal.Product
	searchPattern := "%" + keyword + "%"
	if err := internal.DB.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
		searchPattern, searchPattern, searchPattern).Find(&products).Error; err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   products,
	})
}

// Helper function to extract ID from URL
func extractID(path, prefix string) int {
	idStr := strings.TrimPrefix(path, prefix)
	if idx := strings.Index(idStr, "/"); idx != -1 {
		idStr = idStr[:idx]
	}
	id, _ := strconv.Atoi(idStr)
	return id
}
