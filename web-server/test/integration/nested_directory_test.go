//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/test/testutil"
)

func TestNestedDirectoryCreation(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-nested-%d", time.Now().Unix())
	
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

	// Create nested structure: Root/DirA/DirAA/DirAAA
	
	// Create DirA
	createDirAReq := map[string]interface{}{
		"name":              "DirA",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{"level": 1},
		"encryptedMetaData": "level-1",
	}
	dirAResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirAReq, adminAPIKey)
	dirAID := dirAResult["directoryId"].(string)

	// Create DirAA inside DirA
	createDirAAReq := map[string]interface{}{
		"name":              "DirAA",
		"bucketId":          bucketID,
		"parentDirectoryId": dirAID,
		"metaData":          map[string]interface{}{"level": 2},
		"encryptedMetaData": "level-2",
	}
	dirAAResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirAAReq, adminAPIKey)
	dirAAID := dirAAResult["directoryId"].(string)

	// Create DirAAA inside DirAA
	createDirAAAReq := map[string]interface{}{
		"name":              "DirAAA",
		"bucketId":          bucketID,
		"parentDirectoryId": dirAAID,
		"metaData":          map[string]interface{}{"level": 3},
		"encryptedMetaData": "level-3",
	}
	dirAAAResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createDirAAAReq, adminAPIKey)
	dirAAAID := dirAAAResult["directoryId"].(string)

	// Verify root contains DirA
	getRootReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	rootResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getRootReq, adminAPIKey)
	rootChildDirs := rootResult["childDirectoryList"].([]interface{})
	
	if len(rootChildDirs) != 1 {
		t.Errorf("Expected 1 child in root, got %d", len(rootChildDirs))
	}
	if len(rootChildDirs) > 0 {
		dir := rootChildDirs[0].(map[string]interface{})
		if dir["_id"].(string) != dirAID {
			t.Error("Root child is not DirA")
		}
		if dir["name"].(string) != "DirA" {
			t.Errorf("Expected name 'DirA', got %s", dir["name"].(string))
		}
	}

	// Verify DirA contains DirAA
	getDirAReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": dirAID,
	}
	dirAResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getDirAReq, adminAPIKey)
	dirAChildDirs := dirAResult2["childDirectoryList"].([]interface{})
	
	if len(dirAChildDirs) != 1 {
		t.Errorf("Expected 1 child in DirA, got %d", len(dirAChildDirs))
	}
	if len(dirAChildDirs) > 0 {
		dir := dirAChildDirs[0].(map[string]interface{})
		if dir["_id"].(string) != dirAAID {
			t.Error("DirA child is not DirAA")
		}
	}

	// Verify DirAA contains DirAAA
	getDirAAReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": dirAAID,
	}
	dirAAResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getDirAAReq, adminAPIKey)
	dirAAChildDirs := dirAAResult2["childDirectoryList"].([]interface{})
	
	if len(dirAAChildDirs) != 1 {
		t.Errorf("Expected 1 child in DirAA, got %d", len(dirAAChildDirs))
	}
	if len(dirAAChildDirs) > 0 {
		dir := dirAAChildDirs[0].(map[string]interface{})
		if dir["_id"].(string) != dirAAAID {
			t.Error("DirAA child is not DirAAA")
		}
	}

	// Verify DirAAA is empty
	getDirAAAReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": dirAAAID,
	}
	dirAAAResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getDirAAAReq, adminAPIKey)
	dirAAAChildDirs := dirAAAResult2["childDirectoryList"].([]interface{})
	
	if len(dirAAAChildDirs) != 0 {
		t.Errorf("Expected 0 children in DirAAA, got %d", len(dirAAAChildDirs))
	}
}

