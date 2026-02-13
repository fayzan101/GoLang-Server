package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}
func (h *Hub) BroadcastInventoryUpdate(inventoryID uint, productID uint, warehouseID uint, quantity int, action string) {
	message := map[string]interface{}{
		"type":         "inventory_update",
		"inventory_id": inventoryID,
		"product_id":   productID,
		"warehouse_id": warehouseID,
		"quantity":     quantity,
		"action":       action, // "created", "updated", "deleted", "adjusted"
		"timestamp":    getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling inventory update: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting inventory update: Product %d in Warehouse %d, Quantity: %d", productID, warehouseID, quantity)
}
func (h *Hub) BroadcastLowStockAlert(productID uint, warehouseID uint, currentQuantity int, minStock int, productName string) {
	message := map[string]interface{}{
		"type":             "low_stock_alert",
		"product_id":       productID,
		"warehouse_id":     warehouseID,
		"current_quantity": currentQuantity,
		"min_stock":        minStock,
		"product_name":     productName,
		"timestamp":        getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling low stock alert: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting low stock alert: Product %s (%d), Quantity: %d/%d", productName, productID, currentQuantity, minStock)
}
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// BroadcastWarehouseUpdate sends warehouse updates to all connected clients
func (h *Hub) BroadcastWarehouseUpdate(warehouseID uint, name string, location string, capacity int, action string) {
	message := map[string]interface{}{
		"type":         "warehouse_update",
		"warehouse_id": warehouseID,
		"name":         name,
		"location":     location,
		"capacity":     capacity,
		"action":       action, // "created", "updated", "deleted"
		"timestamp":    getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling warehouse update: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting warehouse update: %s (%d) - %s", name, warehouseID, action)
}

// BroadcastWarehouseCapacityAlert sends alerts when warehouse capacity threshold is reached
func (h *Hub) BroadcastWarehouseCapacityAlert(warehouseID uint, name string, currentStock int, capacity int, utilizationPercent float64) {
	message := map[string]interface{}{
		"type":                "warehouse_capacity_alert",
		"warehouse_id":        warehouseID,
		"warehouse_name":      name,
		"current_stock":       currentStock,
		"capacity":            capacity,
		"utilization_percent": utilizationPercent,
		"timestamp":           getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling warehouse capacity alert: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting warehouse capacity alert: %s (%d), Utilization: %.2f%%", name, warehouseID, utilizationPercent)
}

// BroadcastProductUpdate sends product updates to all connected clients
func (h *Hub) BroadcastProductUpdate(productID uint, name string, sku string, category string, price float64, action string) {
	message := map[string]interface{}{
		"type":       "product_update",
		"product_id": productID,
		"name":       name,
		"sku":        sku,
		"category":   category,
		"price":      price,
		"action":     action, // "created", "updated", "deleted"
		"timestamp":  getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling product update: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting product update: %s (ID: %d) - %s", name, productID, action)
}

// BroadcastProductPriceAlert sends alerts when product price changes significantly
func (h *Hub) BroadcastProductPriceAlert(productID uint, name string, oldPrice float64, newPrice float64, changePercent float64) {
	message := map[string]interface{}{
		"type":           "product_price_alert",
		"product_id":     productID,
		"product_name":   name,
		"old_price":      oldPrice,
		"new_price":      newPrice,
		"change_percent": changePercent,
		"timestamp":      getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling product price alert: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting product price alert: %s (ID: %d), Change: %.2f%%", name, productID, changePercent)
}

// BroadcastSupplierUpdate sends supplier updates to all connected clients
func (h *Hub) BroadcastSupplierUpdate(supplierID uint, name string, email string, phone string, address string, action string) {
	message := map[string]interface{}{
		"type":        "supplier_update",
		"supplier_id": supplierID,
		"name":        name,
		"email":       email,
		"phone":       phone,
		"address":     address,
		"action":      action, // "created", "updated", "deleted"
		"timestamp":   getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling supplier update: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting supplier update: %s (ID: %d) - %s", name, supplierID, action)
}

// BroadcastSupplierStatusAlert sends alerts for supplier status changes
func (h *Hub) BroadcastSupplierStatusAlert(supplierID uint, name string, status string, message string) {
	alertMsg := map[string]interface{}{
		"type":          "supplier_status_alert",
		"supplier_id":   supplierID,
		"supplier_name": name,
		"status":        status, // "active", "inactive", "suspended", "warning"
		"message":       message,
		"timestamp":     getCurrentTimestamp(),
	}

	jsonMessage, err := json.Marshal(alertMsg)
	if err != nil {
		log.Printf("Error marshaling supplier status alert: %v", err)
		return
	}

	h.broadcast <- jsonMessage
	log.Printf("Broadcasting supplier status alert: %s (ID: %d), Status: %s", name, supplierID, status)
}

func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
