package main

import (
	"fmt"
	"os"

	"github.com/chrlesur/ontology-server/internal/logger"
	"github.com/chrlesur/ontology-server/internal/storage"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <ontology_file> <context_file> <metadata_file>")
		os.Exit(1)
	}

	ontologyFile := os.Args[1]
	contextFile := os.Args[2]
	metadataFile := os.Args[3]

	log, err := logger.NewLogger(logger.INFO, "logs")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	memStorage := storage.NewMemoryStorage()
	loader := storage.NewOntologyLoader(memStorage, log)

	err = loader.LoadFiles(ontologyFile, contextFile, metadataFile)
	if err != nil {
		fmt.Printf("Failed to load ontology: %v\n", err)
		os.Exit(1)
	}

	ontologies := memStorage.ListOntologies()
	for _, onto := range ontologies {
		fmt.Printf("Loaded ontology: %s\n", onto.Name)
		fmt.Printf("Number of elements: %d\n", len(onto.Elements))
		fmt.Printf("Number of relations: %d\n", len(onto.Relations))

		elementCountWithContexts := 0
		totalContexts := 0
		for _, elem := range onto.Elements {
			if len(elem.Contexts) > 0 {
				elementCountWithContexts++
				totalContexts += len(elem.Contexts)
			}
		}
		fmt.Printf("Elements with contexts: %d\n", elementCountWithContexts)
		fmt.Printf("Total contexts: %d\n", totalContexts)
		fmt.Printf("Average contexts per element: %.2f\n", float64(totalContexts)/float64(len(onto.Elements)))

		// Afficher quelques éléments et leurs contextes pour vérification
		for i, elem := range onto.Elements {
			if i >= 5 {
				break // Limiter à 5 éléments pour l'affichage
			}
			fmt.Printf("Element: %s, Type: %s, Contexts: %d\n", elem.Name, elem.Type, len(elem.Contexts))
			for j, ctx := range elem.Contexts {
				if j >= 2 {
					break // Limiter à 2 contextes par élément
				}
				contextPreview := ctx.Element
				if len(contextPreview) > 20 {
					contextPreview = contextPreview[:20] + "..."
				}
				fmt.Printf("  Context %d: %s (FileID: %s, Position: %d)\n", j, contextPreview, ctx.FileID, ctx.Position)
			}
		}
	}
}
