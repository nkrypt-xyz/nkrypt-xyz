package redact

import (
	"encoding/json"
	"strings"
)

// SensitiveFields is a list of field names that should be redacted from logs.
var SensitiveFields = []string{
	"password",
	"currentPassword",
	"newPassword",
	"apiKey",
	"api_key",
	"nk-api-key",
	"authorization",
	"token",
	"secret",
	"secretKey",
	"secret_key",
	"accessKey",
	"access_key",
	"cryptData",
	"crypt_data",
	"encryptedMetaData",
	"encrypted_meta_data",
}

// RedactValue returns a redacted string for sensitive data.
func RedactValue() string {
	return "[REDACTED]"
}

// RedactJSON recursively redacts sensitive fields from a JSON object.
// Returns the redacted JSON as a map.
func RedactJSON(data map[string]interface{}) map[string]interface{} {
	redacted := make(map[string]interface{})
	
	for key, value := range data {
		if isSensitiveField(key) {
			redacted[key] = RedactValue()
			continue
		}
		
		// Recursively redact nested objects
		switch v := value.(type) {
		case map[string]interface{}:
			redacted[key] = RedactJSON(v)
		case []interface{}:
			redacted[key] = redactArray(v)
		default:
			redacted[key] = value
		}
	}
	
	return redacted
}

// redactArray recursively redacts sensitive fields from arrays.
func redactArray(arr []interface{}) []interface{} {
	redacted := make([]interface{}, len(arr))
	
	for i, item := range arr {
		switch v := item.(type) {
		case map[string]interface{}:
			redacted[i] = RedactJSON(v)
		case []interface{}:
			redacted[i] = redactArray(v)
		default:
			redacted[i] = item
		}
	}
	
	return redacted
}

// RedactJSONString takes a JSON string, parses it, redacts sensitive fields,
// and returns the redacted JSON string.
func RedactJSONString(jsonStr string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// If parsing fails, return a generic redacted message
		return "[INVALID_JSON]"
	}
	
	redacted := RedactJSON(data)
	redactedJSON, err := json.Marshal(redacted)
	if err != nil {
		return "[ERROR_REDACTING]"
	}
	
	return string(redactedJSON)
}

// RedactHeaders redacts sensitive HTTP headers.
func RedactHeaders(headers map[string][]string) map[string][]string {
	redacted := make(map[string][]string)
	
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if isSensitiveField(lowerKey) {
			redacted[key] = []string{RedactValue()}
		} else {
			redacted[key] = values
		}
	}
	
	return redacted
}

// isSensitiveField checks if a field name is in the sensitive fields list.
func isSensitiveField(fieldName string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, sensitive := range SensitiveFields {
		if strings.ToLower(sensitive) == lowerField {
			return true
		}
	}
	return false
}
