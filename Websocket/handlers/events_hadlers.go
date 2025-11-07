package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"websocket/hub"
	"websocket/models"
)

// EventsHandler recibe POST desde el API REST y reenvía los eventos a todos los clientes WS.
func EventsHandler(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar clave secreta compartida (seguridad básica)
		secret := r.Header.Get("X-WS-SECRET")
		if secret != os.Getenv("WS_SECRET") {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		// Decodificar el cuerpo JSON
		var event models.EventMessage
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			return
		}

		// Enviar el evento al canal Broadcast
		h.Broadcast <- event

		log.Printf("📢 Evento recibido y enviado: %s\n", event.Type)
		w.WriteHeader(http.StatusAccepted)
	}
}
