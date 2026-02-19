package apperror

// CodedError is the base error type with a machine-readable code.
type CodedError struct {
	Code    string
	Message string
}

func (e *CodedError) Error() string {
	return e.Message
}

// UserError represents client errors (4xx).
type UserError struct {
	CodedError
}

func NewUserError(code, message string) *UserError {
	return &UserError{CodedError{Code: code, Message: message}}
}

// DeveloperError represents server errors (5xx).
type DeveloperError struct {
	CodedError
}

func NewDeveloperError(code, message string) *DeveloperError {
	return &DeveloperError{CodedError{Code: code, Message: message}}
}

// ValidationError wraps go-playground/validator errors.
type ValidationError struct {
	Code    string
	Message string
	Details interface{}
}

func (e *ValidationError) Error() string {
	return e.Message
}

// DetectHTTPStatusCode maps an error to the correct HTTP status code.
func DetectHTTPStatusCode(err error) int {
	switch e := err.(type) {
	case *UserError:
		switch e.Code {
		case "API_KEY_EXPIRED", "API_KEY_NOT_FOUND":
			return 401
		case "ACCESS_DENIED", "USER_BANNED":
			return 403
		case "AUTHORIZATION_HEADER_MISSING", "AUTHORIZATION_HEADER_MALFORMATTED":
			return 412
		default:
			return 400
		}
	case *DeveloperError:
		return 500
	case *ValidationError:
		return 400
	default:
		return 500
	}
}

// SerializeError converts an error to the JSON structure expected by clients.
func SerializeError(err error) map[string]interface{} {
	switch e := err.(type) {
	case *UserError:
		return map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
			"details": map[string]interface{}{},
		}
	case *DeveloperError:
		return map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
			"details": map[string]interface{}{},
		}
	case *ValidationError:
		return map[string]interface{}{
			"code":    "VALIDATION_ERROR",
			"message": e.Message,
			"details": e.Details,
		}
	default:
		return map[string]interface{}{
			"code":    "GENERIC_SERVER_ERROR",
			"message": "We have encountered an unexpected server error. It has been logged and administrators will be notified.",
			"details": map[string]interface{}{},
		}
	}
}

