package apperror

import (
	"testing"
)

func TestUserErrorHTTPStatus(t *testing.T) {
	tests := []struct {
		code           string
		expectedStatus int
	}{
		{"API_KEY_EXPIRED", 401},
		{"API_KEY_NOT_FOUND", 401},
		{"ACCESS_DENIED", 403},
		{"USER_BANNED", 403},
		{"AUTHORIZATION_HEADER_MISSING", 412},
		{"AUTHORIZATION_HEADER_MALFORMATTED", 412},
		{"VALIDATION_ERROR", 400},
		{"DUPLICATE_BUCKET_NAME", 400},
		{"UNKNOWN_ERROR", 400}, // Default for UserError
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			err := NewUserError(tt.code, "test message")
			status := DetectHTTPStatusCode(err)
			if status != tt.expectedStatus {
				t.Errorf("Expected status %d for code %s, got %d", tt.expectedStatus, tt.code, status)
			}
		})
	}
}

func TestDeveloperErrorHTTPStatus(t *testing.T) {
	err := NewDeveloperError("INTERNAL_ERROR", "test message")
	status := DetectHTTPStatusCode(err)
	if status != 500 {
		t.Errorf("Expected status 500 for DeveloperError, got %d", status)
	}
}

func TestValidationErrorHTTPStatus(t *testing.T) {
	err := &ValidationError{Code: "VALIDATION_ERROR", Message: "test"}
	status := DetectHTTPStatusCode(err)
	if status != 400 {
		t.Errorf("Expected status 400 for ValidationError, got %d", status)
	}
}

func TestSerializeError(t *testing.T) {
	err := NewUserError("TEST_CODE", "Test message")
	serialized := SerializeError(err)

	if serialized.Code != "TEST_CODE" {
		t.Errorf("Expected code TEST_CODE, got %v", serialized.Code)
	}
	if serialized.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %v", serialized.Message)
	}
}
