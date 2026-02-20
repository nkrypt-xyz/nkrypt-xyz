package model

// UserResponse is the JSON representation of a user in API responses.
type UserResponse struct {
	ID                string          `json:"_id"`
	UserName          string          `json:"userName"`
	DisplayName       string          `json:"displayName"`
	IsBanned          bool            `json:"isBanned"`
	GlobalPermissions map[string]bool `json:"globalPermissions,omitempty"`
}

// SessionResponse is the minimal session representation used in login/assert.
type SessionResponse struct {
	ID string `json:"_id"`
}

// DirectoryResponse is the API representation of a directory.
type DirectoryResponse struct {
	ID                     string      `json:"_id"`
	BucketID               string      `json:"bucketId"`
	ParentDirectoryID      *string     `json:"parentDirectoryId"`
	Name                   string      `json:"name"`
	MetaData               interface{} `json:"metaData"`
	EncryptedMetaData      string      `json:"encryptedMetaData"`
	CreatedByUserIdentifier string     `json:"createdByUserIdentifier"`
	CreatedAt              int64       `json:"createdAt"`
	UpdatedAt              int64       `json:"updatedAt"`
}

// FileResponse is the API representation of a file.
type FileResponse struct {
	ID                       string      `json:"_id"`
	BucketID                 string      `json:"bucketId"`
	ParentDirectoryID        string      `json:"parentDirectoryId"`
	Name                    string      `json:"name"`
	MetaData                interface{} `json:"metaData"`
	EncryptedMetaData       string      `json:"encryptedMetaData"`
	SizeAfterEncryptionBytes int64      `json:"sizeAfterEncryptionBytes"`
	CreatedByUserIdentifier  string      `json:"createdByUserIdentifier"`
	CreatedAt                int64       `json:"createdAt"`
	UpdatedAt                int64       `json:"updatedAt"`
	ContentUpdatedAt         int64       `json:"contentUpdatedAt"`
}

// BucketAuthorizationResponse is one entry in bucketAuthorizations.
type BucketAuthorizationResponse struct {
	UserID      string         `json:"userId"`
	Notes       string         `json:"notes"`
	Permissions map[string]bool `json:"permissions"`
}

// BucketResponse is the API representation of a bucket in list responses.
type BucketResponse struct {
	ID                     string                        `json:"_id"`
	Name                   string                        `json:"name"`
	RootDirectoryID        string                        `json:"rootDirectoryId"`
	CryptSpec              string                        `json:"cryptSpec"`
	CryptData              string                        `json:"cryptData"`
	MetaData               interface{}                   `json:"metaData"`
	BucketAuthorizations   []BucketAuthorizationResponse `json:"bucketAuthorizations"`
	CreatedByUserIdentifier string                       `json:"createdByUserIdentifier"`
	CreatedAt              int64                         `json:"createdAt"`
	UpdatedAt              int64                         `json:"updatedAt"`
}

// UserListItemResponse is a minimal user in list responses.
type UserListItemResponse struct {
	ID          string `json:"_id"`
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
	IsBanned    bool   `json:"isBanned"`
}

// Per-endpoint success response structs (all include HasError for the envelope)

// CreateDirectoryResponse is the response for POST /api/directory/create
type CreateDirectoryResponse struct {
	HasError     bool   `json:"hasError"`
	DirectoryID  string `json:"directoryId"`
}

// GetDirectoryResponse is the response for POST /api/directory/get
type GetDirectoryResponse struct {
	HasError           bool               `json:"hasError"`
	Directory          DirectoryResponse  `json:"directory"`
	ChildDirectoryList []DirectoryResponse `json:"childDirectoryList"`
	ChildFileList      []FileResponse     `json:"childFileList"`
}

// CreateFileResponse is the response for POST /api/file/create
type CreateFileResponse struct {
	HasError bool   `json:"hasError"`
	FileID   string `json:"fileId"`
}

// GetFileResponse is the response for POST /api/file/get
type GetFileResponse struct {
	HasError bool         `json:"hasError"`
	File     FileResponse `json:"file"`
}

// CreateBucketResponse is the response for POST /api/bucket/create
type CreateBucketResponse struct {
	HasError         bool   `json:"hasError"`
	BucketID         string `json:"bucketId"`
	RootDirectoryID  string `json:"rootDirectoryId"`
}

// BucketListResponse is the response for POST /api/bucket/list
type BucketListResponse struct {
	HasError   bool             `json:"hasError"`
	BucketList []BucketResponse `json:"bucketList"`
}

// LoginResponse is the response for POST /api/user/login
type LoginResponse struct {
	HasError bool           `json:"hasError"`
	APIKey   string         `json:"apiKey"`
	User     UserResponse   `json:"user"`
	Session  SessionResponse `json:"session"`
}

// AssertResponse is the response for POST /api/user/assert
type AssertResponse struct {
	HasError bool           `json:"hasError"`
	APIKey   string         `json:"apiKey"`
	User     UserResponse   `json:"user"`
	Session  SessionResponse `json:"session"`
}

// UserListResponse is the response for POST /api/user/list
type UserListResponse struct {
	HasError bool                   `json:"hasError"`
	UserList []UserListItemResponse `json:"userList"`
}

// FindUserResponse is the response for POST /api/user/find
type FindUserResponse struct {
	HasError bool             `json:"hasError"`
	UserList []UserResponse   `json:"userList"`
}

// SessionListResponse is the response for POST /api/user/list-all-sessions
type SessionListResponse struct {
	HasError    bool               `json:"hasError"`
	SessionList []SessionListItem  `json:"sessionList"`
}

// CreateBlobResponse is the response for POST /api/blob/write
type CreateBlobResponse struct {
	HasError bool   `json:"hasError"`
	BlobID   string `json:"blobId"`
}

// WriteQuantizedResponse is the response for POST /api/blob/write-quantized
type WriteQuantizedResponse struct {
	HasError        bool   `json:"hasError"`
	BlobID          string `json:"blobId"`
	BytesTransferred int64 `json:"bytesTransfered"`
}

// AddUserResponse is the response for POST /api/admin/iam/add-user
type AddUserResponse struct {
	HasError bool   `json:"hasError"`
	UserID   string `json:"userId"`
}

// EmptySuccessResponse is used for endpoints that return only { hasError: false }
type EmptySuccessResponse struct {
	HasError bool `json:"hasError"`
}

// MetricsDiskResponse is the disk section of the metrics summary.
type MetricsDiskResponse struct {
	UsedBytes  int64 `json:"usedBytes"`
	TotalBytes int64 `json:"totalBytes"`
}

// MetricsGetSummaryResponse is the response for POST /api/metrics/get-summary
type MetricsGetSummaryResponse struct {
	HasError bool               `json:"hasError"`
	Disk     MetricsDiskResponse `json:"disk"`
}

