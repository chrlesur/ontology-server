package parser

import (
	"fmt"
	"os"

	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/knakk/rdf"
)

// ParseRDF parses an RDF file and returns a slice of OntologyElement structures
func ParseRDF(filename string) ([]models.OntologyElement, error) {
	log.Info(fmt.Sprintf("Starting to parse RDF file: %s", filename))

	file, err := os.Open(filename)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open file: %v", err))
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := rdf.NewTripleDecoder(file, rdf.RDFXML)
	elements := make(map[string]*models.OntologyElement)

	for triple, err := decoder.Decode(); err == nil; triple, err = decoder.Decode() {
		subj := triple.Subj.String()
		pred := triple.Pred.String()
		obj := triple.Obj.String()

		if _, exists := elements[subj]; !exists {
			elements[subj] = &models.OntologyElement{Name: subj}
		}

		switch pred {
		case "http://www.w3.org/1999/02/22-rdf-syntax-ns#type":
			elements[subj].Type = obj
		case "http://www.w3.org/2000/01/rdf-schema#label":
			elements[subj].Name = obj
		case "http://www.w3.org/2000/01/rdf-schema#comment":
			elements[subj].Description = obj
		}
	}

	result := make([]models.OntologyElement, 0, len(elements))
	for _, element := range elements {
		result = append(result, *element)
	}

	log.Info(fmt.Sprintf("Finished parsing RDF file. Found %d elements.", len(result)))
	return result, nil
}
