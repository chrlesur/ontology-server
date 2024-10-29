package search

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/models"
	"github.com/chrlesur/ontology-server/internal/storage"
)

// SearchEngine représente le moteur de recherche
type SearchEngine struct {
	Storage *storage.MemoryStorage
	Logger  *logger.Logger
}

// NewSearchEngine crée une nouvelle instance de SearchEngine
func NewSearchEngine(storage *storage.MemoryStorage, logger *logger.Logger) *SearchEngine {
	return &SearchEngine{
		Storage: storage,
		Logger:  logger,
	}
}

// SearchResult représente un résultat de recherche
type SearchResult struct {
	OntologyID  string
	ElementName string
	ElementType string
	Description string
	Context     string
	Position    int
	Relevance   float64
	Contexts    []models.JSONContext
	Source      *models.SourceMetadata
}

// Search effectue une recherche dans les ontologies
func (se *SearchEngine) Search(query string, ontologyID string, elementType string, contextSize int) ([]SearchResult, error) {
	se.Logger.Info("Starting search with query: " + query)
	query = strings.ToLower(query)
	var results []SearchResult
	var wg sync.WaitGroup
	resultChan := make(chan SearchResult)

	ontologies := se.Storage.ListOntologies()
	se.Logger.Info(fmt.Sprintf("Searching through %d ontologies", len(ontologies)))

	for _, ontology := range ontologies {
		if ontologyID != "" && ontology.ID != ontologyID {
			continue
		}

		wg.Add(1)
		go func(onto *models.Ontology) {
			defer wg.Done()
			se.Logger.Info(fmt.Sprintf("Searching in ontology: %s (Elements: %d)", onto.ID, len(onto.Elements)))
			for _, element := range onto.Elements {
				se.Logger.Info(fmt.Sprintf("Examining element: %s (Type: %s, Contexts: %d)",
					element.Name, element.Type, len(element.Contexts)))

				if elementType != "" && element.Type != elementType {
					continue
				}

				relevance := calculateRelevance(query, element)
				if relevance > 0.3 {
					context := extractContext(element, contextSize)
					position := 0
					if len(element.Positions) > 0 {
						position = element.Positions[0]
					}
					result := SearchResult{
						OntologyID:  onto.ID,
						ElementName: element.Name,
						ElementType: element.Type,
						Description: element.Description,
						Context:     context,
						Position:    position,
						Relevance:   relevance,
						Contexts:    element.Contexts,
						Source:      onto.Source,
					}
					se.Logger.Info(fmt.Sprintf("Found relevant result: %s (Relevance: %.2f, Contexts: %d)",
						result.ElementName, result.Relevance, len(result.Contexts)))
					resultChan <- result
				}
			}
		}(ontology)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		results = append(results, result)
	}

	sortSearchResults(results)

	se.Logger.Info(fmt.Sprintf("Search completed. Found %d results.", len(results)))
	for _, result := range results {
		se.Logger.Info(fmt.Sprintf("Result: %s (Contexts: %d)", result.ElementName, len(result.Contexts)))
	}
	return results, nil
}

// sortSearchResults trie les résultats de recherche par pertinence décroissante
func sortSearchResults(results []SearchResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})
}

// calculateRelevance calcule la pertinence d'un élément par rapport à la requête
func calculateRelevance(query string, element *models.OntologyElement) float64 {
	relevance := 0.0

	nameRelevance := fuzzyMatch(query, element.Name)
	relevance += nameRelevance * 0.6

	typeRelevance := fuzzyMatch(query, element.Type)
	relevance += typeRelevance * 0.3

	descRelevance := fuzzyMatch(query, element.Description)
	relevance += descRelevance * 0.1

	return math.Min(relevance, 1.0)
}

// fuzzyMatch calcule la similarité entre deux chaînes
func fuzzyMatch(s1, s2 string) float64 {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if strings.Contains(s2, s1) {
		return 1.0
	}

	distance := levenshtein.ComputeDistance(s1, s2)
	maxLen := float64(max(len(s1), len(s2)))

	if maxLen == 0 {
		return 0
	}

	similarity := 1 - float64(distance)/maxLen
	if similarity < 0.3 {
		return 0
	}
	return math.Round(similarity*100) / 100 // Arrondir à deux décimales
}

// extractContext extrait le contexte d'un élément
func extractContext(element *models.OntologyElement, contextSize int) string {
	if len(element.Contexts) > 0 {
		context := element.Contexts[0] // Prendre le premier contexte
		before := strings.Join(context.Before, " ")
		after := strings.Join(context.After, " ")
		return fmt.Sprintf("%s [%s] %s", before, context.Element, after)
	}
	// Si pas de contexte disponible, retourner une chaîne vide
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
