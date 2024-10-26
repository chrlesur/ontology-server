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

	log.Info(fmt.Sprintf("Finished parsing JSON file. Found %d contexts.", len(contexts)))
	return contexts, nil
}
