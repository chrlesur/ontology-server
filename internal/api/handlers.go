package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/chrlesur/ontology-server/internal/search"
	"github.com/chrlesur/ontology-server/internal/storage"
	"github.com/gorilla/mux"
)

// Handler encapsule les dépendances nécessaires pour gérer les requêtes API
type Handler struct {
	Storage *storage.MemoryStorage
	Logger  *logger.Logger
	Search  *search.SearchEngine
}

// NewHandler crée une nouvelle instance de Handler avec le stockage, le logger et le moteur de recherche fournis
func NewHandler(storage *storage.MemoryStorage, logger *logger.Logger, search *search.SearchEngine) *Handler {
	return &Handler{Storage: storage, Logger: logger, Search: search}
}

// GetOntology récupère une ontologie par son ID
func (h *Handler) GetOntology(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.Logger.Info(fmt.Sprintf("Getting ontology with ID: %s", id))

	ontology, err := h.Storage.GetOntology(id)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting ontology: %v", err))
		WriteJSONError(w, NewAPIError(ErrNotFound, MsgResourceNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ontology)
}

// AddOntology ajoute une nouvelle ontologie
func (h *Handler) AddOntology(w http.ResponseWriter, r *http.Request) {
	var ontology models.Ontology
	err := json.NewDecoder(r.Body).Decode(&ontology)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error decoding ontology: %v", err))
		WriteJSONError(w, NewAPIError(ErrBadRequest, MsgInvalidInput))
		return
	}

	// Générer un ID unique
	ontology.ID = fmt.Sprintf("onto_%d", time.Now().UnixNano())

	h.Logger.Info(fmt.Sprintf("Adding new ontology: %s with ID: %s", ontology.Name, ontology.ID))

	err = h.Storage.AddOntology(&ontology)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error adding ontology: %v", err))
		WriteJSONError(w, NewAPIError(ErrInternalServerError, MsgInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ontology)
}

// UpdateOntology met à jour une ontologie existante
func (h *Handler) UpdateOntology(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var ontology models.Ontology
	err := json.NewDecoder(r.Body).Decode(&ontology)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error decoding ontology: %v", err))
		WriteJSONError(w, NewAPIError(ErrBadRequest, MsgInvalidInput))
		return
	}

	ontology.ID = id
	h.Logger.Info(fmt.Sprintf("Updating ontology: %s", id))

	err = h.Storage.UpdateOntology(&ontology)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error updating ontology: %v", err))
		WriteJSONError(w, NewAPIError(ErrInternalServerError, MsgInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ontology)
}

// DeleteOntology supprime une ontologie par son ID
func (h *Handler) DeleteOntology(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.Logger.Info(fmt.Sprintf("Deleting ontology: %s", id))

	err := h.Storage.DeleteOntology(id)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error deleting ontology: %v", err))
		WriteJSONError(w, NewAPIError(ErrInternalServerError, MsgInternalServerError))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListOntologies récupère la liste de toutes les ontologies
func (h *Handler) ListOntologies(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Listing all ontologies")

	ontologies := h.Storage.ListOntologies()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ontologies)
}

// SearchOntologies effectue une recherche dans les ontologies
func (h *Handler) SearchOntologies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	ontologyID := r.URL.Query().Get("ontology_id")

	if query == "" {
		WriteJSONError(w, NewAPIError(ErrBadRequest, "Query parameter 'q' is required"))
		return
	}

	h.Logger.Info(fmt.Sprintf("Searching ontologies with query: %s", query))

	// Déboguer: afficher toutes les ontologies
	allOntologies := h.Storage.ListOntologies()
	h.Logger.Info(fmt.Sprintf("All ontologies: %+v", allOntologies))

	results := h.Search.Search(query, ontologyID)

	// Déboguer: afficher les résultats
	h.Logger.Info(fmt.Sprintf("Search results: %+v", results))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ElementDetailsHandler récupère les détails d'un élément spécifique
func (h *Handler) ElementDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	elementName := vars["element_id"]

	h.Logger.Info(fmt.Sprintf("Getting details for element: %s", elementName))

	element, err := h.Storage.GetElement(elementName)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting element details: %v", err))
		WriteJSONError(w, NewAPIError(ErrNotFound, MsgResourceNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(element)
}
