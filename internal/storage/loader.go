package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	l.logger.Info(fmt.Sprintf("Starting to load files: ontology=%s, context=%s, metadata=%s", ontologyFile, contextFile, metadataFile))

	// Charger les métadonnées
	metadata, err := l.loadMetadata(metadataFile)
	if err != nil {
		l.logger.Error(fmt.Sprintf("Failed to load metadata: %v", err))
		return fmt.Errorf("failed to load metadata: %w", err)
	}
	l.logger.Info("Metadata loaded successfully")

	// Charger l'ontologie
	elements, relations, err := l.loadOntologyFile(ontologyFile)
	if err != nil {
		l.logger.Error(fmt.Sprintf("Failed to load ontology file: %v", err))
		return fmt.Errorf("failed to load ontology file: %w", err)
	}
	l.logger.Info(fmt.Sprintf("Ontology loaded successfully: %d elements, %d relations", len(elements), len(relations)))

	// Charger les contextes si présents
	if contextFile != "" {
		if err := l.enrichWithContexts(elements, contextFile, metadata.Files); err != nil {
			l.logger.Error(fmt.Sprintf("Failed to load contexts: %v", err))
			return fmt.Errorf("failed to load contexts: %w", err)
		}
		l.logger.Info("Contexts loaded and associated successfully")
	} else {
		l.logger.Info("No context file provided, skipping context loading")
	}

	// Créer et stocker l'ontologie
	ontology := &models.Ontology{
		ID:         fmt.Sprintf("onto_%d", time.Now().UnixNano()),
		Name:       metadata.OntologyFile,
		Filename:   ontologyFile,
		Format:     filepath.Ext(ontologyFile)[1:],
		ImportedAt: metadata.ProcessingDate,
		Elements:   elements,
		Relations:  relations,
		Source:     metadata,
	}

	if err := l.storage.AddOntology(ontology); err != nil {
		l.logger.Error(fmt.Sprintf("Failed to add ontology to storage: %v", err))
		return fmt.Errorf("failed to add ontology to storage: %w", err)
	}
	l.logger.Info(fmt.Sprintf("Ontology added to storage successfully with ID: %s", ontology.ID))

	return nil
}

func (l *OntologyLoader) loadMetadata(filename string) (*models.SourceMetadata, error) {
	l.logger.Info(fmt.Sprintf("Loading metadata from file: %s", filename))
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata models.SourceMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	l.logger.Info(fmt.Sprintf("Metadata loaded successfully: OntologyFile=%s, ProcessingDate=%v", metadata.OntologyFile, metadata.ProcessingDate))
	return &metadata, nil
}

func (l *OntologyLoader) loadOntologyFile(filename string) ([]*models.OntologyElement, []*models.Relation, error) {
	l.logger.Info(fmt.Sprintf("Loading ontology file: %s", filename))
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

		l.logger.Info(fmt.Sprintf("TSV file parsed successfully: %d elements, %d relations", len(elementPtrs), len(relationPtrs)))
		return elementPtrs, relationPtrs, nil

	default:
		return nil, nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(filename))
	}
}

func (l *OntologyLoader) enrichWithContexts(elements []*models.OntologyElement, contextFile string, fileInfos map[string]models.FileInfo) error {
	contexts, err := parser.ParseJSON(contextFile)
	if err != nil {
		return fmt.Errorf("failed to parse context file: %w", err)
	}
	l.logger.Info(fmt.Sprintf("Parsed %d contexts from JSON file", len(contexts)))

	// Trier les contextes par position
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Position < contexts[j].Position
	})

	for _, elem := range elements {
		l.logger.Info(fmt.Sprintf("Processing element: %s", elem.Name))
		l.logger.Info(fmt.Sprintf("Element positions: %v", elem.Positions))

		elementContexts := make(map[int]*models.JSONContext)

		for _, pos := range elem.Positions {
			ctx := findContextForPosition(pos, contexts)
			if ctx != nil {
				elementContexts[ctx.Position] = ctx
				l.logger.Info(fmt.Sprintf("Found context for %s at position %d (Context position: %d, length: %d)",
					elem.Name, pos, ctx.Position, ctx.Length))
			} else {
				l.logger.Info(fmt.Sprintf("No context found for %s at position %d", elem.Name, pos))
			}
		}

		// Convertir la map en slice
		elem.Contexts = make([]models.JSONContext, 0, len(elementContexts))
		for _, ctx := range elementContexts {
			elem.Contexts = append(elem.Contexts, *ctx)
		}

		if len(elem.Contexts) > 0 {
			l.logger.Info(fmt.Sprintf("Associated %d unique contexts to element '%s'", len(elem.Contexts), elem.Name))
		} else {
			l.logger.Warning(fmt.Sprintf("No contexts found for element '%s'", elem.Name))
		}
	}

	return nil
}

func findContextForPosition(pos int, contexts []models.JSONContext) *models.JSONContext {
	for i := len(contexts) - 1; i >= 0; i-- {
		ctx := &contexts[i]
		if pos >= ctx.Position && pos < ctx.Position+ctx.Length {
			return ctx
		}
		if ctx.Position+ctx.Length < pos {
			// Comme les contextes sont triés, si nous dépassons la position recherchée, nous pouvons arrêter
			break
		}
	}
	return nil
}
