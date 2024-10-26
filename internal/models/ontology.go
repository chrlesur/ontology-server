package models

import (
	"time"
)

// OntologyElement représente un élément individuel dans l'ontologie
type OntologyElement struct {
	Name        string
	Type        string
	Positions   []int
	Description string
}

// Relation représente une relation entre deux éléments de l'ontologie
type Relation struct {
	Source      string
	Type        string
	Target      string
	Description string
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
}
