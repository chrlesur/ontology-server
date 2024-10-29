// internal/storage/loader.go

package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/chrlesur/ontology-server/internal/parser"
)

type OntologyLoader struct {
	storage *MemoryStorage
	logger  *logger.Logger
}

func NewOntologyLoader(storage *MemoryStorage, logger *logger.Logger) *OntologyLoader {
	return &OntologyLoader{
		storage: storage,
		logger:  logger,
	}
}

// LoadFiles charge une ontologie avec ses métadonnées et contextes
func (l *OntologyLoader) LoadFiles(ontologyFile, contextFile, metadataFile string) error {
	// 1. Charger les métadonnées
	metadata, err := l.loadMetadata(metadataFile)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	// 2. Charger l'ontologie
	elements, relations, err := l.loadOntologyFile(ontologyFile)
	if err != nil {
		return fmt.Errorf("failed to load ontology file: %w", err)
	}

	// 3. Charger les contextes si présents
	if contextFile != "" {
		if err := l.enrichWithContexts(elements, contextFile); err != nil {
			return fmt.Errorf("failed to load contexts: %w", err)
		}
	}

	// 4. Créer et stocker l'ontologie
	ontology := &models.Ontology{
		ID:         fmt.Sprintf("onto_%d", time.Now().UnixNano()),
		Name:       metadata.SourceFile,
		Filename:   ontologyFile,
		Format:     filepath.Ext(ontologyFile)[1:],
		SHA256:     metadata.SHA256Hash,
		ImportedAt: time.Now(),
		Elements:   elements,
		Relations:  relations,
		Source:     metadata,
	}

	if err := l.storage.AddOntology(ontology); err != nil {
		return fmt.Errorf("failed to add ontology to storage: %w", err)
	}

	l.logger.Info(fmt.Sprintf("Successfully loaded ontology %s from %s", ontology.ID, metadata.SourceFile))
	return nil
}

func (l *OntologyLoader) loadMetadata(filename string) (*models.SourceMetadata, error) {
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

func (l *OntologyLoader) loadOntologyFile(filename string) ([]*models.OntologyElement, []*models.Relation, error) {
	switch {
	case strings.HasSuffix(filename, ".tsv"):
		elements, relations, err := parser.ParseTSV(filename)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse TSV: %w", err)
		}

		// Convertir les slices en slices de pointeurs
		elementPtrs := make([]*models.OntologyElement, len(elements))
		for i := range elements {
			elementPtrs[i] = &elements[i]
		}

		relationPtrs := make([]*models.Relation, len(relations))
		for i := range relations {
			relationPtrs[i] = &relations[i]
		}

		return elementPtrs, relationPtrs, nil

	default:
		return nil, nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(filename))
	}
}

func (l *OntologyLoader) enrichWithContexts(elements []*models.OntologyElement, contextFile string) error {
	contexts, err := parser.ParseJSON(contextFile)
	if err != nil {
		return fmt.Errorf("failed to parse context file: %w", err)
	}

	// Créer une map pour un accès rapide aux éléments
	elementMap := make(map[string]*models.OntologyElement)
	for _, elem := range elements {
		elementMap[elem.Name] = elem
	}

	// Associer les contextes aux éléments
	for _, ctx := range contexts {
		if elem, exists := elementMap[ctx.Element]; exists {
			elem.Contexts = append(elem.Contexts, ctx)
		}
	}

	return nil
}
