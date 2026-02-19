package model

import "time"

// Bucket represents the buckets table.
type Bucket struct {
	ID               string
	Name             string
	CryptSpec        string
	CryptData        string
	MetaData         []byte // JSONB
	CreatedByUserID  string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// BucketPermission represents a row in bucket_user_permissions.
type BucketPermission struct {
	ID                     int64
	BucketID               string
	UserID                 string
	Notes                  string
	PermModify             bool
	PermManageAuthorization bool
	PermDestroy            bool
	PermViewContent        bool
	PermManageContent      bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// BucketListItem is used for /api/bucket/list response (with rootDirectoryId and bucketAuthorizations).
type BucketListItem struct {
	ID                     string
	Name                   string
	RootDirectoryID        string
	CryptSpec              string
	CryptData              string
	MetaData               []byte
	CreatedByUserID        string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	BucketAuthorizations   []BucketAuthorizationItem
}

// BucketAuthorizationItem is one entry in bucketAuthorizations array (API response format).
type BucketAuthorizationItem struct {
	UserID       string         `json:"userId"`
	Notes        string         `json:"notes"`
	Permissions  map[string]bool `json:"permissions"`
}
