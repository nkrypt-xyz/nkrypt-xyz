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
// For now, this is a simplified implementation that returns placeholder values.
// A full implementation would query MinIO bucket statistics or use admin APIs.
func (s *MetricsService) GetDiskUsage(ctx context.Context) (*model.DiskUsage, error) {
	// TODO: Implement actual MinIO bucket usage query
	// For MVP, return placeholder values
	// In production, use MinIO admin client or bucket info APIs
	return &model.DiskUsage{
		UsedBytes:  0,
		TotalBytes: 0,
	}, nil
}
