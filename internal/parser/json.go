package parser

import (
	"encoding/json"
	"fmt"
	"os"

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

	// Calculer StartOffset et EndOffset pour chaque contexte
	for i := range contexts {
		ctx := &contexts[i]
		ctx.StartOffset = ctx.Position - len(ctx.Before)
		ctx.EndOffset = ctx.Position + ctx.Length + len(ctx.After) - 1
		log.Info(fmt.Sprintf("Context for '%s': StartOffset=%d, EndOffset=%d", ctx.Element, ctx.StartOffset, ctx.EndOffset))
	}

	log.Info(fmt.Sprintf("Parsed %d contexts from JSON file", len(contexts)))
	for i, ctx := range contexts {
		log.Info(fmt.Sprintf("Context %d: Element=%s, Position=%d, StartOffset=%d, EndOffset=%d",
			i, ctx.Element, ctx.Position, ctx.StartOffset, ctx.EndOffset))
	}

	return contexts, nil
}
