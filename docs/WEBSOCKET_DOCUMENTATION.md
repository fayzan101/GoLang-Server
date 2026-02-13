# WebSocket Real-Time Inventory Updates

## Overview
This implementation provides real-time inventory updates using WebSockets. When inventory changes occur, all connected clients receive instant notifications about stock adjustments and low stock alerts.

## Features

### 1. Real-Time Inventory Updates
- Instant notifications when inventory is adjusted
- Broadcasts to all connected clients simultaneously
- Includes product ID, warehouse ID, quantity, and action type

### 2. Low Stock Alerts
- Automatic alerts when stock falls below minimum threshold
- Includes product name and current vs. minimum stock levels
- Can trigger browser notifications

### 3. Multi-Client Support
- Supports unlimited concurrent WebSocket connections
- Thread-safe connection management
- Automatic reconnection handling

## WebSocket Endpoint

**URL:** `ws://localhost:3000/ws/inventory`

**Protocol:** WebSocket (ws://)

## Message Types

### 1. Inventory Update
Sent when inventory is created, updated, or adjusted.

```json
{
  "type": "inventory_update",
  "inventory_id": 123,
  "product_id": 45,
  "warehouse_id": 2,
  "quantity": 150,
  "action": "adjusted",
  "timestamp": "2026-02-13T10:30:45Z"
}
```

**Actions:**
- `created` - New inventory record created
- `updated` - Existing inventory updated
- `adjusted` - Inventory adjusted via `/inventory/adjust` endpoint
- `deleted` - Inventory record deleted

### 2. Low Stock Alert
Sent when inventory quantity falls at or below the minimum stock level.

```json
{
  "type": "low_stock_alert",
  "product_id": 45,
  "warehouse_id": 2,
  "current_quantity": 8,
  "min_stock": 10,
  "product_name": "Widget A",
  "timestamp": "2026-02-13T10:30:45Z"
}
```

## Usage

### JavaScript Client Example

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:3000/ws/inventory');

// Connection opened
ws.onopen = () => {
    console.log('Connected to inventory updates');
};

// Listen for messages
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    
    if (data.type === 'inventory_update') {
        console.log(`Inventory updated: Product ${data.product_id}, Quantity: ${data.quantity}`);
    } else if (data.type === 'low_stock_alert') {
        console.log(`Low stock alert: ${data.product_name} - ${data.current_quantity}/${data.min_stock}`);
    }
};

// Handle errors
ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

// Connection closed
ws.onclose = () => {
    console.log('Disconnected from server');
};
```

### Demo HTML Client

A fully functional demo HTML client is available at:
**`web/inventory-realtime.html`**

To use:
1. Start the server: `go run main.go`
2. Open `web/inventory-realtime.html` in your browser
3. Click "Connect" to establish WebSocket connection
4. Make inventory changes via API endpoints to see real-time updates

## Testing the WebSocket

### Step 1: Start the Server
```bash
go mod tidy  # Install dependencies including gorilla/websocket
go run main.go
```

### Step 2: Connect a Client
Open the demo HTML client or create your own WebSocket connection to `ws://localhost:3000/ws/inventory`

### Step 3: Trigger Updates
Use the REST API to make inventory changes:

```bash
# Adjust inventory (will trigger WebSocket broadcast)
curl -X POST http://localhost:3000/inventory/adjust \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "warehouse_id": 1,
    "quantity": -5,
    "reason": "Sale"
  }'
```

### Step 4: Observe Real-Time Updates
All connected WebSocket clients will receive the update instantly!

## Architecture

### Components

1. **Hub (`internal/websocket/hub.go`)**
   - Manages all active WebSocket connections
   - Broadcasts messages to all clients
   - Handles client registration/unregistration
   - Thread-safe operations using mutex

2. **Client (`internal/websocket/client.go`)**
   - Represents individual WebSocket connection
   - Handles read/write operations
   - Implements ping/pong for connection health
   - Automatic cleanup on disconnect

3. **Handler (`internal/websocket/handler.go`)**
   - HTTP to WebSocket upgrade
   - Global hub instance management
   - Connection request handling

4. **Integration (`internal/inventory/handlers.go`)**
   - Broadcasts updates after inventory changes
   - Checks for low stock conditions
   - Integrates seamlessly with existing handlers

## Configuration

### Connection Settings
Defined in `internal/websocket/client.go`:

```go
writeWait      = 10 * time.Second  // Write timeout
pongWait       = 60 * time.Second  // Read timeout
pingPeriod     = 54 * time.Second  // Ping interval
maxMessageSize = 512               // Max message size
```

### CORS Settings
Currently allows all origins (development mode). For production, modify `internal/websocket/handler.go`:

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // Restrict to specific origins in production
        origin := r.Header.Get("Origin")
        return origin == "https://yourdomain.com"
    },
}
```

## Production Considerations

1. **CORS Configuration**
   - Restrict `CheckOrigin` to trusted domains
   - Implement authentication tokens

2. **Rate Limiting**
   - Add rate limiting to prevent abuse
   - Limit connections per IP/user

3. **Message Filtering**
   - Allow clients to subscribe to specific products/warehouses
   - Reduce unnecessary bandwidth

4. **Monitoring**
   - Track active connections count
   - Monitor message throughput
   - Log connection/disconnection events

5. **SSL/TLS**
   - Use `wss://` instead of `ws://` in production
   - Implement proper certificate management

## Troubleshooting

### WebSocket Connection Fails
- Ensure server is running on port 3000
- Check firewall settings
- Verify URL is correct: `ws://localhost:3000/ws/inventory`

### No Updates Received
- Verify WebSocket connection is established
- Check server logs for broadcast messages
- Ensure inventory changes are being made through the API

### Connection Drops
- Check network stability
- Increase `pongWait` timeout if needed
- Implement client-side reconnection logic

## Future Enhancements

- [ ] Client-side filtering (subscribe to specific products)
- [ ] Authentication/Authorization
- [ ] Message history/replay
- [ ] Compression for large broadcasts
- [ ] Horizontal scaling with Redis pub/sub
- [ ] WebSocket connection metrics dashboard

## Dependencies

- **gorilla/websocket** v1.5.3 - High-performance WebSocket library

Install with:
```bash
go get github.com/gorilla/websocket@v1.5.3
go mod tidy
```
