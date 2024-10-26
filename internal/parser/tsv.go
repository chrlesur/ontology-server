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
func ParseTSV(filename string) ([]models.OntologyElement, error) {
	log.Info(fmt.Sprintf("Starting to parse TSV file: %s", filename))

	file, err := os.Open(filename)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open file: %v", err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'         // Use tab as delimiter
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	var elements []models.OntologyElement

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(fmt.Sprintf("Error reading TSV record: %v", err))
			return nil, fmt.Errorf("error reading TSV record: %w", err)
		}

		if len(record) < 4 {
			log.Warning(fmt.Sprintf("Skipping invalid record: %v", record))
			continue
		}

		positions, err := parsePositions(record[3])
		if err != nil {
			log.Warning(fmt.Sprintf("Error parsing positions for record %v: %v", record, err))
			continue
		}

		element := models.OntologyElement{
			Name:        strings.TrimSpace(record[0]),
			Type:        strings.TrimSpace(record[1]),
			Description: strings.TrimSpace(record[2]),
			Positions:   positions,
		}

		elements = append(elements, element)
	}

	log.Info(fmt.Sprintf("Finished parsing TSV file. Found %d elements.", len(elements)))
	return elements, nil
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
