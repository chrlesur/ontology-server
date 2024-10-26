package api

import (
	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/search"
	"github.com/chrlesur/ontology-server/internal/storage"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, storage *storage.MemoryStorage, logger *logger.Logger) {
	searchEngine := search.NewSearchEngine(storage, logger)
	handler := NewHandler(storage, logger, searchEngine)

	router.GET("/ontologies", handler.ListOntologies)
	router.POST("/ontologies", handler.AddOntology)
	router.GET("/ontologies/:id", handler.GetOntology)
	router.PUT("/ontologies/:id", handler.UpdateOntology)
	router.DELETE("/ontologies/:id", handler.DeleteOntology)
	router.POST("/ontologies/load", handler.LoadOntology)

	router.GET("/search", handler.SearchOntologies)

	router.GET("/elements/details/:element_id", handler.ElementDetailsHandler)
	router.GET("/elements/relations/:element_name", handler.GetElementRelations)
}
