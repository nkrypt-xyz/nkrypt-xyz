//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestFileCreate(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-create-%d", time.Now().Unix())
	// Create bucket
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	// Create a file
	createFileReq := map[string]interface{}{
		"name":              "test-file.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{"mimeType": "text/plain"},
		"encryptedMetaData": "encrypted-file-metadata",
	}
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)

	fileID, ok := result["fileId"].(string)
	if !ok || len(fileID) != 16 {
		t.Errorf("Expected 16-char fileId, got %v", fileID)
	}
}

func TestFileGet(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-get-%d", time.Now().Unix())
	// Create bucket and file
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createFileReq := map[string]interface{}{
		"name":              "test-file-get.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{"size": 1024},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)
	fileID := createResult["fileId"].(string)

	// Get file
	getReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
	}
	result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/get", getReq, adminAPIKey)

	file, ok := result["file"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected file object")
	}

	if file["_id"].(string) != fileID {
		t.Error("File ID mismatch")
	}
	if name, ok := file["name"].(string); !ok || name != "test-file-get.txt" {
		t.Errorf("Expected name='test-file-get.txt', got %v", name)
	}

	metaData := file["metaData"].(map[string]interface{})
	if size, ok := metaData["size"].(float64); !ok || size != 1024 {
		t.Errorf("Expected metaData.size=1024, got %v", size)
	}
}

func TestFileRename(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-rename-%d", time.Now().Unix())
	// Create bucket and file
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createFileReq := map[string]interface{}{
		"name":              "old-filename.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)
	fileID := createResult["fileId"].(string)

	// Rename it
	renameReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
		"name":     "new-filename.txt",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/rename", renameReq, adminAPIKey)

	// Verify via get
	getReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
	}
	getResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/get", getReq, adminAPIKey)
	file := getResult["file"].(map[string]interface{})

	if name, ok := file["name"].(string); !ok || name != "new-filename.txt" {
		t.Errorf("Expected name='new-filename.txt', got %v", name)
	}
}

func TestFileDelete(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-delete-%d", time.Now().Unix())
	// Create bucket and file
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createFileReq := map[string]interface{}{
		"name":              "file-to-delete.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)
	fileID := createResult["fileId"].(string)

	// Delete it
	deleteReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/delete", deleteReq, adminAPIKey)

	// Verify it's gone
	getDirReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	getDirResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getDirReq, adminAPIKey)
	childFiles := getDirResult["childFileList"].([]interface{})

	for _, f := range childFiles {
		file := f.(map[string]interface{})
		if file["_id"].(string) == fileID {
			t.Error("Deleted file still exists")
		}
	}
}

func TestFileSetMetadata(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-setmeta-%d", time.Now().Unix())
	// Create bucket and file
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createFileReq := map[string]interface{}{
		"name":              "file-for-metadata.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{"version": 1},
		"encryptedMetaData": "initial-encrypted",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)
	fileID := createResult["fileId"].(string)

	// Update metadata
	setMetaReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
		"metaData": map[string]interface{}{"version": 2, "updated": true},
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/set-metadata", setMetaReq, adminAPIKey)

	// Update encrypted metadata separately
	setEncMetaReq := map[string]interface{}{
		"bucketId":          bucketID,
		"fileId":            fileID,
		"encryptedMetaData": "updated-encrypted",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/set-encrypted-metadata", setEncMetaReq, adminAPIKey)

	// Verify via get
	getReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
	}
	getResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/get", getReq, adminAPIKey)
	file := getResult["file"].(map[string]interface{})

	metaData := file["metaData"].(map[string]interface{})
	if version, ok := metaData["version"].(float64); !ok || version != 2 {
		t.Errorf("Expected version=2, got %v", version)
	}
	if updated, ok := metaData["updated"].(bool); !ok || !updated {
		t.Errorf("Expected updated=true, got %v", updated)
	}
	if encMeta, ok := file["encryptedMetaData"].(string); !ok || encMeta != "updated-encrypted" {
		t.Errorf("Expected encryptedMetaData='updated-encrypted', got %v", encMeta)
	}
}

func TestFileMove(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-move-%d", time.Now().Unix())
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
		"name":              "source-dir",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createAResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirAReq, adminAPIKey)
	sourceDirID := createAResult["directoryId"].(string)

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

	// Create file in source directory
	createFileReq := map[string]interface{}{
		"name":              "file-to-move.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": sourceDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createFileResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)
	fileID := createFileResult["fileId"].(string)

	// Move file
	moveReq := map[string]interface{}{
		"bucketId":             bucketID,
		"fileId":               fileID,
		"newParentDirectoryId": targetDirID,
		"newName":              "file-to-move.txt",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/move", moveReq, adminAPIKey)

	// Verify it's no longer in source directory
	getSourceReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": sourceDirID,
	}
	getSourceResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getSourceReq, adminAPIKey)
	sourceChildFiles := getSourceResult["childFileList"].([]interface{})

	for _, f := range sourceChildFiles {
		file := f.(map[string]interface{})
		if file["_id"].(string) == fileID {
			t.Error("Moved file still in source directory")
		}
	}

	// Verify it's now in target directory
	getTargetReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": targetDirID,
	}
	getTargetResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getTargetReq, adminAPIKey)
	targetChildFiles := getTargetResult["childFileList"].([]interface{})

	found := false
	for _, f := range targetChildFiles {
		file := f.(map[string]interface{})
		if file["_id"].(string) == fileID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Moved file not found in target directory")
	}
}

func TestFileGetAfterDelete(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-file-getdel-%d", time.Now().Unix())
	// Create bucket and file
	createBucketReq := map[string]interface{}{
		"name":      bucketName,
		"cryptSpec": "aes-256-gcm",
		"cryptData": "test-crypt-data",
		"metaData":  map[string]interface{}{},
	}
	bucketResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/bucket/create", createBucketReq, adminAPIKey)
	bucketID := bucketResult["bucketId"].(string)
	rootDirID := bucketResult["rootDirectoryId"].(string)

	createFileReq := map[string]interface{}{
		"name":              "file-to-delete-check.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted-data",
	}
	createResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createFileReq, adminAPIKey)
	fileID := createResult["fileId"].(string)

	// Delete file
	deleteReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/delete", deleteReq, adminAPIKey)

	// Try to get deleted file - should fail
	getReq := map[string]interface{}{
		"bucketId": bucketID,
		"fileId":   fileID,
	}
	_, result, _ := testutil.CallPostJSON(httpClient, baseURL+"/api/file/get", getReq, adminAPIKey)
	testutil.AssertErrorCode(t, result, "FILE_NOT_IN_BUCKET")
}
