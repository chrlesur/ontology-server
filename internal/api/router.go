package api

import (
	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/search"
	"github.com/chrlesur/ontology-server/internal/storage"
	"github.com/gorilla/mux"
)

func NewRouter(storage *storage.MemoryStorage, logger *logger.Logger) *mux.Router {
	searchEngine := search.NewSearchEngine(storage)
	handler := NewHandler(storage, logger, searchEngine)
	router := mux.NewRouter()

	router.HandleFunc("/ontologies", handler.ListOntologies).Methods("GET")
	router.HandleFunc("/ontologies", handler.AddOntology).Methods("POST")
	router.HandleFunc("/ontologies/{id}", handler.GetOntology).Methods("GET")
	router.HandleFunc("/ontologies/{id}", handler.UpdateOntology).Methods("PUT")
	router.HandleFunc("/ontologies/{id}", handler.DeleteOntology).Methods("DELETE")
	router.HandleFunc("/search", handler.SearchOntologies).Methods("GET")
	router.HandleFunc("/elements/{element_id}", handler.ElementDetailsHandler).Methods("GET")

	return router
}
