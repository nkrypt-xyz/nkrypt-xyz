//go:build integration

package integration

import (
	"testing"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestHealthz(t *testing.T) {
	resp, err := httpClient.Get(baseURL + "/healthz")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestReadyz(t *testing.T) {
	resp, err := httpClient.Get(baseURL + "/readyz")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMetrics(t *testing.T) {
	resp, err := httpClient.Get(baseURL + "/metrics")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMetricsGetSummary(t *testing.T) {
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/metrics/get-summary", map[string]interface{}{}, adminAPIKey)

	disk, ok := result["disk"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected disk object in response")
	}

	if _, ok := disk["usedBytes"].(float64); !ok {
		t.Error("Expected disk.usedBytes")
	}
	if _, ok := disk["totalBytes"].(float64); !ok {
		t.Error("Expected disk.totalBytes")
	}
}
