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
				log.Info(fmt.Sprintf("Found element %s with %d contexts", elementName, len(element.Contexts)))
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

	log.Info(fmt.Sprintf("Loaded %d elements and %d relations from ontology file", len(elements), len(relations)))

	// Charger le fichier de contexte JSON si fourni
	var contexts []models.JSONContext
	if contextFile != "" {
		contexts, err = parser.ParseJSON(contextFile)
		if err != nil {
			log.Error(fmt.Sprintf("Error parsing context file: %v", err))
			return fmt.Errorf("error parsing context file: %w", err)
		}
		log.Info(fmt.Sprintf("Loaded %d contexts from JSON file", len(contexts)))
	}

	// Fonction helper pour vérifier si un élément est présent dans un contexte
	elementInContext := func(elem string, ctx models.JSONContext) bool {
		elemLower := strings.ToLower(elem)
		for _, word := range ctx.Before {
			if strings.ToLower(word) == elemLower {
				return true
			}
		}
		for _, word := range ctx.After {
			if strings.ToLower(word) == elemLower {
				return true
			}
		}
		return strings.ToLower(ctx.Element) == elemLower
	}

	// Associer les contextes aux éléments
	totalAssociations := 0
	for _, elem := range elements {
		contextMap := make(map[int]models.JSONContext)
		for _, ctx := range contexts {
			if elementInContext(elem.Name, ctx) {
				for _, pos := range elem.Positions {
					if pos >= ctx.StartOffset && pos <= ctx.EndOffset {
						if _, exists := contextMap[ctx.Position]; !exists {
							contextMap[ctx.Position] = ctx
							totalAssociations++
							log.Info(fmt.Sprintf("Associated new context (position %d) to element %s", ctx.Position, elem.Name))
						}
						break
					}
				}
			}
		}
		// Convertir la map en slice pour l'élément
		elem.Contexts = make([]models.JSONContext, 0, len(contextMap))
		for _, ctx := range contextMap {
			elem.Contexts = append(elem.Contexts, ctx)
		}
		log.Info(fmt.Sprintf("Element %s has %d unique contexts after association", elem.Name, len(elem.Contexts)))
	}

	log.Info(fmt.Sprintf("Associated a total of %d unique contexts to elements", totalAssociations))

	// Créer une nouvelle ontologie avec les éléments et relations parsés
	ontology := &models.Ontology{
		ID:         fmt.Sprintf("onto_%d", time.Now().UnixNano()),
		Name:       filepath.Base(ontologyFile),
		Filename:   ontologyFile,
		Format:     filepath.Ext(ontologyFile)[1:],
		Size:       0,
		SHA256:     "",
		ImportedAt: time.Now(),
		Elements:   elements,
		Relations:  relations,
	}

	log.Info(fmt.Sprintf("Created new ontology with ID: %s", ontology.ID))
	log.Info(fmt.Sprintf("Number of elements in ontology: %d", len(ontology.Elements)))
	log.Info(fmt.Sprintf("Number of relations in ontology: %d", len(ontology.Relations)))

	// Vérification finale des contextes
	for _, elem := range ontology.Elements {
		log.Info(fmt.Sprintf("Element %s has %d unique contexts after ontology creation", elem.Name, len(elem.Contexts)))
	}

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
