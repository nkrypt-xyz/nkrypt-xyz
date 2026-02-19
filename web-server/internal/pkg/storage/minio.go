package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
)

type MinIOClient struct {
	client      *minio.Client
	bucketName  string
	redisClient *redis.Client
}

type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

func NewMinIOClient(endpoint, accessKey, secretKey, bucketName string, useSSL bool, redisClient *redis.Client) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{
		client:      client,
		bucketName:  bucketName,
		redisClient: redisClient,
	}, nil
}

func (m *MinIOClient) EnsureBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{})
	}
	return nil
}

// UploadBlob streams a full blob upload
func (m *MinIOClient) UploadBlob(ctx context.Context, blobID string, reader io.Reader, size int64) (int64, error) {
	objectKey := "blobs/" + blobID

	info, err := m.client.PutObject(ctx, m.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return 0, err
	}

	return info.Size, nil
}

// DownloadBlob returns a reader for the blob and its size
func (m *MinIOClient) DownloadBlob(ctx context.Context, blobID string) (io.ReadCloser, int64, error) {
	objectKey := "blobs/" + blobID

	obj, err := m.client.GetObject(ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}

	stat, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, 0, err
	}

	return obj, stat.Size, nil
}

// GetBlobSize returns the size of a blob
func (m *MinIOClient) GetBlobSize(ctx context.Context, blobID string) (int64, error) {
	objectKey := "blobs/" + blobID
	stat, err := m.client.StatObject(ctx, m.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return 0, err
	}
	return stat.Size, nil
}

// DeleteBlob removes a blob from storage
func (m *MinIOClient) DeleteBlob(ctx context.Context, blobID string) error {
	objectKey := "blobs/" + blobID
	return m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
}

// AppendBlobChunk appends a chunk to a blob using temporary chunk storage.
// Chunks are stored as separate objects and tracked in Redis.
func (m *MinIOClient) AppendBlobChunk(ctx context.Context, blobID string, offset int64, reader io.Reader) (int64, error) {
	// Store chunk as temporary object with offset in the key
	chunkKey := m.getChunkKey(blobID, offset)
	
	info, err := m.client.PutObject(ctx, m.bucketName, chunkKey, reader, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return 0, err
	}
	
	// Track this chunk offset in Redis
	if m.redisClient != nil {
		redisKey := fmt.Sprintf("blob:chunks:%s", blobID)
		
		// Get existing offsets
		var offsets []int64
		data, err := m.redisClient.Get(ctx, redisKey).Result()
		if err == nil {
			_ = json.Unmarshal([]byte(data), &offsets)
		}
		
		// Add new offset
		offsets = append(offsets, offset)
		
		// Save back to Redis (24 hour TTL)
		offsetsJSON, _ := json.Marshal(offsets)
		_ = m.redisClient.Set(ctx, redisKey, offsetsJSON, 24*time.Hour).Err()
	}
	
	return info.Size, nil
}

// ComposeChunksToBlob composes all chunks into the final blob object.
// Retrieves chunk offsets from Redis and composes them in order.
func (m *MinIOClient) ComposeChunksToBlob(ctx context.Context, blobID string) error {
	finalKey := "blobs/" + blobID
	redisKey := fmt.Sprintf("blob:chunks:%s", blobID)
	
	// Get chunk offsets from Redis
	var offsets []int64
	if m.redisClient != nil {
		data, err := m.redisClient.Get(ctx, redisKey).Result()
		if err != nil {
			return fmt.Errorf("failed to get chunk offsets from Redis: %w", err)
		}
		if err := json.Unmarshal([]byte(data), &offsets); err != nil {
			return fmt.Errorf("failed to parse chunk offsets: %w", err)
		}
	}
	
	// Sort offsets to ensure correct order
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
	
	// Build list of source objects (chunks in order)
	sources := make([]minio.CopySrcOptions, len(offsets))
	for i, offset := range offsets {
		chunkKey := m.getChunkKey(blobID, offset)
		sources[i] = minio.CopySrcOptions{
			Bucket: m.bucketName,
			Object: chunkKey,
		}
	}
	
	// Compose chunks into final object
	_, err := m.client.ComposeObject(ctx, minio.CopyDestOptions{
		Bucket: m.bucketName,
		Object: finalKey,
	}, sources...)
	if err != nil {
		return err
	}
	
	// Clean up chunk files and Redis tracking
	for _, offset := range offsets {
		chunkKey := m.getChunkKey(blobID, offset)
		_ = m.client.RemoveObject(ctx, m.bucketName, chunkKey, minio.RemoveObjectOptions{})
	}
	if m.redisClient != nil {
		_ = m.redisClient.Del(ctx, redisKey).Err()
	}
	
	return nil
}

// getChunkKey returns the MinIO object key for a chunk
func (m *MinIOClient) getChunkKey(blobID string, offset int64) string {
	return fmt.Sprintf("blobs/%s.chunk.%d", blobID, offset)
}
