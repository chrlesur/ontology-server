package parser

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrlesur/ontology-server/internal/models"
)

func TestParseOWL(t *testing.T) {
	// Create a temporary OWL file for testing
	content := `
    <?xml version="1.0"?>
    <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
             xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#"
             xmlns:owl="http://www.w3.org/2002/07/owl#">
      <owl:Class rdf:about="http://example.org/Class1">
        <rdfs:label>Class 1</rdfs:label>
        <rdfs:comment>This is class 1</rdfs:comment>
      </owl:Class>
    </rdf:RDF>
    `

	tmpfile, err := ioutil.TempFile("", "test.owl")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Test the ParseOWL function
	elements, err := ParseOWL(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseOWL returned an error: %v", err)
	}

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	expected := models.OntologyElement{
		Name:        "Class 1",
		Type:        "http://www.w3.org/2002/07/owl#Class",
		Description: "This is class 1",
	}

	if elements[0].Name != expected.Name || elements[0].Type != expected.Type || elements[0].Description != expected.Description {
		t.Errorf("ParseOWL returned unexpected result. Got %v, want %v", elements[0], expected)
	}
}
