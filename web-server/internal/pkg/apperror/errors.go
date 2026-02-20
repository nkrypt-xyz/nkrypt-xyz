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

// SerializedError is the JSON structure for error responses.
type SerializedError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

// SerializeError converts an error to the JSON structure expected by clients.
func SerializeError(err error) SerializedError {
	switch e := err.(type) {
	case *UserError:
		return SerializedError{
			Code:    e.Code,
			Message: e.Message,
			Details: map[string]interface{}{},
		}
	case *DeveloperError:
		return SerializedError{
			Code:    e.Code,
			Message: e.Message,
			Details: map[string]interface{}{},
		}
	case *ValidationError:
		return SerializedError{
			Code:    "VALIDATION_ERROR",
			Message: e.Message,
			Details: e.Details,
		}
	default:
		return SerializedError{
			Code:    "GENERIC_SERVER_ERROR",
			Message: "We have encountered an unexpected server error. It has been logged and administrators will be notified.",
			Details: map[string]interface{}{},
		}
	}
}

