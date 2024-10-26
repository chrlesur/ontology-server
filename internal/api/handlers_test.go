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
	"github.com/gorilla/mux"
)

func setupTestHandler() *Handler {
	storage := storage.NewMemoryStorage()
	logger, _ := logger.NewLogger(logger.INFO, "test_logs")
	searchEngine := search.NewSearchEngine(storage)
	return NewHandler(storage, logger, searchEngine)
}

func TestGetOntology(t *testing.T) {
	h := setupTestHandler()

	// Add a test ontology
	testOntology := &models.Ontology{ID: "test1", Name: "Test Ontology"}
	h.Storage.AddOntology(testOntology)

	// Test successful retrieval
	req, _ := http.NewRequest("GET", "/ontologies/test1", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/ontologies/{id}", h.GetOntology)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedOntology models.Ontology
	json.Unmarshal(rr.Body.Bytes(), &returnedOntology)
	if returnedOntology.ID != testOntology.ID {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedOntology.ID, testOntology.ID)
	}

	// Test non-existent ontology
	req, _ = http.NewRequest("GET", "/ontologies/nonexistent", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAddOntology(t *testing.T) {
	h := setupTestHandler()

	newOntology := models.Ontology{Name: "New Test Ontology"}
	body, _ := json.Marshal(newOntology)
	req, _ := http.NewRequest("POST", "/ontologies", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/ontologies", h.AddOntology).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var returnedOntology models.Ontology
	json.Unmarshal(rr.Body.Bytes(), &returnedOntology)
	if returnedOntology.Name != newOntology.Name {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedOntology.Name, newOntology.Name)
	}
}

func TestSearchOntologies(t *testing.T) {
	h := setupTestHandler()

	// Add test ontologies with elements
	h.Storage.AddOntology(&models.Ontology{
		ID:   "test1",
		Name: "Test Ontology 1",
		Elements: []*models.OntologyElement{
			{Name: "Test Element 1", Type: "Type1", Description: "Description 1"},
		},
	})
	h.Storage.AddOntology(&models.Ontology{
		ID:   "test2",
		Name: "Test Ontology 2",
		Elements: []*models.OntologyElement{
			{Name: "Test Element 2", Type: "Type2", Description: "Description 2"},
		},
	})

	req, _ := http.NewRequest("GET", "/search?q=Test", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/search", h.SearchOntologies).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []search.SearchResult
	json.Unmarshal(rr.Body.Bytes(), &results)
	if len(results) != 2 {
		t.Errorf("handler returned unexpected number of results: got %v want %v", len(results), 2)
	}
}

func TestElementDetailsHandler(t *testing.T) {
	h := setupTestHandler()

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

	req, _ := http.NewRequest("GET", "/elements/Test Element", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/elements/{element_id}", h.ElementDetailsHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedElement models.OntologyElement
	json.Unmarshal(rr.Body.Bytes(), &returnedElement)
	if returnedElement.Name != "Test Element" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedElement.Name, "Test Element")
	}
}
