//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestAdminAddUser(t *testing.T) {
	userName := fmt.Sprintf("newtestuser%d", time.Now().Unix())
	addUserReq := map[string]interface{}{
		"displayName": "New Test User",
		"userName":    userName,
		"password":    "TestPass123!",
	}

	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)

	userID, ok := result["userId"].(string)
	if !ok || len(userID) != 16 {
		t.Errorf("Expected 16-char userId, got %v", userID)
	}

	// Login as new user to verify
	loginReq := map[string]interface{}{
		"userName": userName,
		"password": "TestPass123!",
	}
	loginResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")

	user := loginResult["user"].(map[string]interface{})
	if actualUserName, ok := user["userName"].(string); !ok || actualUserName != userName {
		t.Errorf("Expected userName=%s, got %v", userName, actualUserName)
	}
}

func TestAdminAddUserDuplicate(t *testing.T) {
	userName := fmt.Sprintf("duplicateuser%d", time.Now().Unix())
	// Add user
	addUserReq := map[string]interface{}{
		"displayName": "Duplicate Test",
		"userName":    userName,
		"password":    "TestPass123!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)

	// Try adding again with same username
	_, result, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)
	testutil.AssertErrorCode(t, result, "DUPLICATE_USERNAME")
}

func TestAdminSetGlobalPermissions(t *testing.T) {
	userName := fmt.Sprintf("permtestuser%d", time.Now().Unix())
	// Add a test user
	addUserReq := map[string]interface{}{
		"displayName": "Perm Test User",
		"userName":    userName,
		"password":    "TestPass123!",
	}
	addResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)
	userID := addResult["userId"].(string)

	// Set permissions
	setPermReq := map[string]interface{}{
		"userId": userID,
		"globalPermissions": map[string]bool{
			"CREATE_USER":   true,
			"CREATE_BUCKET": false,
		},
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/set-global-permissions", setPermReq, adminAPIKey)

	// Verify via find
	findReq := map[string]interface{}{
		"filters": []map[string]interface{}{
			{"by": "userId", "userId": userID},
		},
		"includeGlobalPermissions": true,
	}
	findResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/find", findReq, adminAPIKey)

	users := findResult["userList"].([]interface{})
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}

	user := users[0].(map[string]interface{})
	perms := user["globalPermissions"].(map[string]interface{})

	if createUser, ok := perms["CREATE_USER"].(bool); !ok || !createUser {
		t.Error("Expected CREATE_USER=true")
	}
	if createBucket, ok := perms["CREATE_BUCKET"].(bool); !ok || createBucket {
		t.Error("Expected CREATE_BUCKET=false")
	}
}

func TestAdminSetBanningStatus(t *testing.T) {
	userName := fmt.Sprintf("bantestuser%d", time.Now().Unix())
	// Add a test user
	addUserReq := map[string]interface{}{
		"displayName": "Ban Test User",
		"userName":    userName,
		"password":    "TestPass123!",
	}
	addResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)
	userID := addResult["userId"].(string)

	// Ban the user
	banReq := map[string]interface{}{
		"userId":   userID,
		"isBanned": true,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/set-banning-status", banReq, adminAPIKey)

	// Try to login as banned user - should fail
	loginReq := map[string]interface{}{
		"userName": userName,
		"password": "TestPass123!",
	}
	_, result, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/login", loginReq, "")
	testutil.AssertErrorCode(t, result, "USER_BANNED")

	// Unban the user
	unbanReq := map[string]interface{}{
		"userId":   userID,
		"isBanned": false,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/set-banning-status", unbanReq, adminAPIKey)

	// Should be able to login now
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
}

func TestAdminOverwriteUserPassword(t *testing.T) {
	userName := fmt.Sprintf("pwdtestuser%d", time.Now().Unix())
	// Add a test user
	addUserReq := map[string]interface{}{
		"displayName": "Password Test User",
		"userName":    userName,
		"password":    "OldPass123!",
	}
	addResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)
	userID := addResult["userId"].(string)

	// Login to get API key
	loginReq := map[string]interface{}{
		"userName": userName,
		"password": "OldPass123!",
	}
	loginResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
	oldAPIKey := loginResult["apiKey"].(string)

	// Overwrite password
	overwriteReq := map[string]interface{}{
		"userId":      userID,
		"newPassword": "NewPass456!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/overwrite-user-password", overwriteReq, adminAPIKey)

	// Old API key should be expired
	_, result, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, oldAPIKey)
	testutil.AssertErrorCode(t, result, "API_KEY_EXPIRED")

	// Should be able to login with new password
	newLoginReq := map[string]interface{}{
		"userName": userName,
		"password": "NewPass456!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", newLoginReq, "")
}
