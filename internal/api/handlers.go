package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	fileID := c.Query("file_id") // Ajoutez cette ligne
	ontologyID := c.Query("ontology_id")
	elementType := c.Query("type")
	contextSize := 5 // Valeur par défaut, vous pouvez la rendre configurable si nécessaire

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	h.Logger.Info(fmt.Sprintf("Searching ontologies with query: %s, fileID: %s", query, fileID))

	results, err := h.Search.Search(query, ontologyID, elementType, contextSize, fileID)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error during search: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred during the search"})
		return
	}

	// Log des résultats côté serveur
	h.Logger.Info(fmt.Sprintf("Search results: %+v", results))

	// Enrichir les résultats avec les métadonnées
	enrichedResults := make([]gin.H, len(results))
	for i, result := range results {
		element, err := h.Storage.GetElement(result.ElementName)
		if err == nil && element != nil {
			ontology, _ := h.Storage.GetOntology(result.OntologyID)
			var sourceFile string
			var resultFileID string
			var sourceMetadata *models.SourceMetadata
			if ontology != nil && ontology.Source != nil {
				sourceMetadata = ontology.Source
				// Utiliser le fileID de la requête s'il est fourni, sinon chercher dans les contextes
				if fileID != "" {
					if fileInfo, exists := sourceMetadata.Files[fileID]; exists {
						resultFileID = fileID
						sourceFile = fileInfo.SourceFile
					}
				} else {
					// Logique existante pour trouver le FileID
					for _, context := range element.Contexts {
						if fileInfo, exists := sourceMetadata.Files[context.FileID]; exists {
							resultFileID = context.FileID
							sourceFile = fileInfo.SourceFile
							break
						}
					}
				}
				h.Logger.Info(fmt.Sprintf("File info for %s: ID=%s, SourceFile=%s", result.ElementName, resultFileID, sourceFile))
			}
			enrichedResults[i] = gin.H{
				"ElementName": result.ElementName,
				"ElementType": result.ElementType,
				"Description": result.Description,
				"OntologyID":  result.OntologyID,
				"Contexts":    element.Contexts,
				"FileID":      resultFileID,
				"SourceFile":  sourceFile,
				"SourceMetadata": gin.H{
					"ontology_file":   sourceMetadata.OntologyFile,
					"processing_date": sourceMetadata.ProcessingDate,
					"files":           sourceMetadata.Files,
				},
			}
		}
	}
	c.JSON(http.StatusOK, enrichedResults)
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
	// Fichier d'ontologie principal
	ontologyFile, err := c.FormFile("ontologyFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No ontology file uploaded"})
		return
	}

	// Fichier de métadonnées (obligatoire)
	metadataFile, err := c.FormFile("metadataFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No metadata file uploaded"})
		return
	}

	// Fichier de contexte (optionnel)
	var contextTempFile string
	contextFile, err := c.FormFile("contextFile")

	// Sauvegarder temporairement les fichiers
	ontologyTempFile := filepath.Join(os.TempDir(), ontologyFile.Filename)
	if err := c.SaveUploadedFile(ontologyFile, ontologyTempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ontology file"})
		return
	}
	defer os.Remove(ontologyTempFile)

	metadataTempFile := filepath.Join(os.TempDir(), metadataFile.Filename)
	if err := c.SaveUploadedFile(metadataFile, metadataTempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata file"})
		return
	}
	defer os.Remove(metadataTempFile)

	if contextFile != nil {
		contextTempFile = filepath.Join(os.TempDir(), contextFile.Filename)
		if err := c.SaveUploadedFile(contextFile, contextTempFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save context file"})
			return
		}
		defer os.Remove(contextTempFile)
	}

	// Charger l'ontologie avec les métadonnées
	err = h.Storage.LoadOntologyFromFile(ontologyTempFile, contextTempFile, metadataTempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load ontology: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ontology loaded successfully"})
}

// GetElementRelations récupère les relations d'un élément spécifique
func (h *Handler) GetElementRelations(c *gin.Context) {
	h.Logger.Info("GetElementRelations endpoint called")

	// Décoder le nom de l'élément depuis l'URL
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
		// Retourner un tableau vide avec status 200 si aucune relation n'est trouvée
		c.JSON(http.StatusOK, []models.Relation{})
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

// Ajouter un endpoint pour récupérer les métadonnées d'une ontologie
func (h *Handler) GetOntologyMetadata(c *gin.Context) {
	id := c.Param("id")

	ontology, err := h.Storage.GetOntology(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ontology not found"})
		return
	}

	if ontology.Source == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No metadata available for this ontology"})
		return
	}

	c.JSON(http.StatusOK, ontology.Source)
}

// ViewSourceFile gère l'affichage des fichiers source
func (h *Handler) ViewSourceFile(c *gin.Context) {
	// Récupérer le chemin du fichier depuis les query params
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file path provided"}) // gin.H pour format JSON cohérent
		return
	}

	// Vérifier que le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"}) // gin.H pour format JSON cohérent
		return
	}

	// Déterminer le type MIME
	ext := strings.ToLower(filepath.Ext(filePath))
	var contentType string
	switch ext {
	case ".pdf":
		contentType = "application/pdf"
	case ".md":
		contentType = "text/markdown"
	case ".txt":
		contentType = "text/plain"
	case ".html":
		contentType = "text/html"
	default:
		contentType = "application/octet-stream"
	}

	// Pour les fichiers markdown, convertir en HTML si nécessaire
	if ext == ".md" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		// Si vous voulez ajouter une conversion Markdown vers HTML ici
		// Vous pouvez utiliser une bibliothèque comme blackfriday

		c.Header("Content-Type", "text/html")
		c.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
		c.String(http.StatusOK, `
            <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <style>
                    body { 
                        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
                        line-height: 1.6;
                        max-width: 800px;
                        margin: 0 auto;
                        padding: 20px;
                    }
                    pre {
                        background: #f5f5f5;
                        padding: 15px;
                        border-radius: 5px;
                    }
                </style>
            </head>
            <body>
                <pre>%s</pre>
            </body>
            </html>
        `, string(file))
		return
	}

	// Pour tous les autres types de fichiers
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	c.File(filePath)
}

// GetOntologyFiles récupère la liste des fichiers de toutes les ontologies
func (h *Handler) GetOntologyFiles(c *gin.Context) {
	h.Logger.Info("Getting list of ontology files")

	ontologies := h.Storage.ListOntologies()
	fileList := make(map[string]string)

	for _, onto := range ontologies {
		if onto.Source != nil {
			for fileID, fileInfo := range onto.Source.Files {
				fileList[fileID] = fileInfo.SourceFile
			}
		}
	}

	c.JSON(http.StatusOK, fileList)
}
