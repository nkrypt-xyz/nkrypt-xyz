package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// CallPostJSON sends a POST request with JSON body and returns the response and parsed body
func CallPostJSON(client *http.Client, url string, body interface{}, apiKey string) (*http.Response, map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return resp, nil, err
	}
	resp.Body.Close()

	return resp, result, nil
}

// CallPostJSONExpectSuccess sends POST and asserts hasError=false, status 200
func CallPostJSONExpectSuccess(t *testing.T, client *http.Client, url string, body interface{}, apiKey string) map[string]interface{} {
	t.Helper()
	resp, result, err := CallPostJSON(client, url, body, apiKey)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	hasError, ok := result["hasError"].(bool)
	if !ok || hasError {
		t.Fatalf("Expected hasError=false, got %v. Full response: %+v", hasError, result)
	}

	return result
}

// CallPostRaw sends a POST request with raw body (for blob uploads)
func CallPostRaw(client *http.Client, url string, body io.Reader, headers map[string]string, apiKey string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	return client.Do(req)
}

// AssertErrorCode checks the response has hasError=true with specific error code
func AssertErrorCode(t *testing.T, result map[string]interface{}, expectedCode string) {
	t.Helper()
	hasError, ok := result["hasError"].(bool)
	if !ok || !hasError {
		t.Fatalf("Expected hasError=true, got %v", hasError)
	}

	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error object, got %v", result["error"])
	}

	code, ok := errorObj["code"].(string)
	if !ok {
		t.Fatalf("Expected error code string, got %v", errorObj["code"])
	}

	if code != expectedCode {
		t.Fatalf("Expected error code %s, got %s", expectedCode, code)
	}
}

// ParseJSONResponse parses JSON from http.Response into the provided interface
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	return json.NewDecoder(resp.Body).Decode(v)
}
