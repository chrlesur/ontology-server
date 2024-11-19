package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
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
	loader     *OntologyLoader
}

// NewMemoryStorage initializes and returns a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	ms := &MemoryStorage{
		ontologies: make(map[string]*models.Ontology),
	}
	ms.loader = NewOntologyLoader(ms, log)
	return ms
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

func (ms *MemoryStorage) GetElementRelations(elementName string) ([]*models.Relation, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	normalizedName := normalizeElementName(elementName)
	log.Info(fmt.Sprintf("Searching relations for normalized element name: %s", normalizedName))

	var relations []*models.Relation
	for _, ontology := range ms.ontologies {
		for _, relation := range ontology.Relations {
			if normalizeElementName(relation.Source) == normalizedName ||
				normalizeElementName(relation.Target) == normalizedName {
				relations = append(relations, relation)
				log.Info(fmt.Sprintf("Found relation: %s -> %s -> %s",
					relation.Source, relation.Type, relation.Target))
			}
		}
	}

	if len(relations) == 0 {
		return nil, fmt.Errorf("no relations found for element: %s", elementName)
	}

	return relations, nil
}

// LoadOntologyFromFile loads an ontology from files including metadata
func (ms *MemoryStorage) LoadOntologyFromFile(ontologyFile, contextFile, metadataFile string) error {
	return ms.loader.LoadFiles(ontologyFile, contextFile, metadataFile)
}

func loadSourceMetadata(filename string) (*models.SourceMetadata, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata models.SourceMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	return &metadata, nil
}

func normalizeElementName(name string) string {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) == 2 && (parts[0] == "est" || parts[0] == "a") {
		return parts[0] + " " + strings.ReplaceAll(parts[1], "_", " ")
	}

	name = strings.ReplaceAll(name, "_", " ")

	prefixes := []string{
		"l", "d", "j", "m", "t", "s", "c", "n", "qu",
		"jusqu", "lorsqu", "puisqu", "quoiqu", "quelqu",
	}

	for _, prefix := range prefixes {
		pattern := fmt.Sprintf(`\b%s \b`, prefix)
		replacement := fmt.Sprintf("%s'", prefix)
		name = regexp.MustCompile(pattern).ReplaceAllString(name, replacement)
	}

	name = strings.ReplaceAll(name, "aujourd hui", "aujourd'hui")
	name = strings.Join(strings.Fields(name), " ")

	return name
}

func deduplicateTypes(types string) string {
	typeSlice := strings.Split(types, "/")
	typeMap := make(map[string]string)
	for _, t := range typeSlice {
		t = strings.TrimSpace(t)
		normalizedType := normalizeType(t)
		if existingType, exists := typeMap[normalizedType]; exists {
			if len(t) > len(existingType) {
				typeMap[normalizedType] = t
			}
		} else {
			typeMap[normalizedType] = t
		}
	}

	var uniqueTypes []string
	for _, t := range typeMap {
		uniqueTypes = append(uniqueTypes, t)
	}

	sort.Strings(uniqueTypes)
	return strings.Join(uniqueTypes, "/")
}

func normalizeType(t string) string {
	t = strings.ReplaceAll(t, "_", " ")
	t = strings.Join(strings.Fields(t), " ")
	return strings.ToLower(t)
}

// Fonction helper pour vérifier si un élément est présent dans un contexte
func elementInContext(elem string, ctx models.JSONContext) bool {
	log.Info(fmt.Sprintf("Checking if element '%s' is in context of '%s'", elem, ctx.Element))
	elemLower := strings.ToLower(elem)
	contextText := strings.ToLower(strings.Join(append(ctx.Before, ctx.Element), " ") + " " + strings.Join(ctx.After, " "))

	log.Info(fmt.Sprintf("Element: %s, Context: %s", elemLower, contextText))

	// Vérification de la correspondance exacte
	if strings.Contains(contextText, elemLower) {
		log.Info(fmt.Sprintf("Exact match found for '%s' in context", elem))
		return true
	}

	// Vérification avec les underscores remplacés par des espaces
	elemWithoutUnderscore := strings.ReplaceAll(elemLower, "_", " ")
	if strings.Contains(contextText, elemWithoutUnderscore) {
		log.Info(fmt.Sprintf("Match found for '%s' without underscores in context", elem))
		return true
	}

	// Vérification des parties individuelles du nom de l'élément
	elemParts := strings.FieldsFunc(elemLower, func(r rune) bool {
		return r == '_' || r == ' '
	})

	matchCount := 0
	for _, part := range elemParts {
		if strings.Contains(contextText, part) {
			matchCount++
			log.Info(fmt.Sprintf("Partial match found for part '%s' of '%s' in context", part, elem))
		}
	}

	// Si plus de la moitié des parties correspondent, considérez-le comme une correspondance
	if float64(matchCount)/float64(len(elemParts)) > 0.5 {
		log.Info(fmt.Sprintf("Sufficient partial matches found for '%s' in context", elem))
		return true
	}

	// Vérification spéciale pour les éléments contenant "est" ou "a"
	if strings.Contains(elemLower, "est_") || strings.Contains(elemLower, "a_") {
		parts := strings.SplitN(elemLower, "_", 2)
		if len(parts) == 2 && strings.Contains(contextText, parts[1]) {
			log.Info(fmt.Sprintf("Special case match found for '%s' in context", elem))
			return true
		}
	}

	log.Info(fmt.Sprintf("No match found for '%s' in context", elem))
	return false
}

func (ms *MemoryStorage) GetElementContexts(elementName string) ([]models.JSONContext, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	log.Info(fmt.Sprintf("GetElementContext Called for: %s", elementName))
	normalizedName := normalizeElementName(elementName)
	log.Info(fmt.Sprintf("Normalized name: %s", normalizedName))

	for _, ontology := range ms.ontologies {
		for _, elem := range ontology.Elements {
			log.Info(fmt.Sprintf("Checking element: %s (normalized: %s)", elem.Name, normalizeElementName(elem.Name)))
			if normalizeElementName(elem.Name) == normalizedName {
				log.Info(fmt.Sprintf("GetElementContext found for %s with %d contexts", elem.Name, len(elem.Contexts)))
				if len(elem.Contexts) == 0 {
					log.Warning(fmt.Sprintf("Element %s found but has no contexts", elem.Name))
				}
				for i, ctx := range elem.Contexts {
					log.Info(fmt.Sprintf("Context %d for %s: Position=%d, Element=%s", i, elem.Name, ctx.Position, ctx.Element))
				}
				return elem.Contexts, nil
			}
		}
	}

	log.Warning(fmt.Sprintf("Element not found: %s", elementName))
	return nil, fmt.Errorf("element not found: %s", elementName)
}
