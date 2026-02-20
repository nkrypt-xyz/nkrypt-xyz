package handler

import (
	"net/http"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type MetricsHandler struct {
	metricsSvc *service.MetricsService
}

func NewMetricsHandler(metricsSvc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsSvc: metricsSvc}
}

// GetSummary handles POST /api/metrics/get-summary
func (h *MetricsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	authData := middleware.GetAuthData(r.Context())
	if authData == nil {
		SendErrorResponse(w, apperror.NewUserError("ACCESS_DENIED", "Authentication required"))
		return
	}

	diskUsage, err := h.metricsSvc.GetDiskUsage(r.Context())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendSuccess(w, &model.MetricsGetSummaryResponse{
		HasError: false,
		Disk: model.MetricsDiskResponse{
			UsedBytes:  diskUsage.UsedBytes,
			TotalBytes: diskUsage.TotalBytes,
		},
	})
}
