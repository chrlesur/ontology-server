package parser

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/chrlesur/ontology-server/internal/models"
)

func TestParseTSV(t *testing.T) {
	// Créer un fichier TSV temporaire pour le test
	content := `Element1	Type1	Target1	Description1
Element2	Type2	Target2	Description2
Element3	Type3	Target3	Description3`

	tmpfile, err := ioutil.TempFile("", "test.tsv")
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

	// Tester la fonction ParseTSV
	elements, relations, err := ParseTSV(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseTSV returned an error: %v", err)
	}

	// Vérifier les éléments
	expectedElements := []models.OntologyElement{
		{Name: "Element1", Type: "Type1", Description: "Description1"},
		{Name: "Element2", Type: "Type2", Description: "Description2"},
		{Name: "Element3", Type: "Type3", Description: "Description3"},
	}

	if !reflect.DeepEqual(elements, expectedElements) {
		t.Errorf("ParseTSV returned unexpected elements. Got %v, want %v", elements, expectedElements)
	}

	// Vérifier les relations
	expectedRelations := []models.Relation{
		{Source: "Element1", Type: "Type1", Target: "Target1", Description: "Description1"},
		{Source: "Element2", Type: "Type2", Target: "Target2", Description: "Description2"},
		{Source: "Element3", Type: "Type3", Target: "Target3", Description: "Description3"},
	}

	if !reflect.DeepEqual(relations, expectedRelations) {
		t.Errorf("ParseTSV returned unexpected relations. Got %v, want %v", relations, expectedRelations)
	}
}
