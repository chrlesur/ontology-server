package parser

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
)

var log *logger.Logger

func init() {
	var err error
	logDir := filepath.Join(".", "logs") // Use a relative path
	log, err = logger.NewLogger(logger.INFO, logDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

// ParseTSV parses a TSV file and returns a slice of Element structures
func ParseTSV(filename string) ([]models.OntologyElement, []models.Relation, error) {
	log.Info(fmt.Sprintf("Starting to parse TSV file: %s", filename))

	file, err := os.Open(filename)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open file: %v", err))
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'         // Use tab as delimiter
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	var elements []models.OntologyElement
	var relations []models.Relation

	lineNumber := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(fmt.Sprintf("Error reading TSV record: %v", err))
			return nil, nil, fmt.Errorf("error reading TSV record: %w", err)
		}

		lineNumber++

		if len(record) < 4 {
			log.Warning(fmt.Sprintf("Skipping invalid record on line %d: %v", lineNumber, record))
			continue
		}

		// Parse positions
		var positions []int
		positionsStr := strings.Split(record[3], ",")
		for _, pos := range positionsStr {
			if p, err := strconv.Atoi(strings.TrimSpace(pos)); err == nil {
				positions = append(positions, p)
			}
		}

		// Create an OntologyElement
		element := models.OntologyElement{
			Name:        strings.TrimSpace(record[0]),
			Type:        strings.TrimSpace(record[1]),
			Description: strings.TrimSpace(record[2]),
			Positions:   positions,
		}
		elements = append(elements, element)

		// Create a Relation
		relation := models.Relation{
			Source:      strings.TrimSpace(record[0]),
			Type:        strings.TrimSpace(record[1]),
			Target:      strings.TrimSpace(record[2]),
			Description: strings.TrimSpace(record[2]),
		}
		relations = append(relations, relation)

		log.Info(fmt.Sprintf("Parsed line %d: Element: %s, Type: %s, Description: %s",
			lineNumber, element.Name, element.Type, element.Description))
	}
	log.Info(fmt.Sprintf("Finished parsing TSV file. Found %d elements and %d relations.", len(elements), len(relations)))
	return elements, relations, nil
}

func parsePositions(positionsStr string) ([]int, error) {
	positionsStr = strings.TrimSpace(positionsStr)
	if positionsStr == "" {
		return []int{}, nil
	}

	positionStrs := strings.Split(positionsStr, ",")
	positions := make([]int, 0, len(positionStrs))

	for _, pos := range positionStrs {
		pos = strings.TrimSpace(pos)
		if pos == "" {
			continue
		}
		position, err := strconv.Atoi(pos)
		if err != nil {
			return nil, fmt.Errorf("invalid position: %s", pos)
		}
		positions = append(positions, position)
	}

	return positions, nil
}
