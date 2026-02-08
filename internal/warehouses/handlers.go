package warehouses

import (
	"encoding/json"
	"myapp/internal"
	"net/http"
	"strconv"
	"strings"
)

// CreateWarehouse - POST /warehouses
func CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var warehouse internal.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := internal.DB.Create(&warehouse).Error; err != nil {
		http.Error(w, "Failed to create warehouse", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("CREATE", "Warehouse", warehouse.ID, "system", "Created new warehouse")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   warehouse,
	})
}

// ListWarehouses - GET /warehouses
func ListWarehouses(w http.ResponseWriter, r *http.Request) {
	var warehouses []internal.Warehouse
	if err := internal.DB.Find(&warehouses).Error; err != nil {
		http.Error(w, "Failed to fetch warehouses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   warehouses,
	})
}

// GetWarehouse - GET /warehouses/{id}
func GetWarehouse(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/warehouses/")
	if id == 0 {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	var warehouse internal.Warehouse
	if err := internal.DB.First(&warehouse, id).Error; err != nil {
		http.Error(w, "Warehouse not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   warehouse,
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
