package storage

import (
	"fmt"
	"sync"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
)

var log *logger.Logger

func init() {
	var err error
	log, err = logger.NewLogger(logger.INFO, "logs")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

// MemoryStorage represents an in-memory storage for ontologies
type MemoryStorage struct {
	ontologies map[string]*models.Ontology
	mutex      sync.RWMutex
}

// NewMemoryStorage initializes and returns a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ontologies: make(map[string]*models.Ontology),
	}
}

// AddOntology adds a new ontology to the storage
func (ms *MemoryStorage) AddOntology(ontology *models.Ontology) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.ontologies[ontology.ID]; exists {
		return fmt.Errorf("ontology with ID %s already exists", ontology.ID)
	}

	ms.ontologies[ontology.ID] = ontology
	log.Info(fmt.Sprintf("Added ontology with ID: %s", ontology.ID))
	return nil
}

// GetOntology retrieves an ontology by its ID
func (ms *MemoryStorage) GetOntology(id string) (*models.Ontology, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	ontology, exists := ms.ontologies[id]
	if !exists {
		return nil, fmt.Errorf("ontology with ID %s not found", id)
	}

	return ontology, nil
}

// UpdateOntology updates an existing ontology
func (ms *MemoryStorage) UpdateOntology(ontology *models.Ontology) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.ontologies[ontology.ID]; !exists {
		return fmt.Errorf("ontology with ID %s not found", ontology.ID)
	}

	ms.ontologies[ontology.ID] = ontology
	log.Info(fmt.Sprintf("Updated ontology with ID: %s", ontology.ID))
	return nil
}

// DeleteOntology removes an ontology by its ID
func (ms *MemoryStorage) DeleteOntology(id string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.ontologies[id]; !exists {
		return fmt.Errorf("ontology with ID %s not found", id)
	}

	delete(ms.ontologies, id)
	log.Info(fmt.Sprintf("Deleted ontology with ID: %s", id))
	return nil
}

// ListOntologies returns a list of all stored ontologies
func (ms *MemoryStorage) ListOntologies() []*models.Ontology {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	ontologies := make([]*models.Ontology, 0, len(ms.ontologies))
	for _, ontology := range ms.ontologies {
		ontologies = append(ontologies, ontology)
	}

	return ontologies
}

func (ms *MemoryStorage) GetElement(elementName string) (*models.OntologyElement, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	for _, ontology := range ms.ontologies {
		for _, element := range ontology.Elements {
			if element.Name == elementName {
				return element, nil
			}
		}
	}
	return nil, fmt.Errorf("element not found")
}
