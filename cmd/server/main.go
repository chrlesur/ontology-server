package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrlesur/ontology-server/internal/api"
	"github.com/chrlesur/ontology-server/internal/config"
	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logLevel := logger.INFO // You might want to parse this from the config
	l, err := logger.NewLogger(logLevel, cfg.Logging.Directory)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	l.Info("Starting Ontology Server")

	// Initialize storage
	memoryStorage := storage.NewMemoryStorage()

	// Initialize router
	router := api.NewRouter(memoryStorage)

	// Prepare HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	l.Info(fmt.Sprintf("Server listening on %s", addr))

	// Start the server
	err = http.ListenAndServe(addr, router)
	if err != nil {
		l.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
