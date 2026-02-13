package warehouses

import (
	"encoding/json"
	"myapp/internal"
	"myapp/internal/websocket"
	"net/http"
	"strconv"
	"strings"
)

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

	// Broadcast warehouse creation via WebSocket
	if hub := websocket.GetHub(); hub != nil {
		hub.BroadcastWarehouseUpdate(warehouse.ID, warehouse.Name, warehouse.Location, warehouse.Capacity, "created")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   warehouse,
	})
}
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

func UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
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

	var updates internal.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update warehouse fields
	if updates.Name != "" {
		warehouse.Name = updates.Name
	}
	if updates.Location != "" {
		warehouse.Location = updates.Location
	}
	if updates.Capacity != 0 {
		warehouse.Capacity = updates.Capacity
	}

	if err := internal.DB.Save(&warehouse).Error; err != nil {
		http.Error(w, "Failed to update warehouse", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("UPDATE", "Warehouse", warehouse.ID, "system", "Updated warehouse")

	// Broadcast warehouse update via WebSocket
	if hub := websocket.GetHub(); hub != nil {
		hub.BroadcastWarehouseUpdate(warehouse.ID, warehouse.Name, warehouse.Location, warehouse.Capacity, "updated")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   warehouse,
	})
}

func DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
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

	// Store warehouse info before deletion for WebSocket broadcast
	warehouseName := warehouse.Name
	warehouseLocation := warehouse.Location
	warehouseCapacity := warehouse.Capacity
	warehouseID := warehouse.ID

	if err := internal.DB.Delete(&warehouse).Error; err != nil {
		http.Error(w, "Failed to delete warehouse", http.StatusInternalServerError)
		return
	}

	internal.LogAudit("DELETE", "Warehouse", warehouseID, "system", "Deleted warehouse")

	// Broadcast warehouse deletion via WebSocket
	if hub := websocket.GetHub(); hub != nil {
		hub.BroadcastWarehouseUpdate(warehouseID, warehouseName, warehouseLocation, warehouseCapacity, "deleted")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Warehouse deleted successfully",
	})
}
