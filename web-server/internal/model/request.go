package model

// LoginRequest matches the API spec for /api/user/login.
type LoginRequest struct {
	UserName string `json:"userName" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type AssertRequest struct{}

type LogoutRequest struct {
	Message string `json:"message" validate:"required,min=4,max=124"`
}

type LogoutAllSessionsRequest struct {
	Message string `json:"message" validate:"required,min=4,max=124"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName" validate:"required,min=4,max=128"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8,max=32"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=32"`
}

type FindUserFilter struct {
	By       string `json:"by" validate:"required,oneof=userName userId"`
	UserName string `json:"userName,omitempty"`
	UserID   string `json:"userId,omitempty"`
}

type FindUserRequest struct {
	Filters                  []FindUserFilter `json:"filters" validate:"required,dive"`
	IncludeGlobalPermissions bool             `json:"includeGlobalPermissions"`
}

// Bucket requests
type CreateBucketRequest struct {
	Name      string      `json:"name" validate:"required,min=1,max=64"`
	CryptSpec string      `json:"cryptSpec" validate:"required,min=1,max=64"`
	CryptData string      `json:"cryptData" validate:"required,min=1,max=2048"`
	MetaData  interface{} `json:"metaData" validate:"required"`
}

type RenameBucketRequest struct {
	BucketID string `json:"bucketId" validate:"required,len=16,alphanum"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
}

type SetBucketMetaDataRequest struct {
	BucketID string      `json:"bucketId" validate:"required,len=16,alphanum"`
	MetaData interface{} `json:"metaData" validate:"required"`
}

type SetBucketAuthorizationRequest struct {
	TargetUserID     string          `json:"targetUserId" validate:"required,len=16,alphanum"`
	BucketID         string          `json:"bucketId" validate:"required,len=16,alphanum"`
	PermissionsToSet map[string]bool `json:"permissionsToSet" validate:"required"`
}

type DestroyBucketRequest struct {
	BucketID string `json:"bucketId" validate:"required,len=16,alphanum"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
}

// Directory requests
type CreateDirectoryRequest struct {
	Name                string      `json:"name" validate:"required,min=1,max=256"`
	BucketID            string      `json:"bucketId" validate:"required,len=16,alphanum"`
	ParentDirectoryID   string      `json:"parentDirectoryId" validate:"required,len=16,alphanum"`
	MetaData            interface{} `json:"metaData" validate:"required"`
	EncryptedMetaData   string      `json:"encryptedMetaData" validate:"required,min=1,max=1048576"`
}

type GetDirectoryRequest struct {
	BucketID    string `json:"bucketId" validate:"required,len=16,alphanum"`
	DirectoryID string `json:"directoryId" validate:"required,len=16,alphanum"`
}

type RenameDirectoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=256"`
	BucketID    string `json:"bucketId" validate:"required,len=16,alphanum"`
	DirectoryID string `json:"directoryId" validate:"required,len=16,alphanum"`
}

type MoveDirectoryRequest struct {
	BucketID             string `json:"bucketId" validate:"required,len=16,alphanum"`
	DirectoryID          string `json:"directoryId" validate:"required,len=16,alphanum"`
	NewParentDirectoryID string `json:"newParentDirectoryId" validate:"required,len=16,alphanum"`
	NewName              string `json:"newName" validate:"required,min=1,max=256"`
}

type DeleteDirectoryRequest struct {
	BucketID    string `json:"bucketId" validate:"required,len=16,alphanum"`
	DirectoryID string `json:"directoryId" validate:"required,len=16,alphanum"`
}

type SetDirectoryMetaDataRequest struct {
	MetaData    interface{} `json:"metaData" validate:"required"`
	BucketID    string      `json:"bucketId" validate:"required,len=16,alphanum"`
	DirectoryID string      `json:"directoryId" validate:"required,len=16,alphanum"`
}

type SetDirectoryEncryptedMetaDataRequest struct {
	EncryptedMetaData string `json:"encryptedMetaData" validate:"required,min=1,max=1048576"`
	BucketID          string `json:"bucketId" validate:"required,len=16,alphanum"`
	DirectoryID       string `json:"directoryId" validate:"required,len=16,alphanum"`
}

// File requests
type CreateFileRequest struct {
	Name              string      `json:"name" validate:"required,min=1,max=256"`
	BucketID          string      `json:"bucketId" validate:"required,len=16,alphanum"`
	ParentDirectoryID string      `json:"parentDirectoryId" validate:"required,len=16,alphanum"`
	MetaData          interface{} `json:"metaData" validate:"required"`
	EncryptedMetaData string      `json:"encryptedMetaData" validate:"required,min=1,max=1048576"`
}

type GetFileRequest struct {
	BucketID string `json:"bucketId" validate:"required,len=16,alphanum"`
	FileID   string `json:"fileId" validate:"required,len=16,alphanum"`
}

type RenameFileRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=256"`
	BucketID string `json:"bucketId" validate:"required,len=16,alphanum"`
	FileID   string `json:"fileId" validate:"required,len=16,alphanum"`
}

type MoveFileRequest struct {
	BucketID             string `json:"bucketId" validate:"required,len=16,alphanum"`
	FileID               string `json:"fileId" validate:"required,len=16,alphanum"`
	NewParentDirectoryID string `json:"newParentDirectoryId" validate:"required,len=16,alphanum"`
	NewName              string `json:"newName" validate:"required,min=1,max=256"`
}

type DeleteFileRequest struct {
	BucketID string `json:"bucketId" validate:"required,len=16,alphanum"`
	FileID   string `json:"fileId" validate:"required,len=16,alphanum"`
}

type SetFileMetaDataRequest struct {
	MetaData interface{} `json:"metaData" validate:"required"`
	BucketID string      `json:"bucketId" validate:"required,len=16,alphanum"`
	FileID   string      `json:"fileId" validate:"required,len=16,alphanum"`
}

type SetFileEncryptedMetaDataRequest struct {
	EncryptedMetaData string `json:"encryptedMetaData" validate:"required,min=1,max=1048576"`
	BucketID          string `json:"bucketId" validate:"required,len=16,alphanum"`
	FileID            string `json:"fileId" validate:"required,len=16,alphanum"`
}

// Admin requests
type AddUserRequest struct {
	DisplayName string `json:"displayName" validate:"required,min=4,max=128"`
	UserName    string `json:"userName" validate:"required,min=4,max=32"`
	Password    string `json:"password" validate:"required,min=8,max=32"`
}

type SetGlobalPermissionsRequest struct {
	UserID            string          `json:"userId" validate:"required,len=16,alphanum"`
	GlobalPermissions map[string]bool `json:"globalPermissions" validate:"required"`
}

type SetBanningStatusRequest struct {
	UserID   string `json:"userId" validate:"required,len=16,alphanum"`
	IsBanned bool   `json:"isBanned"`
}

type OverwriteUserPasswordRequest struct {
	UserID      string `json:"userId" validate:"required,len=16,alphanum"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=32"`
}

