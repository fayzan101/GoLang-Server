package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var GlobalHub *Hub

func InitHub() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	log.Println("WebSocket hub initialized and running")
}
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(GlobalHub, conn)
	GlobalHub.register <- client
	go client.WritePump()
	go client.ReadPump()
}
func GetHub() *Hub {
	return GlobalHub
}
