package service

import (
	"context"
	"io"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/randstr"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/storage"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
)

type BlobService struct {
	blobRepo      *repository.BlobRepository
	storageClient *storage.MinIOClient
}

func NewBlobService(blobRepo *repository.BlobRepository, storageClient *storage.MinIOClient) *BlobService {
	return &BlobService{blobRepo: blobRepo, storageClient: storageClient}
}

func (s *BlobService) CreateInProgressBlob(ctx context.Context, bucketID, fileID, cryptoMetaHeaderContent, createdByUserID string) (*model.Blob, error) {
	blobID, err := randstr.GenerateID(16)
	if err != nil {
		return nil, apperror.NewDeveloperError("ID_GENERATION_FAILED", "Failed to generate blob ID.")
	}
	blob := &model.Blob{
		ID:                      blobID,
		BucketID:                bucketID,
		FileID:                  fileID,
		CryptoMetaHeaderContent: cryptoMetaHeaderContent,
		Status:                  "started",
		CreatedByUserID:         createdByUserID,
	}
	if err := s.blobRepo.Create(ctx, blob); err != nil {
		return nil, err
	}
	return blob, nil
}

func (s *BlobService) GetInProgressBlob(ctx context.Context, bucketID, fileID, blobID string) (*model.Blob, error) {
	blob, err := s.blobRepo.FindInProgressBlob(ctx, bucketID, fileID, blobID)
	if err != nil || blob == nil {
		return nil, apperror.NewUserError("BLOB_INVALID", "No in-progress blob found with the given ID")
	}
	return blob, nil
}

func (s *BlobService) MarkBlobFinished(ctx context.Context, blobID string) error {
	return s.blobRepo.MarkFinished(ctx, blobID)
}

func (s *BlobService) MarkBlobErroneous(ctx context.Context, blobID string) error {
	return s.blobRepo.MarkErroneous(ctx, blobID)
}

func (s *BlobService) FindLatestFinishedBlob(ctx context.Context, bucketID, fileID string) (*model.Blob, error) {
	return s.blobRepo.FindLatestFinishedBlob(ctx, bucketID, fileID)
}

func (s *BlobService) StreamBlobToWriter(ctx context.Context, blobID string, w io.Writer) (int64, error) {
	reader, _, err := s.storageClient.DownloadBlob(ctx, blobID)
	if err != nil {
		return 0, err
	}
	defer reader.Close()
	return io.Copy(w, reader)
}

func (s *BlobService) GetBlobSize(ctx context.Context, blobID string) (int64, error) {
	return s.storageClient.GetBlobSize(ctx, blobID)
}

func (s *BlobService) UploadBlobFromReader(ctx context.Context, blobID string, reader io.Reader, size int64) error {
	_, err := s.storageClient.UploadBlob(ctx, blobID, reader, size)
	return err
}

func (s *BlobService) RemoveAllOtherBlobs(ctx context.Context, bucketID, fileID, keepBlobID string) error {
	blobs, err := s.blobRepo.ListBlobsForFileExcluding(ctx, bucketID, fileID, keepBlobID)
	if err != nil {
		return err
	}
	for _, blob := range blobs {
		_ = s.storageClient.DeleteBlob(ctx, blob.ID) // Ignore errors (ENOENT is ok)
	}
	return s.blobRepo.DeleteBlobsForFileExcluding(ctx, bucketID, fileID, keepBlobID)
}

func (s *BlobService) RemoveAllBlobsOfFile(ctx context.Context, bucketID, fileID string) error {
	blobs, err := s.blobRepo.ListBlobsForFile(ctx, bucketID, fileID)
	if err != nil {
		return err
	}
	for _, blob := range blobs {
		_ = s.storageClient.DeleteBlob(ctx, blob.ID)
	}
	return s.blobRepo.DeleteAllBlobsForFile(ctx, bucketID, fileID)
}

// AppendChunkToBlob appends a chunk to an in-progress blob at the specified offset.
// Returns the number of bytes written.
func (s *BlobService) AppendChunkToBlob(ctx context.Context, blobID string, offset int64, reader io.Reader) (int64, error) {
	return s.storageClient.AppendBlobChunk(ctx, blobID, offset, reader)
}

// FinalizeChunkedBlob composes all uploaded chunks into the final blob.
// This should be called after the last chunk is uploaded.
func (s *BlobService) FinalizeChunkedBlob(ctx context.Context, blobID string) error {
	return s.storageClient.ComposeChunksToBlob(ctx, blobID)
}
