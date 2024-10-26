package parser

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/chrlesur/ontology-server/internal/models"
)

func TestParseTSV(t *testing.T) {
	// Create a temporary TSV file for testing
	content := `Element1	Type1	Description1	1,2,3
Element2	Type2	Description2	4,5,6
Invalid	Row
Element3	Type3	Description3	7,8,9`

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

	// Test the ParseTSV function
	elements, err := ParseTSV(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseTSV returned an error: %v", err)
	}

	expected := []models.OntologyElement{
		{Name: "Element1", Type: "Type1", Description: "Description1", Positions: []int{1, 2, 3}},
		{Name: "Element2", Type: "Type2", Description: "Description2", Positions: []int{4, 5, 6}},
		{Name: "Element3", Type: "Type3", Description: "Description3", Positions: []int{7, 8, 9}},
	}

	if !reflect.DeepEqual(elements, expected) {
		t.Errorf("ParseTSV returned unexpected result. Got %v, want %v", elements, expected)
	}
}

func TestParsePositions(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"1,2,3", []int{1, 2, 3}, false},
		{"", []int{}, false},
		{"1, 2, 3", []int{1, 2, 3}, false},
		{"1,a,3", nil, true},
		{" 1 , 2 , 3 ", []int{1, 2, 3}, false},
	}

	for _, tt := range tests {
		result, err := parsePositions(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parsePositions(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("parsePositions(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
