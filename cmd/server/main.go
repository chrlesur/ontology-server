package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrlesur/ontology-server/internal/api"
	"github.com/chrlesur/ontology-server/internal/config"
	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Set Gin mode based on config
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Use Gin's logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize storage
	memoryStorage := storage.NewMemoryStorage()

	// Setup API routes
	apiGroup := router.Group("/api")
	api.SetupRoutes(apiGroup, memoryStorage, l)

	// Serve static files
	router.NoRoute(gin.WrapH(http.FileServer(http.Dir("./web"))))

	// Dans votre fonction main ou de configuration du routeur
	router.Use(cors.Default())

	// Prepare HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	l.Info(fmt.Sprintf("Server listening on %s", addr))

	// Start the server
	err = router.Run(addr)
	if err != nil {
		l.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
