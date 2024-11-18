package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chrlesur/ontology-server/internal/models"
)

func TestNewMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()
	if ms == nil {
		t.Error("NewMemoryStorage returned nil")
	}
	if ms.ontologies == nil {
		t.Error("ontologies map is not initialized")
	}
}

func TestAddOntology(t *testing.T) {
	ms := NewMemoryStorage()
	ontology := &models.Ontology{ID: "test1", Name: "Test Ontology"}

	err := ms.AddOntology(ontology)
	if err != nil {
		t.Errorf("Failed to add ontology: %v", err)
	}

	// Try to add the same ontology again
	err = ms.AddOntology(ontology)
	if err == nil {
		t.Error("Expected error when adding duplicate ontology, got nil")
	}
}

func TestGetOntology(t *testing.T) {
	ms := NewMemoryStorage()
	ontology := &models.Ontology{ID: "test1", Name: "Test Ontology"}
	ms.AddOntology(ontology)

	retrieved, err := ms.GetOntology("test1")
	if err != nil {
		t.Errorf("Failed to get ontology: %v", err)
	}
	if retrieved.ID != ontology.ID || retrieved.Name != ontology.Name {
		t.Error("Retrieved ontology does not match the original")
	}

	_, err = ms.GetOntology("nonexistent")
	if err == nil {
		t.Error("Expected error when getting nonexistent ontology, got nil")
	}
}

func TestUpdateOntology(t *testing.T) {
	ms := NewMemoryStorage()
	ontology := &models.Ontology{ID: "test1", Name: "Test Ontology"}
	ms.AddOntology(ontology)

	updatedOntology := &models.Ontology{ID: "test1", Name: "Updated Test Ontology"}
	err := ms.UpdateOntology(updatedOntology)
	if err != nil {
		t.Errorf("Failed to update ontology: %v", err)
	}

	retrieved, _ := ms.GetOntology("test1")
	if retrieved.Name != "Updated Test Ontology" {
		t.Error("Ontology was not updated correctly")
	}

	nonexistentOntology := &models.Ontology{ID: "nonexistent", Name: "Nonexistent"}
	err = ms.UpdateOntology(nonexistentOntology)
	if err == nil {
		t.Error("Expected error when updating nonexistent ontology, got nil")
	}
}

func TestDeleteOntology(t *testing.T) {
	ms := NewMemoryStorage()
	ontology := &models.Ontology{ID: "test1", Name: "Test Ontology"}
	ms.AddOntology(ontology)

	err := ms.DeleteOntology("test1")
	if err != nil {
		t.Errorf("Failed to delete ontology: %v", err)
	}

	_, err = ms.GetOntology("test1")
	if err == nil {
		t.Error("Expected error when getting deleted ontology, got nil")
	}

	err = ms.DeleteOntology("nonexistent")
	if err == nil {
		t.Error("Expected error when deleting nonexistent ontology, got nil")
	}
}

func TestListOntologies(t *testing.T) {
	ms := NewMemoryStorage()
	ontology1 := &models.Ontology{ID: "test1", Name: "Test Ontology 1"}
	ontology2 := &models.Ontology{ID: "test2", Name: "Test Ontology 2"}
	ms.AddOntology(ontology1)
	ms.AddOntology(ontology2)

	ontologies := ms.ListOntologies()
	if len(ontologies) != 2 {
		t.Errorf("Expected 2 ontologies, got %d", len(ontologies))
	}

	foundOntology1 := false
	foundOntology2 := false
	for _, o := range ontologies {
		if o.ID == "test1" {
			foundOntology1 = true
		}
		if o.ID == "test2" {
			foundOntology2 = true
		}
	}

	if !foundOntology1 || !foundOntology2 {
		t.Error("ListOntologies did not return all added ontologies")
	}
}

