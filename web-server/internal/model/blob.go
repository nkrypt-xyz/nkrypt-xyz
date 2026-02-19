package model

import "time"

// Blob represents the blobs table.
type Blob struct {
	ID                       string
	BucketID                 string
	FileID                   string
	CryptoMetaHeaderContent  string
	StartedAt                time.Time
	FinishedAt               *time.Time
	Status                   string // 'started', 'finished', 'error'
	CreatedByUserID          string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
