.PHONY: generate clean help

# Generate API documentation
generate:
	@echo "🔍 Generating API documentation..."
	@cd generator && go run main.go
	@echo "✨ Done! Documentation is in api-docs/docs/"

# Clean generated documentation
clean:
	@echo "🧹 Cleaning generated documentation..."
	@rm -rf docs/*.md
	@echo "✨ Cleaned!"

# Show help
help:
	@echo "API Documentation Generator"
	@echo ""
	@echo "Usage:"
	@echo "  make generate    Generate API documentation from Go source code"
	@echo "  make clean       Remove all generated documentation files"
	@echo "  make help        Show this help message"
