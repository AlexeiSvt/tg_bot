package healtcheck

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func StartHealthCheckServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" 
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "eTriathlon Bot is running! ✅")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	log.Printf("Health check server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Health check server error: %v", err)
	}
}