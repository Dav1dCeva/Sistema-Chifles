package hub

import (
	"encoding/json"
	"log"
	"websocket/models"
)

// Hub mantiene los clientes conectados y distribuye mensajes.
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan models.EventMessage
	Register   chan *Client
	Unregister chan *Client
}

// NewHub inicializa un nuevo Hub.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan models.EventMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run inicia el bucle principal del Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Nuevo cliente conectado")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Println("Cliente desconectado")
			}
		case message := <-h.Broadcast:
			data, err := json.Marshal(message)
			if err != nil {
				log.Println("Error al serializar mensaje:", err)
				continue
			}
			for c := range h.Clients {
				select {
				case c.Send <- data:
				default:
					close(c.Send)
					delete(h.Clients, c)
				}
			}
		}
	}
}
