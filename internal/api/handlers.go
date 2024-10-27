package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/chrlesur/ontology-server/internal/search"
	"github.com/chrlesur/ontology-server/internal/storage"
	"github.com/gin-gonic/gin"
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
func (h *Handler) GetOntology(c *gin.Context) {
	id := c.Param("id")

	h.Logger.Info(fmt.Sprintf("Getting ontology with ID: %s", id))

	ontology, err := h.Storage.GetOntology(id)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting ontology: %v", err))
		c.JSON(http.StatusNotFound, gin.H{"error": MsgResourceNotFound})
		return
	}

	c.JSON(http.StatusOK, ontology)
}

// AddOntology ajoute une nouvelle ontologie
func (h *Handler) AddOntology(c *gin.Context) {
	var ontology models.Ontology
	if err := c.ShouldBindJSON(&ontology); err != nil {
		h.Logger.Error(fmt.Sprintf("Error decoding ontology: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": MsgInvalidInput})
		return
	}

	// Générer un ID unique
	ontology.ID = fmt.Sprintf("onto_%d", time.Now().UnixNano())

	h.Logger.Info(fmt.Sprintf("Adding new ontology: %s with ID: %s", ontology.Name, ontology.ID))

	err := h.Storage.AddOntology(&ontology)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error adding ontology: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": MsgInternalServerError})
		return
	}

	c.JSON(http.StatusCreated, ontology)
}

// UpdateOntology met à jour une ontologie existante
func (h *Handler) UpdateOntology(c *gin.Context) {
	id := c.Param("id")

	var ontology models.Ontology
	if err := c.ShouldBindJSON(&ontology); err != nil {
		h.Logger.Error(fmt.Sprintf("Error decoding ontology: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": MsgInvalidInput})
		return
	}

	ontology.ID = id
	h.Logger.Info(fmt.Sprintf("Updating ontology: %s", id))

	err := h.Storage.UpdateOntology(&ontology)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error updating ontology: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": MsgInternalServerError})
		return
	}

	c.JSON(http.StatusOK, ontology)
}

// DeleteOntology supprime une ontologie par son ID
func (h *Handler) DeleteOntology(c *gin.Context) {
	id := c.Param("id")

	h.Logger.Info(fmt.Sprintf("Deleting ontology: %s", id))

	err := h.Storage.DeleteOntology(id)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error deleting ontology: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": MsgInternalServerError})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListOntologies récupère la liste de toutes les ontologies
func (h *Handler) ListOntologies(c *gin.Context) {
	h.Logger.Info("Listing all ontologies")

	ontologies := h.Storage.ListOntologies()

	// Créer une structure pour l'affichage
	type OntologyInfo struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ElementCount int    `json:"elementCount"`
		ContextCount int    `json:"contextCount"`
	}

	var ontologyInfos []OntologyInfo

	for _, onto := range ontologies {
		contextCount := 0
		for _, elem := range onto.Elements {
			contextCount += len(elem.Contexts)
		}

		ontologyInfos = append(ontologyInfos, OntologyInfo{
			ID:           onto.ID,
			Name:         onto.Name,
			ElementCount: len(onto.Elements),
			ContextCount: contextCount,
		})
	}

	// Assurez-vous de toujours renvoyer un tableau, même s'il est vide
	if ontologyInfos == nil {
		ontologyInfos = []OntologyInfo{}
	}

	c.JSON(http.StatusOK, ontologyInfos)
}

// SearchOntologies effectue une recherche dans les ontologies
func (h *Handler) SearchOntologies(c *gin.Context) {
	query := c.Query("q")
	ontologyID := c.Query("ontology_id")
	elementType := c.Query("type")
	contextSize := 5 // Valeur par défaut, vous pouvez la rendre configurable si nécessaire

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	h.Logger.Info(fmt.Sprintf("Searching ontologies with query: %s", query))

	results, err := h.Search.Search(query, ontologyID, elementType, contextSize)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error during search: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred during the search"})
		return
	}

	// Log des résultats côté serveur
	h.Logger.Info(fmt.Sprintf("Search results: %+v", results))

	// Assurez-vous que chaque résultat inclut les contextes
	for i, result := range results {
		element, err := h.Storage.GetElement(result.ElementName)
		if err == nil && element != nil {
			results[i].Contexts = element.Contexts
		}
	}

	c.JSON(http.StatusOK, results)
}

// ElementDetailsHandler récupère les détails d'un élément spécifique
func (h *Handler) ElementDetailsHandler(c *gin.Context) {
	elementName := c.Param("element_id")

	h.Logger.Info(fmt.Sprintf("Getting details for element: %s", elementName))

	element, err := h.Storage.GetElement(elementName)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting element details: %v", err))
		c.JSON(http.StatusNotFound, gin.H{"error": MsgResourceNotFound})
		return
	}

	// Si les contextes sont nuls, essayez de les récupérer séparément
	if element.Contexts == nil {
		contexts, err := h.Storage.GetElementContexts(elementName)
		if err == nil {
			element.Contexts = contexts
		}
	}
	c.JSON(http.StatusOK, element)
}

func (h *Handler) LoadOntology(c *gin.Context) {
	ontologyFile, err := c.FormFile("ontologyFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No ontology file uploaded"})
		return
	}

	// Sauvegarder le fichier d'ontologie temporairement
	ontologyTempFile := filepath.Join(os.TempDir(), ontologyFile.Filename)
	if err := c.SaveUploadedFile(ontologyFile, ontologyTempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ontology file"})
		return
	}
	defer os.Remove(ontologyTempFile)

	var contextTempFile string
	contextFile, err := c.FormFile("contextFile")
	if err == nil {
		// Un fichier de contexte a été fourni
		contextTempFile = filepath.Join(os.TempDir(), contextFile.Filename)
		if err := c.SaveUploadedFile(contextFile, contextTempFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save context file"})
			return
		}
		defer os.Remove(contextTempFile)
	}

	// Charger l'ontologie
	err = h.Storage.LoadOntologyFromFile(ontologyTempFile, contextTempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load ontology: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ontology loaded successfully"})
}

// GetElementRelations récupère les relations d'un élément spécifique
func (h *Handler) GetElementRelations(c *gin.Context) {
	h.Logger.Info("GetElementRelations endpoint called")
	encodedElementName := c.Param("element_name")
	elementName, err := url.QueryUnescape(encodedElementName)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error decoding element name: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element name"})
		return
	}

	h.Logger.Info(fmt.Sprintf("Getting relations for element: %s", elementName))

	relations, err := h.Storage.GetElementRelations(elementName)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting element relations: %v", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Relations not found"})
		return
	}

	h.Logger.Info(fmt.Sprintf("Found %d relations for element: %s", len(relations), elementName))
	c.JSON(http.StatusOK, relations)
}

func (h *Handler) GetElementContexts(c *gin.Context) {
	elementName := c.Param("element_name")
	contexts, err := h.Storage.GetElementContexts(elementName)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting contexts for element %s: %v", elementName, err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
		return
	}
	c.JSON(http.StatusOK, contexts)
}
