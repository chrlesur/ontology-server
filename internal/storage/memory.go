package storage

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
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

	// Normaliser et fusionner les éléments
	normalizedElements := make(map[string]*models.OntologyElement)
	for _, elem := range elements {
		normalizedName := normalizeElementName(elem.Name)
		if existingElem, exists := normalizedElements[normalizedName]; exists {
			// Fusionner les éléments
			existingElem.Positions = append(existingElem.Positions, elem.Positions...)
			if len(elem.Description) > len(existingElem.Description) {
				existingElem.Description = elem.Description
			}
			// Fusionner et dédupliquer les types
			combinedTypes := existingElem.Type + "/" + elem.Type
			existingElem.Type = deduplicateTypes(combinedTypes)
		} else {
			// Dédupliquer les types pour les nouveaux éléments aussi
			elem.Type = deduplicateTypes(elem.Type)
			normalizedElements[normalizedName] = elem
			// Conserver le nom original
			elem.OriginalName = elem.Name
			// Utiliser le nom normalisé comme nouveau nom
			elem.Name = normalizedName
		}
	}

	// Convertir la map en slice
	elements = make([]*models.OntologyElement, 0, len(normalizedElements))
	for _, elem := range normalizedElements {
		elements = append(elements, elem)
	}

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

	// Associer les contextes aux éléments
	totalAssociations := 0
	for _, elem := range elements {
		contextMap := make(map[int]models.JSONContext)
		normalizedElemName := normalizeElementName(elem.Name)
		for _, ctx := range contexts {
			if elementInContext(normalizedElemName, ctx) {
				// Vérifier si au moins une position de l'élément est dans la plage du contexte
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

// Fonction helper pour normaliser les noms d'éléments
func normalizeElementName(name string) string {
	// Remplacer les underscores par des espaces, sauf pour certains préfixes
	parts := strings.SplitN(name, "_", 2)
	if len(parts) == 2 && (parts[0] == "est" || parts[0] == "a") {
		return parts[0] + " " + strings.ReplaceAll(parts[1], "_", " ")
	}

	// Remplacer les underscores par des espaces pour les autres cas
	name = strings.ReplaceAll(name, "_", " ")

	// Liste des préfixes qui peuvent être suivis d'une apostrophe en français
	prefixes := []string{
		"l", "d", "j", "m", "t", "s", "c", "n", "qu",
		"jusqu", "lorsqu", "puisqu", "quoiqu", "quelqu",
	}

	// Remplacer les espaces par des apostrophes pour ces préfixes
	for _, prefix := range prefixes {
		pattern := fmt.Sprintf(`\b%s \b`, prefix)
		replacement := fmt.Sprintf("%s'", prefix)
		name = regexp.MustCompile(pattern).ReplaceAllString(name, replacement)
	}

	// Gestion spéciale pour "aujourd hui"
	name = strings.ReplaceAll(name, "aujourd hui", "aujourd'hui")

	// Supprimer les espaces multiples
	name = strings.Join(strings.Fields(name), " ")

	// Ne pas mettre en minuscules pour préserver la casse originale
	return name
}
func deduplicateTypes(types string) string {
	typeSlice := strings.Split(types, "/")
	typeMap := make(map[string]string)
	for _, t := range typeSlice {
		t = strings.TrimSpace(t)
		normalizedType := normalizeType(t)
		if existingType, exists := typeMap[normalizedType]; exists {
			// Garder la version la plus longue du type
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
	// Remplacer les underscores par des espaces
	t = strings.ReplaceAll(t, "_", " ")
	// Supprimer les espaces multiples
	t = strings.Join(strings.Fields(t), " ")
	// Mettre en minuscules
	return strings.ToLower(t)
}

// Fonction helper pour vérifier si un élément est présent dans un contexte
func elementInContext(elem string, ctx models.JSONContext) bool {
	elemLower := strings.ToLower(elem)
	contextText := strings.ToLower(strings.Join(append(ctx.Before, ctx.Element), " ") + " " + strings.Join(ctx.After, " "))

	// Vérification de la correspondance exacte
	if strings.Contains(contextText, elemLower) {
		return true
	}

	// Vérification avec les underscores remplacés par des espaces
	elemWithoutUnderscore := strings.ReplaceAll(elemLower, "_", " ")
	if strings.Contains(contextText, elemWithoutUnderscore) {
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
		}
	}

	// Si plus de la moitié des parties correspondent, considérez-le comme une correspondance
	if float64(matchCount)/float64(len(elemParts)) > 0.5 {
		return true
	}

	// Vérification spéciale pour les éléments contenant "est" ou "a"
	if strings.Contains(elemLower, "est_") || strings.Contains(elemLower, "a_") {
		parts := strings.SplitN(elemLower, "_", 2)
		if len(parts) == 2 && strings.Contains(contextText, parts[1]) {
			return true
		}
	}

	return false
}
func (ms *MemoryStorage) GetElementContexts(elementName string) ([]models.JSONContext, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	log.Info("GetElementContext Called for: " + elementName)
	normalizedName := normalizeElementName(elementName)

	for _, ontology := range ms.ontologies {
		for _, elem := range ontology.Elements {
			if normalizeElementName(elem.Name) == normalizedName {
				log.Info(fmt.Sprintf("GetElementContext found for %s with %d contexts", elem.Name, len(elem.Contexts)))
				return elem.Contexts, nil
			}
		}
	}

	return nil, fmt.Errorf("element not found: %s", elementName)
}
