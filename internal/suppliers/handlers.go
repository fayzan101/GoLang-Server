package suppliers

import (
	"encoding/json"
	"myapp/internal"
	"net/http"
	"strconv"
	"strings"
)

// CreateSupplier - POST /suppliers
func CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var supplier internal.Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := internal.DB.Create(&supplier).Error; err != nil {
		http.Error(w, "Failed to create supplier", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("CREATE", "Supplier", supplier.ID, "system", "Created new supplier")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   supplier,
	})
}

// ListSuppliers - GET /suppliers
func ListSuppliers(w http.ResponseWriter, r *http.Request) {
	var suppliers []internal.Supplier
	if err := internal.DB.Find(&suppliers).Error; err != nil {
		http.Error(w, "Failed to fetch suppliers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   suppliers,
	})
}

// UpdateSupplier - PUT /suppliers/{id}
func UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/suppliers/")
	if id == 0 {
		http.Error(w, "Invalid supplier ID", http.StatusBadRequest)
		return
	}

	var supplier internal.Supplier
	if err := internal.DB.First(&supplier, id).Error; err != nil {
		http.Error(w, "Supplier not found", http.StatusNotFound)
		return
	}

	var updates internal.Supplier
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := internal.DB.Model(&supplier).Updates(updates).Error; err != nil {
		http.Error(w, "Failed to update supplier", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("UPDATE", "Supplier", supplier.ID, "system", "Updated supplier")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   supplier,
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
