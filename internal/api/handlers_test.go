package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/chrlesur/ontology-server/internal/search"
	"github.com/chrlesur/ontology-server/internal/storage"
	"github.com/gin-gonic/gin"
)

func setupTestHandler() (*Handler, *gin.Engine) {
	storage := storage.NewMemoryStorage()
	logger, _ := logger.NewLogger(logger.INFO, "test_logs")
	searchEngine := search.NewSearchEngine(storage, logger)
	handler := NewHandler(storage, logger, searchEngine)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	return handler, router
}

func TestGetOntology(t *testing.T) {
	h, router := setupTestHandler()

	// Add a test ontology
	testOntology := &models.Ontology{ID: "test1", Name: "Test Ontology"}
	h.Storage.AddOntology(testOntology)

	// Test successful retrieval
	router.GET("/ontologies/:id", h.GetOntology)
	req, _ := http.NewRequest("GET", "/ontologies/test1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var returnedOntology models.Ontology
	err := json.Unmarshal(w.Body.Bytes(), &returnedOntology)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if returnedOntology.ID != testOntology.ID {
		t.Errorf("Expected ontology ID %s, got %s", testOntology.ID, returnedOntology.ID)
	}

	// Test non-existent ontology
	req, _ = http.NewRequest("GET", "/ontologies/nonexistent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAddOntology(t *testing.T) {
	h, router := setupTestHandler()

	router.POST("/ontologies", h.AddOntology)

	newOntology := models.Ontology{Name: "New Test Ontology"}
	body, _ := json.Marshal(newOntology)
	req, _ := http.NewRequest("POST", "/ontologies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var returnedOntology models.Ontology
	err := json.Unmarshal(w.Body.Bytes(), &returnedOntology)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if returnedOntology.Name != newOntology.Name {
		t.Errorf("Expected ontology name %s, got %s", newOntology.Name, returnedOntology.Name)
	}
}

func TestSearchOntologies(t *testing.T) {
	h, router := setupTestHandler()

	router.GET("/search", h.SearchOntologies)

	req, _ := http.NewRequest("GET", "/search?q=Test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var results []search.SearchResult
	err := json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Add more specific assertions based on your expected search results
}

func TestElementDetailsHandler(t *testing.T) {
	h, router := setupTestHandler()

	// Add a test ontology with an element
	testElement := &models.OntologyElement{
		Name:        "Test Element",
		Type:        "TestType",
		Description: "This is a test element",
	}
	testOntology := &models.Ontology{
		ID:       "test1",
		Name:     "Test Ontology",
		Elements: []*models.OntologyElement{testElement},
	}
	h.Storage.AddOntology(testOntology)

	router.GET("/elements/:element_id", h.ElementDetailsHandler)

	req, _ := http.NewRequest("GET", "/elements/Test Element", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var returnedElement models.OntologyElement
	err := json.Unmarshal(w.Body.Bytes(), &returnedElement)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if returnedElement.Name != "Test Element" {
		t.Errorf("Expected element name 'Test Element', got '%s'", returnedElement.Name)
	}
}

// Ajouter après les tests existants...

func TestViewSourceFile(t *testing.T) {
	h, router := setupTestHandler()

	// Configurer la route pour le test
	router.GET("/view-source", h.ViewSourceFile)

	// Créer quelques fichiers de test temporaires
	tmpDir := t.TempDir()

	// Fichier texte
	txtContent := "Test content"
	txtPath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(txtPath, []byte(txtContent), 0644); err != nil {
		t.Fatalf("Failed to create test text file: %v", err)
	}

	// Fichier Markdown
	mdContent := "# Test Title\nTest content"
	mdPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to create test markdown file: %v", err)
	}

	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedType string
		expectedBody map[string]string // Changé pour correspondre à gin.H
	}{
		{
			name:         "View Text File",
			path:         txtPath,
			expectedCode: http.StatusOK,
			expectedType: "text/plain",
			expectedBody: nil, // pas de body JSON pour les fichiers
		},
		{
			name:         "View Markdown File",
			path:         mdPath,
			expectedCode: http.StatusOK,
			expectedType: "text/html",
			expectedBody: nil, // pas de body JSON pour les fichiers
		},
		{
			name:         "File Not Found",
			path:         filepath.Join(tmpDir, "nonexistent.txt"),
			expectedCode: http.StatusNotFound,
			expectedType: "application/json",
			expectedBody: map[string]string{"error": "File not found"},
		},
		{
			name:         "No Path Provided",
			path:         "",
			expectedCode: http.StatusBadRequest,
			expectedType: "application/json",
			expectedBody: map[string]string{"error": "No file path provided"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Construire la requête
			urlPath := "/view-source"
			if tt.path != "" {
				urlPath = fmt.Sprintf("/view-source?path=%s", url.QueryEscape(tt.path)) // Correction ici
			}

			req, _ := http.NewRequest("GET", urlPath, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Vérifier le code de statut
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
			}

			// Vérifier le type de contenu
			contentType := w.Header().Get("Content-Type")
			if !strings.Contains(contentType, tt.expectedType) {
				t.Errorf("Expected content type %s, got %s", tt.expectedType, contentType)
			}

			// Vérifier le contenu selon le type de réponse
			switch tt.expectedType {
			case "application/json":
				var response map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}
				if response["error"] != tt.expectedBody["error"] {
					t.Errorf("Expected error message %v, got %v", tt.expectedBody["error"], response["error"])
				}
			case "text/html":
				// Pour le markdown, vérifier que le contenu HTML est non vide
				if w.Body.Len() == 0 {
					t.Error("Expected non-empty HTML content")
				}
			default:
				// Pour les autres types, vérifier que le contenu n'est pas vide
				if w.Body.Len() == 0 {
					t.Error("Expected non-empty content")
				}
			}

			// Vérifier l'en-tête Content-Disposition pour les fichiers
			if tt.expectedCode == http.StatusOK {
				contentDisposition := w.Header().Get("Content-Disposition")
				expectedFilename := filepath.Base(tt.path)
				expected := fmt.Sprintf("inline; filename=%s", expectedFilename)
				if contentDisposition != expected {
					t.Errorf("Expected Content-Disposition %s, got %s", expected, contentDisposition)
				}
			}
		})
	}
}
