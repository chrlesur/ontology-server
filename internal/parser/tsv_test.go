package parser

import (
	"io/ioutil"
	"os"
	"reflect"
	"sort"
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
		{Name: "Element1", Type: "Type1", Description: "Description1", Positions: []int{}, Contexts: []models.JSONContext{}},
		{Name: "Element2", Type: "Type2", Description: "Description2", Positions: []int{}, Contexts: []models.JSONContext{}},
		{Name: "Element3", Type: "Type3", Description: "Description3", Positions: []int{}, Contexts: []models.JSONContext{}},
	}

	if len(elements) != len(expectedElements) {
		t.Errorf("ParseTSV returned %d elements, expected %d", len(elements), len(expectedElements))
	} else {
		for i, elem := range elements {
			if !compareElements(t, elem, expectedElements[i], i) {
				// L'erreur est déjà rapportée dans compareElements
			}
		}
	}

	// Vérifier les relations
	expectedRelations := []models.Relation{
		{Source: "Element1", Type: "Type1", Target: "Target1", Description: "Description1"},
		{Source: "Element2", Type: "Type2", Target: "Target2", Description: "Description2"},
		{Source: "Element3", Type: "Type3", Target: "Target3", Description: "Description3"},
	}

	// Trier les relations obtenues et attendues
	sortRelations(relations)
	sortRelations(expectedRelations)

	if !reflect.DeepEqual(relations, expectedRelations) {
		t.Errorf("ParseTSV returned unexpected relations.")
		for i, rel := range relations {
			t.Errorf("Relation %d: Got %+v, Expected %+v", i, rel, expectedRelations[i])
		}
	}
}

func compareElements(t *testing.T, got, want models.OntologyElement, index int) bool {
	equal := true
	if got.Name != want.Name {
		t.Errorf("Element %d: Name mismatch. Got %q, want %q", index, got.Name, want.Name)
		equal = false
	}
	if got.OriginalName != want.OriginalName {
		t.Errorf("Element %d: OriginalName mismatch. Got %q, want %q", index, got.OriginalName, want.OriginalName)
		equal = false
	}
	if got.Type != want.Type {
		t.Errorf("Element %d: Type mismatch. Got %q, want %q", index, got.Type, want.Type)
		equal = false
	}
	if got.Description != want.Description {
		t.Errorf("Element %d: Description mismatch. Got %q, want %q", index, got.Description, want.Description)
		equal = false
	}
	if !reflect.DeepEqual(got.Positions, want.Positions) {
		t.Errorf("Element %d: Positions mismatch. Got %v, want %v", index, got.Positions, want.Positions)
		equal = false
	}
	if len(got.Contexts) != len(want.Contexts) {
		t.Errorf("Element %d: Contexts length mismatch. Got %d, want %d", index, len(got.Contexts), len(want.Contexts))
		equal = false
	} else {
		for i, gotCtx := range got.Contexts {
			wantCtx := want.Contexts[i]
			if !reflect.DeepEqual(gotCtx, wantCtx) {
				t.Errorf("Element %d: Context %d mismatch. Got %+v, want %+v", index, i, gotCtx, wantCtx)
				equal = false
			}
		}
	}
	return equal
}

func sortRelations(relations []models.Relation) {
	sort.Slice(relations, func(i, j int) bool {
		return relations[i].Source < relations[j].Source
	})
}

func TestParseTSVEmpty(t *testing.T) {
	content := ""
	tmpfile, err := ioutil.TempFile("", "test_empty.tsv")
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

	elements, relations, err := ParseTSV(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseTSV returned an error for empty file: %v", err)
	}
	if len(elements) != 0 || len(relations) != 0 {
		t.Errorf("Expected 0 elements and relations, got %d elements and %d relations", len(elements), len(relations))
	}
}

func TestParseTSVMalformedLines(t *testing.T) {
	content := `Element1	Type1	Target1	Description1
InvalidLine
Element2	Type2	Target2	Description2	ExtraField
Element3`

	tmpfile, err := ioutil.TempFile("", "test_malformed.tsv")
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

	elements, relations, err := ParseTSV(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseTSV returned an error: %v", err)
	}
	if len(elements) != 2 || len(relations) != 2 {
		t.Errorf("Expected 2 elements and relations, got %d elements and %d relations", len(elements), len(relations))
	}
}

func TestParseTSVInvalidPositions(t *testing.T) {
	content := `Element1	Type1	Target1	Description1	1,2,3
Element2	Type2	Target2	Description2	invalid,4,5
Element3	Type3	Target3	Description3	6,7,8`

	tmpfile, err := ioutil.TempFile("", "test_invalid_positions.tsv")
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

	elements, _, err := ParseTSV(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseTSV returned an error: %v", err)
	}
	if len(elements) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(elements))
	}
	if len(elements[1].Positions) != 2 {
		t.Errorf("Expected 2 valid positions for Element2, got %d", len(elements[1].Positions))
	}
	expectedPositions := []int{4, 5}
	if !reflect.DeepEqual(elements[1].Positions, expectedPositions) {
		t.Errorf("Expected positions %v for Element2, got %v", expectedPositions, elements[1].Positions)
	}
}
