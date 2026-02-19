//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestBucketCreate(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-create-%d", time.Now().Unix())
	createReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{"description": "Test bucket"},
	}

	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createReq, adminAPIKey)

	bucketID, ok := result["bucketId"].(string)
	if !ok || len(bucketID) != 16 {
		t.Errorf("Expected 16-char bucketId, got %v", bucketID)
	}

	rootDirectoryID, ok := result["rootDirectoryId"].(string)
	if !ok || len(rootDirectoryID) != 16 {
		t.Errorf("Expected 16-char rootDirectoryId, got %v", rootDirectoryID)
	}
}

func TestBucketList(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-list-%d", time.Now().Unix())
	// Create a bucket first
	createReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createReq, adminAPIKey)
	bucketID := createResult["bucketId"].(string)

	// List buckets
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/list", map[string]interface{}{}, adminAPIKey)

	bucketList, ok := result["bucketList"].([]interface{})
	if !ok {
		t.Fatal("Expected bucketList array")
	}

	// Find our bucket
	found := false
	for _, b := range bucketList {
		bucket := b.(map[string]interface{})
		if bucket["_id"].(string) == bucketID {
			found = true

			// Verify structure
			if _, ok := bucket["name"].(string); !ok {
				t.Error("Expected bucket.name")
			}
			if _, ok := bucket["rootDirectoryId"].(string); !ok {
				t.Error("Expected bucket.rootDirectoryId")
			}
			if _, ok := bucket["cryptSpec"].(string); !ok {
				t.Error("Expected bucket.cryptSpec")
			}
			if _, ok := bucket["bucketAuthorizations"].([]interface{}); !ok {
				t.Error("Expected bucket.bucketAuthorizations array")
			}
			if _, ok := bucket["createdByUserIdentifier"].(string); !ok {
				t.Error("Expected bucket.createdByUserIdentifier")
			}
			if _, ok := bucket["createdAt"].(float64); !ok {
				t.Error("Expected bucket.createdAt as number (epoch ms)")
			}

			break
		}
	}

	if !found {
		t.Error("Created bucket not found in list")
	}
}

func TestBucketRename(t *testing.T) {
	timestamp := time.Now().Unix()
	oldName := fmt.Sprintf("test-bucket-rename-old-%d", timestamp)
	newName := fmt.Sprintf("test-bucket-rename-new-%d", timestamp)
	
	// Create a bucket
	createReq := map[string]interface{}{
		"name":      oldName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createReq, adminAPIKey)
	bucketID := createResult["bucketId"].(string)

	// Rename it
	renameReq := map[string]interface{}{
		"bucketId": bucketID,
		"name":     newName,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/rename", renameReq, adminAPIKey)

	// Verify via list
	listResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/list", map[string]interface{}{}, adminAPIKey)
	bucketList := listResult["bucketList"].([]interface{})

	for _, b := range bucketList {
		bucket := b.(map[string]interface{})
		if bucket["_id"].(string) == bucketID {
			if name, ok := bucket["name"].(string); !ok || name != newName {
				t.Errorf("Expected name='%s', got %v", newName, name)
			}
			return
		}
	}

	t.Error("Renamed bucket not found")
}

func TestBucketSetMetadata(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-metadata-%d", time.Now().Unix())
	// Create a bucket
	createReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{"version": 1},
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createReq, adminAPIKey)
	bucketID := createResult["bucketId"].(string)

	// Update metadata
	setMetaReq := map[string]interface{}{
		"bucketId": bucketID,
		"metaData": map[string]interface{}{"version": 2, "updated": true},
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/set-metadata", setMetaReq, adminAPIKey)

	// Verify via list
	listResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/list", map[string]interface{}{}, adminAPIKey)
	bucketList := listResult["bucketList"].([]interface{})

	for _, b := range bucketList {
		bucket := b.(map[string]interface{})
		if bucket["_id"].(string) == bucketID {
			metaData := bucket["metaData"].(map[string]interface{})
			if version, ok := metaData["version"].(float64); !ok || version != 2 {
				t.Errorf("Expected version=2, got %v", version)
			}
			if updated, ok := metaData["updated"].(bool); !ok || !updated {
				t.Errorf("Expected updated=true, got %v", updated)
			}
			return
		}
	}

	t.Error("Bucket not found")
}

func TestBucketDestroy(t *testing.T) {
	// Create a bucket
	createReq := map[string]interface{}{
		"name":      "test-bucket-destroy",
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createReq, adminAPIKey)
	bucketID := createResult["bucketId"].(string)

	// Destroy it
	destroyReq := map[string]interface{}{
		"bucketId": bucketID,
		"name":     "test-bucket-destroy",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/destroy", destroyReq, adminAPIKey)

	// Verify it's gone from list
	listResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/list", map[string]interface{}{}, adminAPIKey)
	bucketList := listResult["bucketList"].([]interface{})

	for _, b := range bucketList {
		bucket := b.(map[string]interface{})
		if bucket["_id"].(string) == bucketID {
			t.Error("Destroyed bucket still in list")
		}
	}
}
