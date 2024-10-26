package parser

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseRDF(t *testing.T) {
	// Créer un fichier RDF temporaire pour le test
	content := `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#">
    <rdf:Description rdf:about="http://example.org/Element1">
        <rdf:type rdf:resource="http://example.org/Type1"/>
        <rdfs:label>Element1</rdfs:label>
        <rdfs:comment>Description1</rdfs:comment>
    </rdf:Description>
    <rdf:Description rdf:about="http://example.org/Relation1">
        <rdf:type rdf:resource="http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"/>
        <rdfs:domain rdf:resource="http://example.org/Element1"/>
        <rdfs:range rdf:resource="http://example.org/Element2"/>
    </rdf:Description>
</rdf:RDF>`

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

	// Tester la fonction ParseRDF
	elements, relations, err := ParseRDF(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseRDF returned an error: %v", err)
	}

	// Vérifier les éléments
	if len(elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(elements))
	} else if elements[0].Name != "Element1" || elements[0].Description != "Description1" {
		t.Errorf("Unexpected element: %v", elements[0])
	}

	// Vérifier les relations
	if len(relations) != 1 {
		t.Errorf("Expected 1 relation, got %d", len(relations))
	} else if relations[0].Type != "http://example.org/Relation1" ||
		relations[0].Source != "http://example.org/Element1" ||
		relations[0].Target != "http://example.org/Element2" {
		t.Errorf("Unexpected relation: %v", relations[0])
	}
}
