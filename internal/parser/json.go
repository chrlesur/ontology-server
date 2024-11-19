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

	for i := range contexts {
		ctx := &contexts[i]
		ctx.StartOffset = ctx.FilePosition
		ctx.EndOffset = ctx.FilePosition + ctx.Length - 1

		// Vérification supplémentaire pour s'assurer que FilePosition est défini
		if ctx.FilePosition == 0 {
			ctx.FilePosition = ctx.Position
			log.Warning(fmt.Sprintf("FilePosition was 0 for context of '%s', using Position value: %d", ctx.Element, ctx.Position))
		}

		log.Info(fmt.Sprintf("Context for '%s': FileID=%s, FilePosition=%d, Position=%d, StartOffset=%d, EndOffset=%d, Length=%d",
			ctx.Element, ctx.FileID, ctx.FilePosition, ctx.Position, ctx.StartOffset, ctx.EndOffset, ctx.Length))
	}

	log.Info(fmt.Sprintf("Parsed %d contexts from JSON file", len(contexts)))

	return contexts, nil
}
