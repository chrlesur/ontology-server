package parser

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chrlesur/ontology-server/internal/models"
)

func TestParseRDF(t *testing.T) {
	// Create a temporary RDF file for testing
	content := `
    <?xml version="1.0"?>
    <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
             xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#">
      <rdf:Description rdf:about="http://example.org/resource1">
        <rdf:type rdf:resource="http://example.org/Type1"/>
        <rdfs:label>Resource 1</rdfs:label>
        <rdfs:comment>This is resource 1</rdfs:comment>
      </rdf:Description>
    </rdf:RDF>
    `

	tmpfile, err := ioutil.TempFile("", "test.rdf")
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

	// Test the ParseRDF function
	elements, err := ParseRDF(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseRDF returned an error: %v", err)
	}

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	expected := models.OntologyElement{
		Name:        "Resource 1",
		Type:        "http://example.org/Type1",
		Description: "This is resource 1",
	}

	if elements[0].Name != expected.Name || elements[0].Type != expected.Type || elements[0].Description != expected.Description {
		t.Errorf("ParseRDF returned unexpected result. Got %v, want %v", elements[0], expected)
	}
}
