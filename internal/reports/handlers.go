package reports

import (
	"encoding/json"
	"myapp/internal"
	"net/http"
)

func GetStockSummary(w http.ResponseWriter, r *http.Request) {
	var results []struct {
		WarehouseName string  `json:"warehouse_name"`
		ProductName   string  `json:"product_name"`
		SKU           string  `json:"sku"`
		Quantity      int     `json:"quantity"`
		Value         float64 `json:"value"`
	}

	query := `
		SELECT 
			w.name as warehouse_name,
			p.name as product_name,
			p.sku,
			i.quantity,
			(i.quantity * p.price) as value
		FROM inventories i
		JOIN products p ON i.product_id = p.id
		JOIN warehouses w ON i.warehouse_id = w.id
		ORDER BY w.name, p.name
	`

	if err := internal.DB.Raw(query).Scan(&results).Error; err != nil {
		http.Error(w, "Failed to generate stock summary", http.StatusInternalServerError)
		return
	}
	var totalValue float64
	var totalItems int
	for _, r := range results {
		totalValue += r.Value
		totalItems += r.Quantity
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"items":       results,
			"total_value": totalValue,
			"total_items": totalItems,
		},
	})
}
func GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	var logs []internal.AuditLog
	query := internal.DB.Order("created_at DESC").Limit(100)
	entity := r.URL.Query().Get("entity")
	if entity != "" {
		query = query.Where("entity = ?", entity)
	}
	action := r.URL.Query().Get("action")
	if action != "" {
		query = query.Where("action = ?", action)
	}

	if err := query.Find(&logs).Error; err != nil {
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   logs,
	})
}
