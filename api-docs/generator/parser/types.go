package parser

// API represents the complete API documentation structure
type API struct {
	Endpoints []Endpoint
	Models    map[string]*Model
}

// Endpoint represents a single API endpoint
type Endpoint struct {
	Method            string
	Path              string
	Handler           string
	Description       string
	RequiresAuth      bool
	RequestModel      string
	ResponseModel     string
	ErrorResponses    []ErrorResponse
	GroupName         string // e.g., "user", "admin", "bucket"
}

// Model represents a request or response model
type Model struct {
	Name        string
	Description string
	Fields      []Field
}

// Field represents a struct field in a model
type Field struct {
	Name        string
	JSONName    string
	Type        string
	Required    bool
	Constraints string
	Description string
}

// ErrorResponse represents a possible error response
type ErrorResponse struct {
	Code        string
	Description string
}
