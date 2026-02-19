package testutil

import (
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOHelper provides direct MinIO access for test verification
type MinIOHelper struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOHelper creates a new MinIO helper for tests
func NewMinIOHelper() (*MinIOHelper, error) {
	endpoint := os.Getenv("NK_MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("NK_MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("NK_MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	bucketName := os.Getenv("NK_MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "nkrypt-xyz-dev"
	}

	useSSL := os.Getenv("NK_MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOHelper{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// BlobExists checks if a blob exists in MinIO
func (m *MinIOHelper) BlobExists(ctx context.Context, blobID string) (bool, error) {
	objectName := "blobs/" + blobID
	_, err := m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" || errResp.Code == "NoSuchBucket" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetBlobSize returns the size of a blob in MinIO
func (m *MinIOHelper) GetBlobSize(ctx context.Context, blobID string) (int64, error) {
	objectName := "blobs/" + blobID
	info, err := m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchBucket" {
			return 0, nil
		}
		return 0, err
	}
	return info.Size, nil
}

// GetBlob retrieves blob data from MinIO
func (m *MinIOHelper) GetBlob(ctx context.Context, blobID string) ([]byte, error) {
	objectName := "blobs/" + blobID
	obj, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
}

// DeleteBlob removes a blob from MinIO (for cleanup)
func (m *MinIOHelper) DeleteBlob(ctx context.Context, blobID string) error {
	objectName := "blobs/" + blobID
	return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// ListBlobs lists all blobs with a given prefix
func (m *MinIOHelper) ListBlobs(ctx context.Context, prefix string) ([]string, error) {
	var blobs []string
	objectCh := m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:    "blobs/" + prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		// Remove "blobs/" prefix from the key
		blobID := object.Key[6:]
		blobs = append(blobs, blobID)
	}

	return blobs, nil
}
