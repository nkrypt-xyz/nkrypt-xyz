# API Documentation

This directory contains the automated API documentation system for the nkrypt-xyz web server.

## Quick Start

```bash
# Generate documentation
make generate

# Clean generated docs
make clean
```

## Structure

```
api-docs/
├── docs/               # Generated markdown documentation (commit these!)
│   ├── README.md      # Overview and navigation
│   ├── *-endpoints.md # Endpoint documentation by resource type
│   └── models.md      # Request/response model reference
│
├── generator/          # Documentation generator (Go program)
│   ├── main.go        # Entry point
│   ├── parser/        # AST parsing logic
│   └── writer/        # Markdown generation
│
└── Makefile           # Build commands
```

## How It Works

The generator uses Go's AST parser to:

1. Parse `web-server/internal/router/router.go` to extract all HTTP endpoints
2. Parse `web-server/internal/model/*.go` to understand request/response structures
3. Extract validation rules from struct tags (`validate:"required,min=4,max=32"`)
4. Generate organized Markdown files in `docs/`

**Zero code annotations required** - it works directly from your existing Go code!

## Generated Documentation

The documentation includes:

- **40 endpoints** across 7 resource groups
- Request body schemas with validation constraints
- Authentication requirements
- Response formats
- Complete model reference

View the generated docs: [docs/README.md](./docs/README.md)

## Usage in CI/CD

Add to your build pipeline:

```yaml
# Example GitHub Actions
- name: Generate API docs
  run: |
    cd api-docs
    make generate
    
- name: Check for doc changes
  run: git diff --exit-code api-docs/docs/
```

Or regenerate docs automatically on pre-commit:

```bash
# .git/hooks/pre-commit
#!/bin/bash
cd api-docs && make generate && git add docs/
```

## Maintenance

The generator requires minimal maintenance. When you:

- **Add new endpoints** → Just run `make generate`
- **Change request models** → Just run `make generate`
- **Update validation rules** → Just run `make generate`

The documentation stays in sync with your code automatically!

## Extending

Want to add more features? Edit:

- `generator/parser/parser.go` - To extract more information from code
- `generator/writer/writer.go` - To change markdown formatting
- `generator/parser/types.go` - To add new metadata fields

See [generator/README.md](./generator/README.md) for details.
