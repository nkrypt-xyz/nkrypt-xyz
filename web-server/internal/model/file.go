package model

import "time"

// File represents the files table.
type File struct {
	ID                        string
	BucketID                  string
	ParentDirectoryID         string
	Name                      string
	MetaData                  []byte // JSONB
	EncryptedMetaData         string
	SizeAfterEncryptionBytes  int64
	CreatedByUserID           string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	ContentUpdatedAt          time.Time
}
