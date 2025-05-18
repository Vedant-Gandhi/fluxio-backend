package server

import (
	"fluxio-backend/pkg/common/schema"
	"fluxio-backend/pkg/config"
	"fluxio-backend/pkg/logger"
	"fluxio-backend/pkg/repository"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http"
	"fluxio-backend/pkg/transport/http/controller"
	"fluxio-backend/pkg/transport/http/middleware"
	"fluxio-backend/pkg/transport/http/routes"
	"fmt"
	"os"
)

func NewServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	db, err := pgsql.NewPgSQL(pgsql.PgSQLConfig{
		URL: cfg.Database.GetDatabaseURL(),
	})
	if err != nil {
		fmt.Println("Error loading database:", err)
		os.Exit(1)
	}

	logr := schema.Logger(logger.NewDefaultLogger())

	// Repositories
	userRepo := repository.NewUserRepository(db, &logr)

	videoRepo := repository.NewVideoRepository(db, repository.VideoRepositoryConfig{
		S3RawVideoBucketName:    cfg.VideoCfg.S3RawVideoBucketName,
		S3PublicVideoBucketName: cfg.VideoCfg.S3PublicVideoBucketName,
		S3ThumbnailBucketName:   cfg.VideoCfg.S3ThumbnailBucketName,
		S3Region:                cfg.VideoCfg.S3Region,
		S3AccessKey:             cfg.VideoCfg.S3AccessKey,
		S3SecretKey:             cfg.VideoCfg.S3SecretKey,
		S3Endpoint:              cfg.VideoCfg.S3Endpoint,
	},
		&logr)

	// Services
	if cfg.JWT.Secret == "" {
		fmt.Println("JWT secret is not set in the environment variables.")
		os.Exit(1)
	}

	jwtService := service.NewJWTService(cfg.JWT.Secret, &logr)
	userService := service.NewUserService(userRepo, jwtService, &logr)
	videoService := service.NewVideoService(videoRepo, &logr)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(userService, jwtService, &logr)
	middleware := middleware.NewMiddleware(authMiddleware)

	// Controllers
	authController := controller.NewAuthController(userService, &logr)
	videoController := controller.NewVideoController(videoService, &logr)

	// Pass the raw bucket name since that bucket's callback needs to be handled here
	s3Controller := controller.NewS3CallbackController(cfg.VideoCfg.S3RawVideoBucketName, videoService, &logr)

	// Route registrars
	authRouter := routes.NewAuthRouter(authController, middleware)
	videoRouter := routes.NewVideoRouter(videoController, middleware)
	s3Router := routes.NewAWSCallbackRouter(s3Controller, middleware)

	// Create and start HTTP router
	router := http.NewRouter(
		http.RouterConfig{
			Port:    cfg.Server.Port,
			Address: cfg.Server.Address,
		},
		authRouter,  // Pass the auth router as a route registrar
		videoRouter, // Pass the video router as a route registrar
		s3Router,
	)

	// Start the server
	if err := router.Start(); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
