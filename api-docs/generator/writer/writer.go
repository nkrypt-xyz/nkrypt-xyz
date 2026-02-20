package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/api-docs/generator/parser"
)

// GenerateDocs generates markdown documentation files
func GenerateDocs(api *parser.API, outputPath string) error {
	// Create output directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Group endpoints by category
	groups := groupEndpoints(api.Endpoints)

	// Generate README.md (overview)
	if err := generateReadme(api, groups, outputPath); err != nil {
		return err
	}

	// Generate individual endpoint group files
	for groupName, endpoints := range groups {
		if err := generateGroupFile(groupName, endpoints, api.Models, outputPath); err != nil {
			return err
		}
	}

	// Generate models reference
	if err := generateModelsReference(api.Models, outputPath); err != nil {
		return err
	}

	return nil
}

// groupEndpoints groups endpoints by their group name
func groupEndpoints(endpoints []parser.Endpoint) map[string][]parser.Endpoint {
	groups := make(map[string][]parser.Endpoint)

	for _, endpoint := range endpoints {
		groupName := endpoint.GroupName
		groups[groupName] = append(groups[groupName], endpoint)
	}

	return groups
}

// generateReadme generates the main README.md file
func generateReadme(api *parser.API, groups map[string][]parser.Endpoint, outputPath string) error {
	var sb strings.Builder

	sb.WriteString("# API Documentation\n\n")
	sb.WriteString("This documentation is automatically generated from the Go source code.\n\n")

	// API Overview
	sb.WriteString("## Overview\n\n")
	sb.WriteString(fmt.Sprintf("Total Endpoints: **%d**\n\n", len(api.Endpoints)))

	// Authentication
	sb.WriteString("## Authentication\n\n")
	sb.WriteString("Most endpoints require authentication using an API key.\n\n")
	sb.WriteString("Include the API key in the request header:\n")
	sb.WriteString("```\n")
	sb.WriteString("Authorization: Bearer <your-api-key>\n")
	sb.WriteString("```\n\n")

	// Endpoint Groups
	sb.WriteString("## Endpoint Groups\n\n")

	// Sort groups alphabetically
	var groupNames []string
	for name := range groups {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		endpoints := groups[groupName]
		sb.WriteString(fmt.Sprintf("- [%s](./%s-endpoints.md) - %d endpoints\n",
			strings.Title(groupName), groupName, len(endpoints)))
	}

	sb.WriteString("\n## Models Reference\n\n")
	sb.WriteString("- [Models Reference](./models.md) - Request and response model definitions\n")

	// Write to file
	return os.WriteFile(filepath.Join(outputPath, "README.md"), []byte(sb.String()), 0644)
}

// generateGroupFile generates a markdown file for a group of endpoints
func generateGroupFile(groupName string, endpoints []parser.Endpoint, models map[string]*parser.Model, outputPath string) error {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# %s Endpoints\n\n", strings.Title(groupName)))
	sb.WriteString(fmt.Sprintf("This page documents all **%s** related endpoints.\n\n", groupName))

	// Table of contents
	sb.WriteString("## Table of Contents\n\n")
	for _, endpoint := range endpoints {
		anchor := generateAnchor(endpoint.Method, endpoint.Path)
		sb.WriteString(fmt.Sprintf("- [%s %s](#%s)\n", endpoint.Method, endpoint.Path, anchor))
	}
	sb.WriteString("\n---\n\n")

	// Endpoint details
	for _, endpoint := range endpoints {
		writeEndpoint(&sb, endpoint, models)
		sb.WriteString("\n---\n\n")
	}

	// Write to file
	filename := fmt.Sprintf("%s-endpoints.md", groupName)
	return os.WriteFile(filepath.Join(outputPath, filename), []byte(sb.String()), 0644)
}

