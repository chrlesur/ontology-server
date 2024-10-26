package storage

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/chrlesur/ontology-server/internal/parser"
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

// GetElementRelations récupère toutes les relations d'un élément spécifique
func (ms *MemoryStorage) GetElementRelations(elementName string) ([]*models.Relation, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	var relations []*models.Relation

	for _, ontology := range ms.ontologies {
		for _, relation := range ontology.Relations {
			if relation.Source == elementName || relation.Target == elementName {
				relations = append(relations, relation)
			}
		}
	}

	if len(relations) == 0 {
		return nil, fmt.Errorf("no relations found for element: %s", elementName)
	}

	return relations, nil
}

func (ms *MemoryStorage) LoadOntologyFromFile(ontologyFile, contextFile string) error {
	log.Info(fmt.Sprintf("Loading ontology from file: %s", ontologyFile))

	var elements []*models.OntologyElement
	var relations []*models.Relation
	var err error

	// Charger l'ontologie principale
	switch {
	case strings.HasSuffix(ontologyFile, ".tsv"):
		elementsSlice, relationsSlice, err := parser.ParseTSV(ontologyFile)
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing TSV file: %v", err))
			return fmt.Errorf("error parsing TSV file: %w", err)
		}
		elements = make([]*models.OntologyElement, len(elementsSlice))
		relations = make([]*models.Relation, len(relationsSlice))
		for i := range elementsSlice {
			elements[i] = &elementsSlice[i]
		}
		for i := range relationsSlice {
			relations[i] = &relationsSlice[i]
		}
	case strings.HasSuffix(ontologyFile, ".owl"):
		elementsSlice, relationsSlice, err := parser.ParseOWL(ontologyFile)
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing OWL file: %v", err))
			return fmt.Errorf("error parsing OWL file: %w", err)
		}
		elements = make([]*models.OntologyElement, len(elementsSlice))
		relations = make([]*models.Relation, len(relationsSlice))
		for i := range elementsSlice {
			elements[i] = &elementsSlice[i]
		}
		for i := range relationsSlice {
			relations[i] = &relationsSlice[i]
		}
	case strings.HasSuffix(ontologyFile, ".rdf"):
		elementsSlice, relationsSlice, err := parser.ParseRDF(ontologyFile)
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing RDF file: %v", err))
			return fmt.Errorf("error parsing RDF file: %w", err)
		}
		elements = make([]*models.OntologyElement, len(elementsSlice))
		relations = make([]*models.Relation, len(relationsSlice))
		for i := range elementsSlice {
			elements[i] = &elementsSlice[i]
		}
		for i := range relationsSlice {
			relations[i] = &relationsSlice[i]
		}
	default:
		return fmt.Errorf("unsupported ontology file format")
	}

	// Charger le fichier de contexte JSON si fourni
	var contexts []models.JSONContext
	if contextFile != "" {
		contexts, err = parser.ParseJSON(contextFile)
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing context file: %v", err))
			return fmt.Errorf("error parsing context file: %w", err)
		}
	}

	// Associer les contextes aux éléments
	elementMap := make(map[string]*models.OntologyElement)
	for i := range elements {
		elementMap[elements[i].Name] = elements[i]
	}

	for _, ctx := range contexts {
		if elem, ok := elementMap[ctx.Element]; ok {
			if elem.Contexts == nil {
				elem.Contexts = make([]models.JSONContext, 0)
			}
			elem.Contexts = append(elem.Contexts, ctx)
		}
	}

	// Créer une nouvelle ontologie avec les éléments et relations parsés
	ontology := &models.Ontology{
		ID:         fmt.Sprintf("onto_%d", time.Now().UnixNano()),
		Name:       filepath.Base(ontologyFile),
		Filename:   ontologyFile,
		Format:     filepath.Ext(ontologyFile)[1:], // Remove the dot from the extension
		Size:       0,                              // You might want to get the actual file size
		SHA256:     "",                             // You might want to calculate the file hash
		ImportedAt: time.Now(),
		Elements:   elements,
		Relations:  relations,
	}

	log.Info(fmt.Sprintf("Loaded ontology with ID: %s", ontology.ID))
	log.Info(fmt.Sprintf("Number of elements: %d", len(ontology.Elements)))
	log.Info(fmt.Sprintf("Number of relations: %d", len(ontology.Relations)))

	// Ajouter l'ontologie au stockage
	err = ms.AddOntology(ontology)
	if err != nil {
		log.Error(fmt.Sprintf("Error adding ontology to storage: %v", err))
		return fmt.Errorf("error adding ontology to storage: %w", err)
	}

	log.Info("Ontology successfully loaded and added to storage")
	return nil
}

func (ms *MemoryStorage) GetElementContexts(elementName string) ([]models.JSONContext, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	log.Info("GetElementContext Called")

	for _, ontology := range ms.ontologies {
		for _, elem := range ontology.Elements {
			if elem.Name == elementName {
				log.Info("GetElementContext found")
				return elem.Contexts, nil
			}
		}
	}

	return nil, fmt.Errorf("element not found")
}
