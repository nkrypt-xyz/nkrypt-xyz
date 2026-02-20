# API Documentation Generator

This tool automatically generates Markdown documentation for the nkrypt-xyz web server API by parsing the Go source code.

## How It Works

The generator uses Go's AST (Abstract Syntax Tree) parser to:

1. Parse `router.go` to extract all endpoint definitions
2. Parse model files to understand request/response structures
3. Extract validation rules from struct tags
4. Generate organized Markdown files

## Usage

From the `api-docs/generator` directory:

```bash
# Install dependencies
go mod download

# Generate documentation
go run main.go
```

The generated documentation will be written to `api-docs/docs/`.

## Project Structure

```
api-docs/
├── generator/          # Documentation generator code
│   ├── main.go        # Entry point
│   ├── parser/        # AST parsing logic
│   └── writer/        # Markdown generation
└── docs/              # Generated documentation (output)
    ├── README.md      # Overview
    ├── *-endpoints.md # Endpoint groups
    └── models.md      # Model reference
```

## Adding to Build Process

To automatically regenerate docs, add a Makefile or script:

```makefile
.PHONY: docs
docs:
	cd api-docs/generator && go run main.go
```

## Extending the Generator

- **Add descriptions:** Add Go doc comments above handlers and models
- **Customize output:** Modify `writer/writer.go` templates
- **Add more metadata:** Extend `parser/types.go` with additional fields
