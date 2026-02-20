package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/api-docs/generator/parser"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/api-docs/generator/writer"
)

func main() {
	// Paths relative to the generator location
	webServerPath := filepath.Join("..", "..", "web-server")
	docsOutputPath := filepath.Join("..", "docs")

	// Resolve absolute paths
	absWebServerPath, err := filepath.Abs(webServerPath)
	if err != nil {
		log.Fatalf("Failed to resolve web-server path: %v", err)
	}

	absDocsPath, err := filepath.Abs(docsOutputPath)
	if err != nil {
		log.Fatalf("Failed to resolve docs output path: %v", err)
	}

	fmt.Println("Parsing API from:", absWebServerPath)
	fmt.Println("Generating docs to:", absDocsPath)

	// Parse the API
	api, err := parser.ParseAPI(absWebServerPath)
	if err != nil {
		log.Fatalf("Failed to parse API: %v", err)
	}

	fmt.Printf("Found %d endpoints\n", len(api.Endpoints))

	// Generate documentation
	if err := writer.GenerateDocs(api, absDocsPath); err != nil {
		log.Fatalf("Failed to generate docs: %v", err)
	}

	fmt.Println("Documentation generated successfully!")
}
