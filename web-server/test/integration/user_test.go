//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestUserLogin(t *testing.T) {
	loginReq := map[string]interface{}{
		"userName": "admin",
		"password": "PleaseChangeMe@YourEarliest2Day",
	}

	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")

	// Assert response structure
	if _, ok := result["apiKey"].(string); !ok {
		t.Error("Expected apiKey in response")
	}

	user, ok := result["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user object in response")
	}

	if _, ok := user["_id"].(string); !ok {
		t.Error("Expected user._id")
	}
	if userName, ok := user["userName"].(string); !ok || userName != "admin" {
		t.Errorf("Expected userName=admin, got %v", userName)
	}
	if _, ok := user["displayName"].(string); !ok {
		t.Error("Expected user.displayName")
	}

	session, ok := result["session"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected session object in response")
	}
	if _, ok := session["_id"].(string); !ok {
		t.Error("Expected session._id")
	}
}

func TestUserAssert(t *testing.T) {
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, adminAPIKey)

	// Assert response structure (same as login)
	user, ok := result["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user object in response")
	}

	if userName, ok := user["userName"].(string); !ok || userName != "admin" {
		t.Errorf("Expected userName=admin, got %v", userName)
	}
}

func TestUserLogout(t *testing.T) {
	// First login to get a fresh API key
	loginReq := map[string]interface{}{
		"userName": "admin",
		"password": "PleaseChangeMe@YourEarliest2Day",
	}
	loginResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
	tempAPIKey := loginResult["apiKey"].(string)

	// Logout
	logoutReq := map[string]interface{}{
		"message": "Test logout",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/logout", logoutReq, tempAPIKey)

	// Try to use the old API key - should fail
	_, result, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, tempAPIKey)
	testutil.AssertErrorCode(t, result, "API_KEY_EXPIRED")
}

func TestUserList(t *testing.T) {
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/list", map[string]interface{}{}, adminAPIKey)

	users, ok := result["userList"].([]interface{})
	if !ok {
		t.Fatal("Expected userList array in response")
	}

	if len(users) == 0 {
		t.Error("Expected at least one user (admin)")
	}

	// Check that admin user is in the list
	foundAdmin := false
	for _, u := range users {
		user, ok := u.(map[string]interface{})
		if !ok {
			continue
		}
		if userName, ok := user["userName"].(string); ok && userName == "admin" {
			foundAdmin = true
			break
		}
	}

	if !foundAdmin {
		t.Error("Admin user not found in user list")
	}
}

func TestUserUpdateProfile(t *testing.T) {
	userName := fmt.Sprintf("testprofile%d", time.Now().Unix())
	// Create a test user first
	addUserReq := map[string]interface{}{
		"displayName": "Test User",
		"userName":    userName,
		"password":    "TestPass123!",
	}
	addResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)
	userID := addResult["userId"].(string)

	// Login as the test user
	loginReq := map[string]interface{}{
		"userName": userName,
		"password": "TestPass123!",
	}
	loginResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
	testUserAPIKey := loginResult["apiKey"].(string)

	// Update profile
	updateReq := map[string]interface{}{
		"displayName": "Updated Test User",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/update-profile", updateReq, testUserAPIKey)

	// Verify via assert
	assertResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, testUserAPIKey)
	user := assertResult["user"].(map[string]interface{})
	if displayName, ok := user["displayName"].(string); !ok || displayName != "Updated Test User" {
		t.Errorf("Expected displayName='Updated Test User', got %v", displayName)
	}

	_ = userID // Keep for potential cleanup
}

func TestUserListAllSessions(t *testing.T) {
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/list-all-sessions", map[string]interface{}{}, adminAPIKey)

	sessions, ok := result["sessionList"].([]interface{})
	if !ok {
		t.Fatal("Expected sessionList array in response")
	}

	if len(sessions) == 0 {
		t.Error("Expected at least one session")
	}
}

