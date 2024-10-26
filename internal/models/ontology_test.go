package models

import (
	"testing"
	"time"
)

func TestOntologyStructure(t *testing.T) {
	// Créer une nouvelle ontologie
	ontology := Ontology{
		ID:         "test-ontology",
		Name:       "Test Ontology",
		Filename:   "test.tsv",
		Format:     "TSV",
		Size:       1024,
		SHA256:     "abcdef1234567890",
		ImportedAt: time.Now(),
		Elements: []*OntologyElement{
			{
				Name:        "TestElement",
				Type:        "Concept",
				Positions:   []int{1, 2, 3},
				Description: "This is a test element",
			},
		},
		Relations: []*Relation{
			{
				Source:      "TestElement",
				Type:        "isA",
				Target:      "ParentElement",
				Description: "Test relation",
			},
		},
	}

	// Vérifier que les champs sont correctement définis
	if ontology.ID != "test-ontology" {
		t.Errorf("Expected ID to be 'test-ontology', got '%s'", ontology.ID)
	}

	if len(ontology.Elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(ontology.Elements))
	}

	if len(ontology.Relations) != 1 {
		t.Errorf("Expected 1 relation, got %d", len(ontology.Relations))
	}
}

func TestJSONContextStructure(t *testing.T) {
	// Tester la structure JSONContext
	context := JSONContext{
		Position: 1,
		Before:   []string{"before"},
		After:    []string{"after"},
		Element:  "TestElement",
		Length:   1,
	}

	if context.Element != "TestElement" {
		t.Errorf("Expected Element to be 'TestElement', got '%s'", context.Element)
	}
}
