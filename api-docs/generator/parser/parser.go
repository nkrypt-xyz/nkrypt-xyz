package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ParseAPI parses the web-server codebase and extracts API information
func ParseAPI(webServerPath string) (*API, error) {
	api := &API{
		Endpoints: []Endpoint{},
		Models:    make(map[string]*Model),
	}

	// Parse router to extract endpoints
	endpoints, err := parseRouter(filepath.Join(webServerPath, "internal", "router", "router.go"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse router: %w", err)
	}
	api.Endpoints = endpoints

	// Parse models to extract request/response structures
	modelsPath := filepath.Join(webServerPath, "internal", "model")
	models, err := parseModels(modelsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}
	api.Models = models

	// Match endpoints with their request models
	matchEndpointModels(api)

	return api, nil
}

// parseRouter parses router.go and extracts endpoint definitions
func parseRouter(routerPath string) ([]Endpoint, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, routerPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var endpoints []Endpoint

	// Walk the AST to find route definitions
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for method calls like r.Post("/api/user/login", userHandler.Login)
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if this is a route registration (Get, Post, etc.)
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		method := selExpr.Sel.Name
		if !isHTTPMethod(method) {
			return true
		}

		// Extract path and handler
		if len(callExpr.Args) >= 2 {
			path := extractStringLiteral(callExpr.Args[0])
			handler := extractHandler(callExpr.Args[1])

			if path != "" && handler != "" {
				endpoint := Endpoint{
					Method:       strings.ToUpper(method),
					Path:         path,
					Handler:      handler,
					RequiresAuth: isAuthRequired(path),
					GroupName:    extractGroupName(path),
				}
				endpoints = append(endpoints, endpoint)
			}
		}

		return true
	})

	return endpoints, nil
}

// parseModels parses the model package and extracts struct definitions
func parseModels(modelsPath string) (map[string]*Model, error) {
	models := make(map[string]*Model)

	err := filepath.Walk(modelsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// Extract struct definitions
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				model := parseStruct(typeSpec.Name.Name, structType, genDecl.Doc)
				models[model.Name] = model
			}
		}

		return nil
	})

	return models, err
}

// parseStruct parses a struct type and extracts field information
func parseStruct(name string, structType *ast.StructType, doc *ast.CommentGroup) *Model {
	model := &Model{
		Name:        name,
		Description: extractComment(doc),
		Fields:      []Field{},
	}

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue // Skip embedded fields
		}

		fieldName := field.Names[0].Name
		fieldType := extractTypeName(field.Type)

		// Parse struct tags
		jsonName, required, constraints := parseStructTag(field.Tag)

		f := Field{
			Name:        fieldName,
			JSONName:    jsonName,
			Type:        fieldType,
			Required:    required,
			Constraints: constraints,
			Description: extractComment(field.Doc),
		}

		model.Fields = append(model.Fields, f)
	}

	return model
}

// parseStructTag parses struct tags and extracts validation information
func parseStructTag(tag *ast.BasicLit) (jsonName string, required bool, constraints string) {
	if tag == nil {
		return "", false, ""
	}

	tagValue := strings.Trim(tag.Value, "`")

	// Extract json tag
	jsonTag := extractTag(tagValue, "json")
	if jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		jsonName = parts[0]
	}

	// Extract validation tag
	validateTag := extractTag(tagValue, "validate")
	if validateTag != "" {
		required = strings.Contains(validateTag, "required")
		constraints = validateTag
	}

	return jsonName, required, constraints
}

// extractTag extracts a specific tag value from a struct tag string
func extractTag(tagStr, tagName string) string {
	parts := strings.Split(tagStr, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, tagName+":") {
			value := strings.TrimPrefix(part, tagName+":")
			value = strings.Trim(value, `"`)
			return value
		}
	}
	return ""
}

// matchEndpointModels matches endpoints with their request and response models
func matchEndpointModels(api *API) {
	for i := range api.Endpoints {
		endpoint := &api.Endpoints[i]

		// Try to match request model (try grouped name first, then simple name)
		for _, name := range inferRequestModelNames(endpoint.Path, endpoint.GroupName) {
			if _, exists := api.Models[name]; exists {
				endpoint.RequestModel = name
				break
			}
		}

		// Try to match response model (try grouped name first, then simple name)
		for _, name := range inferResponseModelNames(endpoint.Path, endpoint.GroupName) {
			if _, exists := api.Models[name]; exists {
				endpoint.ResponseModel = name
				break
			}
		}
	}
}