// writeEndpoint writes a single endpoint documentation
func writeEndpoint(sb *strings.Builder, endpoint parser.Endpoint, models map[string]*parser.Model) {
	// Header
	anchor := generateAnchor(endpoint.Method, endpoint.Path)
	sb.WriteString(fmt.Sprintf("## %s %s {#%s}\n\n", endpoint.Method, endpoint.Path, anchor))

	// Description
	if endpoint.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", endpoint.Description))
	}

	// Authentication
	if endpoint.RequiresAuth {
		sb.WriteString("🔒 **Authentication Required**\n\n")
	} else {
		sb.WriteString("🌐 **Public Endpoint** (No authentication required)\n\n")
	}

	// Request Body
	if endpoint.RequestModel != "" {
		sb.WriteString("### Request Body\n\n")
		if model, exists := models[endpoint.RequestModel]; exists {
			writeModelTable(sb, model)
		} else {
			sb.WriteString(fmt.Sprintf("Model: `%s`\n\n", endpoint.RequestModel))
		}
	} else {
		sb.WriteString("### Request Body\n\n")
		sb.WriteString("No request body required.\n\n")
	}

	// Response
	sb.WriteString("### Response\n\n")
	
	if endpoint.ResponseModel != "" {
		sb.WriteString("**Success (200):**\n\n")
		if model, exists := models[endpoint.ResponseModel]; exists {
			sb.WriteString(fmt.Sprintf("Response Model: [`%s`](./models.md#%s)\n\n", 
				endpoint.ResponseModel, 
				strings.ToLower(endpoint.ResponseModel)))
			writeModelTable(sb, model)
		} else {
			sb.WriteString(fmt.Sprintf("Model: `%s` - See [Models Reference](./models.md#%s)\n\n", 
				endpoint.ResponseModel,
				strings.ToLower(endpoint.ResponseModel)))
		}
	} else {
		sb.WriteString("**Success (200):**\n\n")
		sb.WriteString("```json\n")
		sb.WriteString("{\n")
		sb.WriteString("  \"hasError\": false,\n")
		sb.WriteString("  ...\n")
		sb.WriteString("}\n")
		sb.WriteString("```\n\n")
	}

	// Error responses
	sb.WriteString("**Error Responses:**\n\n")
	sb.WriteString("Common error codes:\n")
	sb.WriteString("- `ACCESS_DENIED` - Authentication required or insufficient permissions\n")
	sb.WriteString("- `VALIDATION_ERROR` - Request validation failed\n")
	sb.WriteString("- `NOT_FOUND` - Resource not found\n\n")
}

// writeModelTable writes a table for a model's fields
func writeModelTable(sb *strings.Builder, model *parser.Model) {
	sb.WriteString("| Field | Type | Required | Constraints | Description |\n")
	sb.WriteString("|-------|------|----------|-------------|-------------|\n")

	for _, field := range model.Fields {
		required := "No"
		if field.Required {
			required = "**Yes**"
		}

		jsonName := field.JSONName
		if jsonName == "" {
			jsonName = field.Name
		}

		// Format constraints
		constraints := formatConstraints(field.Constraints)

		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s |\n",
			jsonName,
			field.Type,
			required,
			constraints,
			field.Description,
		))
	}

	sb.WriteString("\n")
}

// generateModelsReference generates a reference file for all models
func generateModelsReference(models map[string]*parser.Model, outputPath string) error {
	var sb strings.Builder

	sb.WriteString("# Models Reference\n\n")
	sb.WriteString("This page documents all request and response models used in the API.\n\n")

	// Sort models alphabetically
	var modelNames []string
	for name := range models {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	// Table of contents
	sb.WriteString("## Table of Contents\n\n")
	for _, name := range modelNames {
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") {
			sb.WriteString(fmt.Sprintf("- [%s](#%s)\n", name, strings.ToLower(name)))
		}
	}
	sb.WriteString("\n---\n\n")

	// Model details
	for _, name := range modelNames {
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") {
			model := models[name]
			sb.WriteString(fmt.Sprintf("## %s\n\n", name))

			if model.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", model.Description))
			}

			writeModelTable(&sb, model)
			sb.WriteString("\n---\n\n")
		}
	}

	return os.WriteFile(filepath.Join(outputPath, "models.md"), []byte(sb.String()), 0644)
}

// Helper functions

func generateAnchor(method, path string) string {
	anchor := strings.ToLower(method + "-" + path)
	anchor = strings.ReplaceAll(anchor, "/", "-")
	anchor = strings.ReplaceAll(anchor, "{", "")
	anchor = strings.ReplaceAll(anchor, "}", "")
	anchor = strings.Trim(anchor, "-")
	return anchor
}

func formatConstraints(constraints string) string {
	if constraints == "" {
		return "-"
	}

	// Parse validation constraints and make them human-readable
	parts := strings.Split(constraints, ",")
	var formatted []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "required" {
			continue // Already shown in Required column
		}

		// Format common constraints
		if strings.HasPrefix(part, "min=") {
			formatted = append(formatted, "Min: "+strings.TrimPrefix(part, "min="))
		} else if strings.HasPrefix(part, "max=") {
			formatted = append(formatted, "Max: "+strings.TrimPrefix(part, "max="))
		} else if strings.HasPrefix(part, "len=") {
			formatted = append(formatted, "Length: "+strings.TrimPrefix(part, "len="))
		} else if strings.HasPrefix(part, "oneof=") {
			formatted = append(formatted, "One of: "+strings.TrimPrefix(part, "oneof="))
		} else {
			formatted = append(formatted, part)
		}
	}

	if len(formatted) == 0 {
		return "-"
	}

	return strings.Join(formatted, ", ")
}
