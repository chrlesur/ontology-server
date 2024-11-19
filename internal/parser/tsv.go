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
	reader.Comma = '\t'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	var elements []models.OntologyElement
	var relations []models.Relation
	elementMap := make(map[string]*models.OntologyElement)

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

		if len(record) < 3 {
			log.Warning(fmt.Sprintf("Skipping invalid record on line %d: %v", lineNumber, record))
			continue
		}

		name := strings.TrimSpace(record[0])
		elemType := strings.TrimSpace(record[1])
		thirdField := strings.TrimSpace(record[2])

		if len(record) >= 3 && !strings.Contains(elemType, ":") { // C'est probablement un élément
			description := thirdField
			var positions []int
			if len(record) > 3 {
				positionsStr := strings.TrimSpace(record[3])
				if positionsStr != "" {
					positions, err = parsePositions(positionsStr)
					if err != nil {
						log.Warning(fmt.Sprintf("Error parsing positions on line %d: %v", lineNumber, err))
					}
				}
			}

			element := models.OntologyElement{
				Name:        name,
				Type:        elemType,
				Description: description,
				Positions:   positions,
				Contexts:    []models.JSONContext{},
			}
			elements = append(elements, element)
			elementMap[name] = &element

			log.Info(fmt.Sprintf("Parsed element on line %d: Name: %s, Type: %s, Description: %s, Positions: %v",
				lineNumber, element.Name, element.Type, element.Description, element.Positions))

		} else { // C'est une relation
			target := thirdField
			description := ""
			if len(record) > 3 {
				description = strings.TrimSpace(record[3])
			}

			relation := models.Relation{
				Source:      name,
				Type:        elemType,
				Target:      target,
				Description: description,
			}
			relations = append(relations, relation)

			log.Info(fmt.Sprintf("Parsed relation on line %d: Source: %s, Type: %s, Target: %s",
				lineNumber, relation.Source, relation.Type, relation.Target))
		}
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
			// Log the error but continue processing other positions
			log.Warning(fmt.Sprintf("Invalid position value: %s", pos))
			continue
		}
		positions = append(positions, position)
	}

	return positions, nil
}
