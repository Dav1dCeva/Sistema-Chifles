package handlers

import (
	"log"
	"net/http"
	"websocket/hub"

	"github.com/gorilla/websocket"
)

// upgrader convierte una conexión HTTP en WebSocket.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Aquí podrías validar ALLOWED_ORIGIN desde el .env
		return true
	},
}

// ServeWs gestiona nuevas conexiones WebSocket.
func ServeWs(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al convertir a WS:", err)
		return
	}
	client := &hub.Client{
		Hub:  h,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	h.Register <- client

	go client.WritePump()
	client.ReadPump()
}
