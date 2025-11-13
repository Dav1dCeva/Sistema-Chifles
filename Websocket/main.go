package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"websocket/handlers"
	"websocket/hub"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar .env
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	h := hub.NewHub()
	go h.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(h, w, r)
	})

	http.HandleFunc("/events", handlers.EventsHandler(h))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "WebSocket server running ✅")
	})

	log.Printf("Servidor WebSocket corriendo en http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error:", err)
	}
}
