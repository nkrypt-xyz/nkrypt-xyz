package redact

import (
	"encoding/json"
	"testing"
)

func TestRedactJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "redact password field",
			input: map[string]interface{}{
				"userName": "testuser",
				"password": "supersecret123",
			},
			expected: map[string]interface{}{
				"userName": "testuser",
				"password": "[REDACTED]",
			},
		},
		{
			name: "redact multiple sensitive fields",
			input: map[string]interface{}{
				"userName":        "testuser",
				"password":        "supersecret123",
				"apiKey":          "abc123xyz",
				"currentPassword": "oldpass",
				"newPassword":     "newpass",
			},
			expected: map[string]interface{}{
				"userName":        "testuser",
				"password":        "[REDACTED]",
				"apiKey":          "[REDACTED]",
				"currentPassword": "[REDACTED]",
				"newPassword":     "[REDACTED]",
			},
		},
		{
			name: "redact nested objects",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "testuser",
					"password": "secret",
				},
				"config": map[string]interface{}{
					"apiKey": "key123",
					"debug":  true,
				},
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "testuser",
					"password": "[REDACTED]",
				},
				"config": map[string]interface{}{
					"apiKey": "[REDACTED]",
					"debug":  true,
				},
			},
		},
		{
			name: "redact arrays with objects",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name":     "user1",
						"password": "pass1",
					},
					map[string]interface{}{
						"name":     "user2",
						"password": "pass2",
					},
				},
			},
			expected: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name":     "user1",
						"password": "[REDACTED]",
					},
					map[string]interface{}{
						"name":     "user2",
						"password": "[REDACTED]",
					},
				},
			},
		},
		{
			name: "preserve non-sensitive fields",
			input: map[string]interface{}{
				"bucketId":   "abc123",
				"bucketName": "my-bucket",
				"metaData": map[string]interface{}{
					"createdAt": "2024-01-01",
				},
			},
			expected: map[string]interface{}{
				"bucketId":   "abc123",
				"bucketName": "my-bucket",
				"metaData": map[string]interface{}{
					"createdAt": "2024-01-01",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactJSON(tt.input)
			resultJSON, _ := json.Marshal(result)
			expectedJSON, _ := json.Marshal(tt.expected)

			if string(resultJSON) != string(expectedJSON) {
				t.Errorf("RedactJSON() = %s, want %s", resultJSON, expectedJSON)
			}
		})
	}
}

func TestRedactJSONString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "redact password in JSON string",
			input:    `{"userName":"testuser","password":"supersecret"}`,
			expected: `{"password":"[REDACTED]","userName":"testuser"}`,
		},
		{
			name:     "invalid JSON returns error marker",
			input:    `{invalid json`,
			expected: `[INVALID_JSON]`,
		},
		{
			name:     "redact cryptData field",
			input:    `{"name":"bucket1","cryptData":"encrypted-data-here"}`,
			expected: `{"cryptData":"[REDACTED]","name":"bucket1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactJSONString(tt.input)
			
			// Parse both to compare as JSON objects (field order doesn't matter)
			var resultObj, expectedObj interface{}
			if err := json.Unmarshal([]byte(result), &resultObj); err != nil {
				if result != tt.expected {
					t.Errorf("RedactJSONString() = %s, want %s", result, tt.expected)
				}
				return
			}
			if err := json.Unmarshal([]byte(tt.expected), &expectedObj); err != nil {
				t.Errorf("Test setup error: invalid expected JSON: %s", tt.expected)
				return
			}
			
			resultJSON, _ := json.Marshal(resultObj)
			expectedJSON, _ := json.Marshal(expectedObj)
			
			if string(resultJSON) != string(expectedJSON) {
				t.Errorf("RedactJSONString() = %s, want %s", resultJSON, expectedJSON)
			}
		})
	}
}

func TestRedactHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string][]string
		expected map[string][]string
	}{
		{
			name: "redact authorization header",
			input: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer token123"},
			},
			expected: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {"[REDACTED]"},
			},
		},
		{
			name: "redact custom API key header",
			input: map[string][]string{
				"Content-Type": {"application/json"},
				"Nk-Api-Key":   {"secret-key-123"},
			},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
				"Nk-Api-Key":   {"[REDACTED]"},
			},
		},
		{
			name: "preserve non-sensitive headers",
			input: map[string][]string{
				"Content-Type": {"application/json"},
				"User-Agent":   {"Mozilla/5.0"},
				"Accept":       {"*/*"},
			},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
				"User-Agent":   {"Mozilla/5.0"},
				"Accept":       {"*/*"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactHeaders(tt.input)
			resultJSON, _ := json.Marshal(result)
			expectedJSON, _ := json.Marshal(tt.expected)

			if string(resultJSON) != string(expectedJSON) {
				t.Errorf("RedactHeaders() = %s, want %s", resultJSON, expectedJSON)
			}
		})
	}
}
