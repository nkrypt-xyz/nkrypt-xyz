package service

import (
	"context"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/storage"
)

type MetricsService struct {
	storageClient *storage.MinIOClient
}

func NewMetricsService(storageClient *storage.MinIOClient) *MetricsService {
	return &MetricsService{storageClient: storageClient}
}

// GetDiskUsage returns disk usage statistics from MinIO.
func (s *MetricsService) GetDiskUsage(ctx context.Context) (*model.DiskUsage, error) {
	usedBytes, totalBytes, err := s.storageClient.GetBucketUsage(ctx)
	if err != nil {
		return nil, err
	}

	return &model.DiskUsage{
		UsedBytes:  usedBytes,
		TotalBytes: totalBytes,
	}, nil
}
