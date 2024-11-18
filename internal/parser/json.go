package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/chrlesur/ontology-server/internal/models"
)

// ParseJSON parses a JSON file and returns a slice of JSONContext structures
func ParseJSON(filename string) ([]models.JSONContext, error) {
	log.Info(fmt.Sprintf("Starting to parse JSON file: %s", filename))

	file, err := os.Open(filename)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open file: %v", err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var contexts []models.JSONContext
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&contexts)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to decode JSON: %v", err))
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Trier les contextes par position
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Position < contexts[j].Position
	})

	// Reconstruire les contextes complets
	for i := range contexts {
		ctx := &contexts[i]
		ctx.StartOffset = ctx.FilePosition - len(ctx.Before)
		ctx.EndOffset = ctx.FilePosition + ctx.Length + len(ctx.After) - 1

		// Compléter le contexte avec les éléments adjacents si nécessaire
		if i > 0 {
			prevCtx := contexts[i-1]
			missingBefore := prevCtx.EndOffset + 1 - ctx.StartOffset
			if missingBefore > 0 && prevCtx.FileID == ctx.FileID {
				ctx.Before = append(prevCtx.After[len(prevCtx.After)-missingBefore:], ctx.Before...)
				ctx.StartOffset = prevCtx.EndOffset + 1
			}
		}
		if i < len(contexts)-1 {
			nextCtx := contexts[i+1]
			missingAfter := nextCtx.StartOffset - ctx.EndOffset - 1
			if missingAfter > 0 && nextCtx.FileID == ctx.FileID {
				ctx.After = append(ctx.After, nextCtx.Before[:missingAfter]...)
				ctx.EndOffset = nextCtx.StartOffset - 1
			}
		}

		log.Info(fmt.Sprintf("Context for '%s': FileID=%s, FilePosition=%d, StartOffset=%d, EndOffset=%d",
			ctx.Element, ctx.FileID, ctx.FilePosition, ctx.StartOffset, ctx.EndOffset))
	}

	log.Info(fmt.Sprintf("Parsed %d contexts from JSON file", len(contexts)))
	for i, ctx := range contexts {
		log.Info(fmt.Sprintf("Context %d: Element=%s, FileID=%s, FilePosition=%d, StartOffset=%d, EndOffset=%d",
			i, ctx.Element, ctx.FileID, ctx.FilePosition, ctx.StartOffset, ctx.EndOffset))
	}

	return contexts, nil
}
