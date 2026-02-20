package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/handler"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/storage"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/repository"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/router"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/server"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Logging configuration
	zerolog.TimeFieldFormat = time.RFC3339Nano
	var logger zerolog.Logger
	switch cfg.Log.Format {
	case "console":
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	default:
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	log.Logger = logger.Level(parseLogLevel(cfg.Log.Level))

	// PostgreSQL connection pool
	dbPool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create PostgreSQL pool")
	}
	defer dbPool.Close()
	
	// Test database connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping PostgreSQL")
	}
	log.Info().Msg("PostgreSQL connection established")

	// Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()
	
	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping Redis")
	}
	log.Info().Msg("Redis connection established")

	// MinIO/Storage client
	minioClient, err := storage.NewMinIOClient(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.BucketName,
		cfg.MinIO.UseSSL,
		redisClient,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create MinIO client")
	}
	if err := minioClient.EnsureBucket(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure MinIO bucket exists")
	}
	log.Info().Msg("MinIO connection established")

	// Repositories
	userRepo := repository.NewUserRepository(dbPool)
	sessionRepo := repository.NewSessionRepository(dbPool)
	bucketRepo := repository.NewBucketRepository(dbPool)
	directoryRepo := repository.NewDirectoryRepository(dbPool)
	fileRepo := repository.NewFileRepository(dbPool)
	blobRepo := repository.NewBlobRepository(dbPool)

	// Services
	sessionSvc := service.NewSessionService(redisClient, sessionRepo, cfg)
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(sessionSvc, userSvc, cfg)
	adminSvc := service.NewAdminService(userRepo, sessionSvc, cfg)
	bucketSvc := service.NewBucketService(bucketRepo, directoryRepo)
	directorySvc := service.NewDirectoryService(directoryRepo, fileRepo)
	fileSvc := service.NewFileService(fileRepo)
	blobSvc := service.NewBlobService(blobRepo, minioClient)
	metricsSvc := service.NewMetricsService(minioClient)

	// Seed default admin
	if err := adminSvc.CreateDefaultAdminIfNotExists(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to create default admin user")
	}

	// Handlers
	userHandler := handler.NewUserHandler(userSvc, sessionSvc, authSvc, cfg)
	adminHandler := handler.NewAdminHandler(adminSvc, userSvc)
	bucketHandler := handler.NewBucketHandler(bucketSvc)
	directoryHandler := handler.NewDirectoryHandler(bucketSvc, directorySvc)
	fileHandler := handler.NewFileHandler(bucketSvc, directorySvc, fileSvc, blobSvc)
	blobHandler := handler.NewBlobHandler(bucketSvc, fileSvc, blobSvc)
	metricsHandler := handler.NewMetricsHandler(metricsSvc)

	// Router & server
	r := router.New(cfg, dbPool, redisClient, authSvc, userHandler, adminHandler, bucketHandler, directoryHandler, fileHandler, blobHandler, metricsHandler)
	srv := server.New(cfg, r)

	if err := srv.ListenAndServe(ctx); err != nil {
		log.Fatal().Err(err).Msg("server exited with error")
	}
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