func TestUserUpdatePassword(t *testing.T) {
	userName := fmt.Sprintf("testpwdchange%d", time.Now().Unix())
	// Create a test user
	addUserReq := map[string]interface{}{
		"displayName": "Password Change Test",
		"userName":    userName,
		"password":    "OldPassword123!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)

	// Login as the test user
	loginReq := map[string]interface{}{
		"userName": userName,
		"password": "OldPassword123!",
	}
	loginResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
	testUserAPIKey := loginResult["apiKey"].(string)

	// Update password
	updatePwdReq := map[string]interface{}{
		"currentPassword": "OldPassword123!",
		"newPassword":     "NewPassword456!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/update-password", updatePwdReq, testUserAPIKey)

	// Old API key should be invalidated
	_, result, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, testUserAPIKey)
	testutil.AssertErrorCode(t, result, "API_KEY_EXPIRED")

	// Should be able to login with new password
	newLoginReq := map[string]interface{}{
		"userName": userName,
		"password": "NewPassword456!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", newLoginReq, "")

	// Old password should not work
	oldLoginReq := map[string]interface{}{
		"userName": userName,
		"password": "OldPassword123!",
	}
	_, oldResult, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/login", oldLoginReq, "")
	testutil.AssertErrorCode(t, oldResult, "PASSWORD_INVALID")
}

func TestUserLogoutAllSessions(t *testing.T) {
	userName := fmt.Sprintf("testlogoutall%d", time.Now().Unix())
	// Create a test user
	addUserReq := map[string]interface{}{
		"displayName": "Logout All Test",
		"userName":    userName,
		"password":    "TestPass123!",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)

	// Login multiple times to create multiple sessions
	loginReq := map[string]interface{}{
		"userName": userName,
		"password": "TestPass123!",
	}
	loginResult1 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
	apiKey1 := loginResult1["apiKey"].(string)

	loginResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/login", loginReq, "")
	apiKey2 := loginResult2["apiKey"].(string)

	// Logout all sessions using the first API key
	logoutAllReq := map[string]interface{}{
		"message": "Logging out all sessions",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/logout-all-sessions", logoutAllReq, apiKey1)

	// Both API keys should now be invalidated
	_, result1, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, apiKey1)
	testutil.AssertErrorCode(t, result1, "API_KEY_EXPIRED")

	_, result2, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/user/assert", map[string]interface{}{}, apiKey2)
	testutil.AssertErrorCode(t, result2, "API_KEY_EXPIRED")
}

func TestUserFind(t *testing.T) {
	userName := fmt.Sprintf("testfind%d", time.Now().Unix())
	// Create a test user
	addUserReq := map[string]interface{}{
		"displayName": "Find Test User",
		"userName":    userName,
		"password":    "TestPass123!",
	}
	addResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/admin/iam/add-user", addUserReq, adminAPIKey)
	userID := addResult["userId"].(string)

	// Find by userName
	findReq := map[string]interface{}{
		"filters": []map[string]interface{}{
			{"by": "userName", "userName": userName},
		},
	}
	findResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/find", findReq, adminAPIKey)

	userList, ok := findResult["userList"].([]interface{})
	if !ok {
		t.Fatal("Expected userList array in response")
	}

	if len(userList) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(userList))
	}

	user := userList[0].(map[string]interface{})
	if user["_id"].(string) != userID {
		t.Errorf("Expected user ID %s, got %s", userID, user["_id"].(string))
	}
	if user["userName"].(string) != userName {
		t.Errorf("Expected userName %s, got %s", userName, user["userName"].(string))
	}

	// Find by userId
	findByIDReq := map[string]interface{}{
		"filters": []map[string]interface{}{
			{"by": "userId", "userId": userID},
		},
	}
	findByIDResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/user/find", findByIDReq, adminAPIKey)

	userListByID, ok := findByIDResult["userList"].([]interface{})
	if !ok || len(userListByID) != 1 {
		t.Fatal("Expected 1 user when finding by userId")
	}

	userByID := userListByID[0].(map[string]interface{})
	if userByID["_id"].(string) != userID {
		t.Errorf("Expected user ID %s, got %s", userID, userByID["_id"].(string))
	}
}
