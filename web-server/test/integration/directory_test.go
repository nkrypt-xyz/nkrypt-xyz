//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestDirectoryCreate(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-dir-create-%d", time.Now().Unix())
	// Create a bucket first
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	// Create a directory
	createDirReq := map[string]interface{}{
		"name":              "test-subdir",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{"type": "folder"},
		"encryptedMetaData": "encrypted-test-data",
	}
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirReq, adminAPIKey)

	dirID, ok := result["directoryId"].(string)
	if !ok || len(dirID) != 16 {
		t.Errorf("Expected 16-char directoryId, got %v", dirID)
	}
}

func TestDirectoryGet(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-dir-get-%d", time.Now().Unix())
	// Create bucket and directory
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createDirReq := map[string]interface{}{
		"name":              "test-dir-get",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirReq, adminAPIKey)
	dirID := createResult["directoryId"].(string)

	// Get directory contents
	getReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getReq, adminAPIKey)

	// Verify structure
	directory, ok := result["directory"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected directory object")
	}
	if directory["_id"].(string) != rootDirID {
		t.Error("Directory ID mismatch")
	}

	childDirs, ok := result["childDirectoryList"].([]interface{})
	if !ok {
		t.Fatal("Expected childDirectoryList array")
	}

	// Should have our created directory
	found := false
	for _, d := range childDirs {
		dir := d.(map[string]interface{})
		if dir["_id"].(string) == dirID {
			found = true
			if name, ok := dir["name"].(string); !ok || name != "test-dir-get" {
				t.Errorf("Expected name='test-dir-get', got %v", name)
			}
			break
		}
	}

	if !found {
		t.Error("Created directory not in childDirectoryList")
	}

	childFiles, ok := result["childFileList"].([]interface{})
	if !ok {
		t.Fatal("Expected childFileList array")
	}
	if len(childFiles) != 0 {
		t.Error("Expected empty childFileList for new directory")
	}
}

func TestDirectoryRename(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-dir-rename-%d", time.Now().Unix())
	// Create bucket and directory
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createDirReq := map[string]interface{}{
		"name":              "old-dir-name",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirReq, adminAPIKey)
	dirID := createResult["directoryId"].(string)

	// Rename it
	renameReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": dirID,
		"name":        "new-dir-name",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/rename", renameReq, adminAPIKey)

	// Verify via get
	getReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	getResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getReq, adminAPIKey)
	childDirs := getResult["childDirectoryList"].([]interface{})

	for _, d := range childDirs {
		dir := d.(map[string]interface{})
		if dir["_id"].(string) == dirID {
			if name, ok := dir["name"].(string); !ok || name != "new-dir-name" {
				t.Errorf("Expected name='new-dir-name', got %v", name)
			}
			return
		}
	}

	t.Error("Renamed directory not found")
}

func TestDirectoryDelete(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-dir-delete-%d", time.Now().Unix())
	// Create bucket and directory
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createDirReq := map[string]interface{}{
		"name":              "dir-to-delete",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirReq, adminAPIKey)
	dirID := createResult["directoryId"].(string)

	// Delete it
	deleteReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": dirID,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/delete", deleteReq, adminAPIKey)

	// Verify it's gone
	getReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	getResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getReq, adminAPIKey)
	childDirs := getResult["childDirectoryList"].([]interface{})

	for _, d := range childDirs {
		dir := d.(map[string]interface{})
		if dir["_id"].(string) == dirID {
			t.Error("Deleted directory still exists")
		}
	}
}

func TestDirectorySetMetadata(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-dir-setmeta-%d", time.Now().Unix())
	// Create bucket and directory
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createDirReq := map[string]interface{}{
		"name":              "dir-for-metadata",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{"version": 1},
		"encryptedMetaData": "initial-encrypted",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirReq, adminAPIKey)
	dirID := createResult["directoryId"].(string)

	// Update metadata
	setMetaReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": dirID,
		"metaData":    map[string]interface{}{"version": 2, "updated": true},
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/set-metadata", setMetaReq, adminAPIKey)

	// Update encrypted metadata separately
	setEncMetaReq := map[string]interface{}{
		"bucketId":          bucketID,
		"directoryId":       dirID,
		"encryptedMetaData": "updated-encrypted",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/set-encrypted-metadata", setEncMetaReq, adminAPIKey)

	// Verify via get
	getReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	getResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getReq, adminAPIKey)
	childDirs := getResult["childDirectoryList"].([]interface{})

	found := false
	for _, d := range childDirs {
		dir := d.(map[string]interface{})
		if dir["_id"].(string) == dirID {
			found = true
			metaData := dir["metaData"].(map[string]interface{})
			if version, ok := metaData["version"].(float64); !ok || version != 2 {
				t.Errorf("Expected version=2, got %v", version)
			}
			if updated, ok := metaData["updated"].(bool); !ok || !updated {
				t.Errorf("Expected updated=true, got %v", updated)
			}
			if encMeta, ok := dir["encryptedMetaData"].(string); !ok || encMeta != "updated-encrypted" {
				t.Errorf("Expected encryptedMetaData='updated-encrypted', got %v", encMeta)
			}
			break
		}
	}

	if !found {
		t.Error("Directory not found after metadata update")
	}
}

func TestDirectoryMove(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-dir-move-%d", time.Now().Unix())
	// Create bucket and directories
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	// Create source directory
	createDirAReq := map[string]interface{}{
		"name":              "dir-to-move",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createAResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirAReq, adminAPIKey)
	dirToMoveID := createAResult["directoryId"].(string)

	// Create target directory
	createDirBReq := map[string]interface{}{
		"name":              "target-dir",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createBResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirBReq, adminAPIKey)
	targetDirID := createBResult["directoryId"].(string)

	// Move directory
	moveReq := map[string]interface{}{
		"bucketId":             bucketID,
		"directoryId":          dirToMoveID,
		"newParentDirectoryId": targetDirID,
		"newName":              "dir-to-move",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/move", moveReq, adminAPIKey)

	// Verify it's no longer in root
	getRootReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	getRootResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getRootReq, adminAPIKey)
	rootChildDirs := getRootResult["childDirectoryList"].([]interface{})

	for _, d := range rootChildDirs {
		dir := d.(map[string]interface{})
		if dir["_id"].(string) == dirToMoveID {
			t.Error("Moved directory still in root")
		}
	}

	// Verify it's now in target directory
	getTargetReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": targetDirID,
	}
	getTargetResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getTargetReq, adminAPIKey)
	targetChildDirs := getTargetResult["childDirectoryList"].([]interface{})

	found := false
	for _, d := range targetChildDirs {
		dir := d.(map[string]interface{})
		if dir["_id"].(string) == dirToMoveID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Moved directory not found in target directory")
	}
}
