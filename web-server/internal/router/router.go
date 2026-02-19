package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/handler"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/middleware"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

// New constructs the chi router with middleware and routes.
func New(cfg *config.Config, db *pgxpool.Pool, redisClient *redis.Client, authSvc *service.AuthService, userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, bucketHandler *handler.BucketHandler, directoryHandler *handler.DirectoryHandler, fileHandler *handler.FileHandler, blobHandler *handler.BlobHandler, metricsHandler *handler.MetricsHandler) http.Handler {
	r := chi.NewRouter()

	// Core middleware stack
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	healthHandler := &handler.HealthHandler{
		DB:    db,
		Redis: redisClient,
	}

	// Health probes
	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)

	// Prometheus metrics
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Post("/user/login", userHandler.Login)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			// User endpoints
			r.Post("/user/assert", userHandler.Assert)
			r.Post("/user/logout", userHandler.Logout)
			r.Post("/user/logout-all-sessions", userHandler.LogoutAllSessions)
			r.Post("/user/list-all-sessions", userHandler.ListAllSessions)
			r.Post("/user/list", userHandler.List)
			r.Post("/user/find", userHandler.Find)
			r.Post("/user/update-profile", userHandler.UpdateProfile)
			r.Post("/user/update-password", userHandler.UpdatePassword)

			// Admin endpoints
			r.Post("/admin/iam/add-user", adminHandler.AddUser)
			r.Post("/admin/iam/set-global-permissions", adminHandler.SetGlobalPermissions)
			r.Post("/admin/iam/set-banning-status", adminHandler.SetBanningStatus)
			r.Post("/admin/iam/overwrite-user-password", adminHandler.OverwriteUserPassword)

			// Bucket endpoints
			r.Post("/bucket/create", bucketHandler.Create)
			r.Post("/bucket/list", bucketHandler.List)
			r.Post("/bucket/rename", bucketHandler.Rename)
			r.Post("/bucket/set-metadata", bucketHandler.SetMetaData)
			r.Post("/bucket/set-authorization", bucketHandler.SetAuthorization)
			r.Post("/bucket/destroy", bucketHandler.Destroy)

			// Directory endpoints
			r.Post("/directory/create", directoryHandler.Create)
			r.Post("/directory/get", directoryHandler.Get)
			r.Post("/directory/rename", directoryHandler.Rename)
			r.Post("/directory/move", directoryHandler.Move)
			r.Post("/directory/delete", directoryHandler.Delete)
			r.Post("/directory/set-metadata", directoryHandler.SetMetaData)
			r.Post("/directory/set-encrypted-metadata", directoryHandler.SetEncryptedMetaData)

			// File endpoints
			r.Post("/file/create", fileHandler.Create)
			r.Post("/file/get", fileHandler.Get)
			r.Post("/file/rename", fileHandler.Rename)
			r.Post("/file/move", fileHandler.Move)
			r.Post("/file/delete", fileHandler.Delete)
			r.Post("/file/set-metadata", fileHandler.SetMetaData)
			r.Post("/file/set-encrypted-metadata", fileHandler.SetEncryptedMetaData)

			// Blob endpoints (streaming)
			r.Post("/blob/read/{bucketId}/{fileId}", blobHandler.Read)
			r.Post("/blob/write/{bucketId}/{fileId}", blobHandler.Write)
			r.Post("/blob/write-quantized/{bucketId}/{fileId}/{blobId}/{offset}/{shouldEnd}", blobHandler.WriteQuantized)

			// Metrics endpoints
			r.Post("/metrics/get-summary", metricsHandler.GetSummary)
		})
	})

	// Catch-all 404
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	})

	_ = cfg

	return r
}

