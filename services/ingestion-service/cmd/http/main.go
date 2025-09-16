package main

import (
	"log"
	"os"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/transport"
)

func main() {
	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	// Bootstrap and run HTTP server
	if err := transport.BootstrapHTTPWithDefaults(addr); err != nil {
		log.Fatalf("Failed to bootstrap HTTP server: %v", err)
	}
}