// inferRequestModelNames returns candidate request model names in order of preference
func inferRequestModelNames(path, group string) []string {
	action := extractActionFromPath(path)
	if action == "" {
		return nil
	}

	simpleName := action + "Request"

	var groupName string
	if len(group) > 0 && group != "system" {
		groupName = strings.ToUpper(group[:1])
		if len(group) > 1 {
			groupName += group[1:]
		}
	}
	groupedName := action + groupName + "Request"

	// For admin/iam/* endpoints, prefer simple name (AddUserRequest, not AddUserAdminRequest)
	if group == "admin" {
		return []string{simpleName, groupedName}
	}
	if groupName != "" && action != "Login" && action != "Assert" && action != "Logout" && action != "LogoutAllSessions" {
		return []string{groupedName, simpleName}
	}
	return []string{simpleName}
}

// inferResponseModelNames returns candidate response model names in order of preference
func inferResponseModelNames(path, group string) []string {
	action := extractActionFromPath(path)
	if action == "" {
		return nil
	}

	simpleName := action + "Response"

	var groupName string
	if len(group) > 0 && group != "system" {
		groupName = strings.ToUpper(group[:1])
		if len(group) > 1 {
			groupName += group[1:]
		}
	}
	groupedName := action + groupName + "Response"

	// For admin/iam/* endpoints, prefer simple name
	if group == "admin" {
		return []string{simpleName, groupedName}
	}
	if groupName != "" && action != "Login" && action != "Assert" && action != "AddUser" {
		return []string{groupedName, simpleName}
	}
	return []string{simpleName}
}

// extractActionFromPath extracts and formats the action from an endpoint path
func extractActionFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	// Get the action (last part of the path)
	lastPart := parts[len(parts)-1]
	
	// Skip path parameters like {bucketId}
	if strings.Contains(lastPart, "{") {
		return ""
	}

	// Convert kebab-case or snake_case to PascalCase
	words := strings.FieldsFunc(lastPart, func(r rune) bool {
		return r == '-' || r == '_'
	})

	var action strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			action.WriteString(strings.ToUpper(word[:1]))
			action.WriteString(word[1:])
		}
	}

	return action.String()
}

// Helper functions

func isHTTPMethod(method string) bool {
	methods := []string{"Get", "Post", "Put", "Delete", "Patch"}
	for _, m := range methods {
		if method == m {
			return true
		}
	}
	return false
}

func extractStringLiteral(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, `"`)
}

func extractHandler(expr ast.Expr) string {
	selExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	return selExpr.Sel.Name
}

func extractTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return extractTypeName(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + extractTypeName(t.Elt)
	case *ast.MapType:
		return "map[" + extractTypeName(t.Key) + "]" + extractTypeName(t.Value)
	case *ast.StarExpr:
		return extractTypeName(t.X)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "unknown"
	}
}

func extractComment(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	var comments []string
	for _, c := range cg.List {
		text := strings.TrimPrefix(c.Text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)
		if text != "" {
			comments = append(comments, text)
		}
	}
	return strings.Join(comments, " ")
}

func isAuthRequired(path string) bool {
	// Public endpoints (no authentication required)
	publicEndpoints := []string{
		"/healthz",
		"/readyz",
		"/metrics",
		"/api/user/login",
		"/user/login", // Both versions (with and without /api prefix)
	}

	for _, public := range publicEndpoints {
		if path == public {
			return false
		}
	}

	// All /api/* and /user/*, /admin/*, etc. endpoints require auth except the ones above
	return strings.HasPrefix(path, "/api/") || 
		strings.HasPrefix(path, "/user/") ||
		strings.HasPrefix(path, "/admin/") ||
		strings.HasPrefix(path, "/bucket/") ||
		strings.HasPrefix(path, "/directory/") ||
		strings.HasPrefix(path, "/file/") ||
		strings.HasPrefix(path, "/blob/") ||
		strings.HasPrefix(path, "/metrics/")
}

func extractGroupName(path string) string {
	// Extract group from path like /user/login -> user or /api/user/login -> user
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	// If path starts with "api", skip it
	if len(parts) > 0 && parts[0] == "api" {
		parts = parts[1:]
	}
	
	// Return first segment as group name
	if len(parts) > 0 {
		firstSegment := parts[0]
		// Health and metrics endpoints are system endpoints
		if firstSegment == "healthz" || firstSegment == "readyz" || firstSegment == "metrics" {
			return "system"
		}
		return firstSegment
	}
	
	return "system"
}
