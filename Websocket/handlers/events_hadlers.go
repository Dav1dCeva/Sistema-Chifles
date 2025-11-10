package handlers

import (
    "encoding/json"
    "net/http"
    "websocket/hub"
)

// Estructura del evento recibido desde el API REST
type Event struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// Recibe eventos desde la API REST
func HandleNotify(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
    var event Event
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Error al leer JSON", http.StatusBadRequest)
        return
    }

    msg, _ := json.Marshal(event)
    h.Broadcast <- msg

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Evento enviado correctamente"))
}
