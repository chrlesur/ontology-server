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
	"unicode"

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
	relationMap := make(map[string]models.Relation)

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

		name := strings.TrimSpace(record[0])
		elemType := strings.TrimSpace(record[1])
		positionsStr := record[len(record)-1]
		description := strings.Join(record[2:len(record)-1], "\t")

		positions, err := parsePositions(positionsStr)
		if err != nil {
			log.Warning(fmt.Sprintf("Error parsing positions on line %d: %v", lineNumber, err))
		}

		element := models.OntologyElement{
			Name:        name,
			Type:        elemType,
			Description: strings.TrimSpace(description),
			Positions:   positions,
		}
		elements = append(elements, element)

		// Dédupliquer les relations
		relationKey := fmt.Sprintf("%s|%s|%s", name, elemType, strings.TrimSpace(record[2]))
		if _, exists := relationMap[relationKey]; !exists {
			relationMap[relationKey] = models.Relation{
				Source:      name,
				Type:        elemType,
				Target:      strings.TrimSpace(record[2]),
				Description: strings.TrimSpace(description),
			}
		}

		log.Info(fmt.Sprintf("Parsed line %d: Element: %s, Type: %s, Description: %s",
			lineNumber, element.Name, element.Type, element.Description))
	}

	// Convertir la map de relations en slice
	relations := make([]models.Relation, 0, len(relationMap))
	for _, relation := range relationMap {
		relations = append(relations, relation)
	}

	log.Info(fmt.Sprintf("Finished parsing TSV file. Found %d elements and %d unique relations.", len(elements), len(relations)))
	return elements, relations, nil
}

func parsePositions(positionsStr string) ([]int, error) {
	positionsStr = strings.TrimSpace(positionsStr)
	if positionsStr == "" {
		return []int{}, nil
	}

	// Vérifier si la chaîne contient des caractères non numériques (à l'exception des virgules et des espaces)
	if strings.IndexFunc(positionsStr, func(r rune) bool {
		return !unicode.IsDigit(r) && r != ',' && r != ' '
	}) != -1 {
		// Si c'est le cas, retourner un slice vide sans erreur
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
			// Ignorer les erreurs de conversion et continuer
			continue
		}
		positions = append(positions, position)
	}

	return positions, nil
}
