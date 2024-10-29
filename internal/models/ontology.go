package models

import (
	"time"
)

// OntologyElement représente un élément individuel dans l'ontologie
type OntologyElement struct {
	Name         string
	OriginalName string
	Type         string
	Positions    []int
	Description  string
	Contexts     []JSONContext
}

// Relation représente une relation entre deux éléments de l'ontologie
type Relation struct {
	Source      string
	Type        string
	Target      string
	Description string
}

// SourceMetadata représente les métadonnées du fichier source
type SourceMetadata struct {
	SourceFile     string    `json:"source_file"`
	Directory      string    `json:"directory"`
	FileDate       time.Time `json:"file_date"`
	SHA256Hash     string    `json:"sha256_hash"`
	OntologyFile   string    `json:"ontology_file"`
	ContextFile    string    `json:"context_file"`
	ProcessingDate time.Time `json:"processing_date"`
}

// Ontology représente une ontologie complète
type Ontology struct {
	ID         string
	Name       string
	Filename   string
	Format     string
	Size       int64
	SHA256     string
	ImportedAt time.Time
	Elements   []*OntologyElement
	Relations  []*Relation
	Source     *SourceMetadata
}
