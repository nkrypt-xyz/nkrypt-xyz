//go:build integration

package integration

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

var (
	baseURL     string
	httpClient  *http.Client
	adminAPIKey string
)

func TestMain(m *testing.M) {
	baseURL = os.Getenv("NK_TEST_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:9041"
	}

	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	// Login as admin to get API key for authenticated tests
	adminAPIKey = loginAsAdmin()

	os.Exit(m.Run())
}

func loginAsAdmin() string {
	loginReq := map[string]interface{}{
		"userName": "admin",
		"password": "PleaseChangeMe@YourEarliest2Day",
	}

	_, result, err := testutil.CallPostJSON(httpClient, baseURL+"/api/user/login", loginReq, "")
	if err != nil {
		panic("Failed to login as admin: " + err.Error())
	}

	apiKey, ok := result["apiKey"].(string)
	if !ok {
		panic("No apiKey in login response")
	}

	return apiKey
}
