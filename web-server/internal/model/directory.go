package model

import "time"

// Directory represents the directories table.
type Directory struct {
	ID                  string
	BucketID            string
	ParentDirectoryID   *string // NULL for root
	Name                string
	MetaData            []byte // JSONB
	EncryptedMetaData   string
	CreatedByUserID     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
