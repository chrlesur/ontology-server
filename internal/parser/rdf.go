package parser

import (
	"fmt"
	"os"

	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/knakk/rdf"
)

// ParseRDF parses an RDF file and returns slices of OntologyElement and Relation
func ParseRDF(filename string) ([]models.OntologyElement, []models.Relation, error) {
	log.Info(fmt.Sprintf("Starting to parse RDF file: %s", filename))

	file, err := os.Open(filename)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open file: %v", err))
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := rdf.NewTripleDecoder(file, rdf.RDFXML)

	elements := make(map[string]*models.OntologyElement)
	relations := make(map[string]*models.Relation)

	for triple, err := decoder.Decode(); err == nil; triple, err = decoder.Decode() {
		subj := triple.Subj.String()
		pred := triple.Pred.String()
		obj := triple.Obj.String()

		switch pred {
		case "http://www.w3.org/1999/02/22-rdf-syntax-ns#type":
			if obj != "http://www.w3.org/1999/02/22-rdf-syntax-ns#Property" {
				if _, exists := elements[subj]; !exists {
					elements[subj] = &models.OntologyElement{Name: subj, Type: obj}
				}
			} else {
				if _, exists := relations[subj]; !exists {
					relations[subj] = &models.Relation{Type: subj}
				}
			}
		case "http://www.w3.org/2000/01/rdf-schema#label":
			if elem, exists := elements[subj]; exists {
				elem.Name = obj
			}
		case "http://www.w3.org/2000/01/rdf-schema#comment":
			if elem, exists := elements[subj]; exists {
				elem.Description = obj
			}
		case "http://www.w3.org/2000/01/rdf-schema#domain":
			if rel, exists := relations[subj]; exists {
				rel.Source = obj
			}
		case "http://www.w3.org/2000/01/rdf-schema#range":
			if rel, exists := relations[subj]; exists {
				rel.Target = obj
			}
		}

		log.Info(fmt.Sprintf("Parsed triple: %s -> %s -> %s", subj, pred, obj))
	}

	// Convert maps to slices
	elementSlice := make([]models.OntologyElement, 0, len(elements))
	for _, elem := range elements {
		elementSlice = append(elementSlice, *elem)
	}
	relationSlice := make([]models.Relation, 0, len(relations))
	for _, rel := range relations {
		relationSlice = append(relationSlice, *rel)
	}

	log.Info(fmt.Sprintf("Finished parsing RDF file. Found %d elements and %d relations.", len(elementSlice), len(relationSlice)))
	return elementSlice, relationSlice, nil
}
