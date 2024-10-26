package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/chrlesur/ontology-server/internal/models"
)

func TestParseJSON(t *testing.T) {
	// Create a temporary JSON file for testing
	testContexts := []models.JSONContext{
		{
			Position: 1,
			Before:   []string{"La"},
			After:    []string{"juridique", "sert", "de", "premier", "plan."},
			Element:  "Qualification_Juridique",
			Length:   1,
		},
		{
			Position: 8,
			Before:   []string{"La", "qualification", "juridique", "sert", "de", "premier", "plan.", "Dans"},
			After:    []string{"cas", "des", "internationales,", "le", "Conseil"},
			Element:  "Joueuse_Internationale",
			Length:   1,
		},
	}

	content, err := json.Marshal(testContexts)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	tmpfile, err := ioutil.TempFile("", "test.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Test the ParseJSON function
	contexts, err := ParseJSON(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseJSON returned an error: %v", err)
	}

	if !reflect.DeepEqual(contexts, testContexts) {
		t.Errorf("ParseJSON returned unexpected result. Got %v, want %v", contexts, testContexts)
	}
}