func TestNestedDirectoryWithFiles(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-nested-files-%d", time.Now().Unix())
	
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

	// Create structure: Root/Projects/2024/Q1 with files at each level
	
	// Create Projects
	createProjectsReq := map[string]interface{}{
		"name":              "Projects",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	projectsResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createProjectsReq, adminAPIKey)
	projectsID := projectsResult["directoryId"].(string)

	// Create 2024 inside Projects
	create2024Req := map[string]interface{}{
		"name":              "2024",
		"bucketId":          bucketID,
		"parentDirectoryId": projectsID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	year2024Result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", create2024Req, adminAPIKey)
	year2024ID := year2024Result["directoryId"].(string)

	// Create Q1 inside 2024
	createQ1Req := map[string]interface{}{
		"name":              "Q1",
		"bucketId":          bucketID,
		"parentDirectoryId": year2024ID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	q1Result := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createQ1Req, adminAPIKey)
	q1ID := q1Result["directoryId"].(string)

	// Add files at each level
	
	// File in root
	createRootFileReq := map[string]interface{}{
		"name":              "README.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	rootFileResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createRootFileReq, adminAPIKey)
	rootFileID := rootFileResult["fileId"].(string)

	// File in Projects
	createProjectsFileReq := map[string]interface{}{
		"name":              "index.txt",
		"bucketId":          bucketID,
		"parentDirectoryId": projectsID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	projectsFileResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createProjectsFileReq, adminAPIKey)
	projectsFileID := projectsFileResult["fileId"].(string)

	// File in Q1
	createQ1FileReq := map[string]interface{}{
		"name":              "report.pdf",
		"bucketId":          bucketID,
		"parentDirectoryId": q1ID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	q1FileResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/file/create", createQ1FileReq, adminAPIKey)
	q1FileID := q1FileResult["fileId"].(string)

	// Verify root has 1 directory and 1 file
	getRootReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": rootDirID,
	}
	rootResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getRootReq, adminAPIKey)
	
	rootDirs := rootResult["childDirectoryList"].([]interface{})
	rootFiles := rootResult["childFileList"].([]interface{})
	
	if len(rootDirs) != 1 {
		t.Errorf("Expected 1 directory in root, got %d", len(rootDirs))
	}
	if len(rootFiles) != 1 {
		t.Errorf("Expected 1 file in root, got %d", len(rootFiles))
	}
	if len(rootFiles) > 0 {
		file := rootFiles[0].(map[string]interface{})
		if file["_id"].(string) != rootFileID {
			t.Error("Root file ID mismatch")
		}
	}

	// Verify Projects has 1 directory and 1 file
	getProjectsReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": projectsID,
	}
	projectsResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getProjectsReq, adminAPIKey)
	
	projectsDirs := projectsResult2["childDirectoryList"].([]interface{})
	projectsFiles := projectsResult2["childFileList"].([]interface{})
	
	if len(projectsDirs) != 1 {
		t.Errorf("Expected 1 directory in Projects, got %d", len(projectsDirs))
	}
	if len(projectsFiles) != 1 {
		t.Errorf("Expected 1 file in Projects, got %d", len(projectsFiles))
	}
	if len(projectsFiles) > 0 {
		file := projectsFiles[0].(map[string]interface{})
		if file["_id"].(string) != projectsFileID {
			t.Error("Projects file ID mismatch")
		}
	}

	// Verify Q1 has 0 directories and 1 file
	getQ1Req := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": q1ID,
	}
	q1Result2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getQ1Req, adminAPIKey)
	
	q1Dirs := q1Result2["childDirectoryList"].([]interface{})
	q1Files := q1Result2["childFileList"].([]interface{})
	
	if len(q1Dirs) != 0 {
		t.Errorf("Expected 0 directories in Q1, got %d", len(q1Dirs))
	}
	if len(q1Files) != 1 {
		t.Errorf("Expected 1 file in Q1, got %d", len(q1Files))
	}
	if len(q1Files) > 0 {
		file := q1Files[0].(map[string]interface{})
		if file["_id"].(string) != q1FileID {
			t.Error("Q1 file ID mismatch")
		}
	}
}

func TestNestedDirectoryOperations(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-nested-ops-%d", time.Now().Unix())
	
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

	// Create structure: Root/A/B and Root/C
	createAReq := map[string]interface{}{
		"name":              "A",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	aResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createAReq, adminAPIKey)
	aID := aResult["directoryId"].(string)

	createBReq := map[string]interface{}{
		"name":              "B",
		"bucketId":          bucketID,
		"parentDirectoryId": aID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	bResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createBReq, adminAPIKey)
	bID := bResult["directoryId"].(string)

	createCReq := map[string]interface{}{
		"name":              "C",
		"bucketId":          bucketID,
		"parentDirectoryId": rootDirID,
		"metaData":          map[string]interface{}{},
		"encryptedMetaData": "encrypted",
	}
	cResult := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/create", createCReq, adminAPIKey)
	cID := cResult["directoryId"].(string)

	// Move B from A to C (Root/A/B => Root/C/B)
	moveBReq := map[string]interface{}{
		"bucketId":             bucketID,
		"directoryId":          bID,
		"newParentDirectoryId": cID,
		"newName":              "B",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/move", moveBReq, adminAPIKey)

	// Verify A is now empty
	getAReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": aID,
	}
	aResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getAReq, adminAPIKey)
	aChildren := aResult2["childDirectoryList"].([]interface{})
	if len(aChildren) != 0 {
		t.Errorf("Expected A to be empty after move, got %d children", len(aChildren))
	}

	// Verify C now contains B
	getCReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": cID,
	}
	cResult2 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getCReq, adminAPIKey)
	cChildren := cResult2["childDirectoryList"].([]interface{})
	if len(cChildren) != 1 {
		t.Errorf("Expected C to have 1 child after move, got %d", len(cChildren))
	}
	if len(cChildren) > 0 {
		child := cChildren[0].(map[string]interface{})
		if child["_id"].(string) != bID {
			t.Error("C's child is not B")
		}
	}

	// Rename B to B-Renamed
	renameBReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": bID,
		"name":        "B-Renamed",
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/rename", renameBReq, adminAPIKey)

	// Verify rename worked
	cResult3 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getCReq, adminAPIKey)
	cChildren2 := cResult3["childDirectoryList"].([]interface{})
	if len(cChildren2) > 0 {
		child := cChildren2[0].(map[string]interface{})
		if child["name"].(string) != "B-Renamed" {
			t.Errorf("Expected name 'B-Renamed', got %s", child["name"].(string))
		}
	}

	// Delete B (now in C)
	deleteBReq := map[string]interface{}{
		"bucketId":    bucketID,
		"directoryId": bID,
	}
	testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/delete", deleteBReq, adminAPIKey)

	// Verify C is now empty
	cResult4 := testutil.CallPostJSONExpectSuccess(t, httpClient, baseURL+"/api/directory/get", getCReq, adminAPIKey)
	cChildren3 := cResult4["childDirectoryList"].([]interface{})
	if len(cChildren3) != 0 {
		t.Errorf("Expected C to be empty after delete, got %d children", len(cChildren3))
	}
}
