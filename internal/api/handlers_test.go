package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