func TestConcurrency(t *testing.T) {
	ms := NewMemoryStorage()
	concurrentOperations := 1000

	// Concurrent additions
	for i := 0; i < concurrentOperations; i++ {
		go func(id int) {
			ontology := &models.Ontology{ID: fmt.Sprintf("test%d", id), Name: fmt.Sprintf("Test Ontology %d", id)}
			ms.AddOntology(ontology)
		}(i)
	}

	time.Sleep(time.Second) // Give time for goroutines to complete

	// Verify all ontologies were added
	ontologies := ms.ListOntologies()
	if len(ontologies) != concurrentOperations {
		t.Errorf("Expected %d ontologies, got %d", concurrentOperations, len(ontologies))
	}

	// Concurrent reads and updates
	for i := 0; i < concurrentOperations; i++ {
		go func(id int) {
			ontologyID := fmt.Sprintf("test%d", id)
			ms.GetOntology(ontologyID)
			updatedOntology := &models.Ontology{ID: ontologyID, Name: fmt.Sprintf("Updated Test Ontology %d", id)}
			ms.UpdateOntology(updatedOntology)
		}(i)
	}

	time.Sleep(time.Second) // Give time for goroutines to complete

	// Verify all ontologies were updated
	for i := 0; i < concurrentOperations; i++ {
		ontology, _ := ms.GetOntology(fmt.Sprintf("test%d", i))
		if ontology.Name != fmt.Sprintf("Updated Test Ontology %d", i) {
			t.Errorf("Ontology %d was not updated correctly", i)
		}
	}
}
func TestLoadOntologyWithMetadata(t *testing.T) {
	ms := NewMemoryStorage()

	// Créer un fichier de métadonnées temporaire avec la nouvelle structure
	metadataContent := `{
        "ontology_file": "test.tsv",
        "processing_date": "2024-10-29T08:54:33.959Z",
        "files": {
            "file1": {
                "id": "file1",
                "source_file": "test.txt",
                "directory": "/test/path",
                "file_date": "2024-10-29T08:54:33.959Z",
                "sha256_hash": "test-hash"
            }
        }
    }`

	metadataFile := filepath.Join(t.TempDir(), "metadata.json")
	if err := os.WriteFile(metadataFile, []byte(metadataContent), 0644); err != nil {
		t.Fatalf("Failed to create test metadata file: %v", err)
	}

	// Créer un fichier TSV minimal pour le test
	tsvContent := "Element1\tType1\tDescription1\t1,2,3"
	tsvFile := filepath.Join(t.TempDir(), "test.tsv")
	if err := os.WriteFile(tsvFile, []byte(tsvContent), 0644); err != nil {
		t.Fatalf("Failed to create test TSV file: %v", err)
	}

	// Test du chargement
	err := ms.LoadOntologyFromFile(tsvFile, "", metadataFile)
	if err != nil {
		t.Errorf("Failed to load ontology with metadata: %v", err)
	}

	// Vérifier que les ontologies sont chargées avec les métadonnées
	ontologies := ms.ListOntologies()
	if len(ontologies) != 1 {
		t.Errorf("Expected 1 ontology, got %d", len(ontologies))
	}

	// Vérifier les métadonnées
	ontology := ontologies[0]
	if ontology.Source == nil {
		t.Error("Expected source metadata to be present")
	} else {
		if ontology.Source.OntologyFile != "test.tsv" {
			t.Errorf("Expected ontology file 'test.tsv', got '%s'", ontology.Source.OntologyFile)
		}
		if len(ontology.Source.Files) != 1 {
			t.Errorf("Expected 1 file in metadata, got %d", len(ontology.Source.Files))
		}
		fileInfo, exists := ontology.Source.Files["file1"]
		if !exists {
			t.Error("Expected file info for 'file1' to exist")
		} else {
			if fileInfo.SourceFile != "test.txt" {
				t.Errorf("Expected source file 'test.txt', got '%s'", fileInfo.SourceFile)
			}
			if fileInfo.SHA256Hash != "test-hash" {
				t.Errorf("Expected SHA256 hash 'test-hash', got '%s'", fileInfo.SHA256Hash)
			}
		}
	}
}
