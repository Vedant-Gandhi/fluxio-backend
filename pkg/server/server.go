package server

import (
	"fluxio-backend/pkg/config"
	"fluxio-backend/pkg/repository"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http"
	"fluxio-backend/pkg/transport/http/controller"
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

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo)

	// Controllers
	authController := controller.NewAuthController(userService)

	// Route registrars
	authRouter := routes.NewAuthRouter(authController)

	// Create and start HTTP router
	router := http.NewRouter(
		http.RouterConfig{
			Port:    cfg.Server.Port,
			Address: cfg.Server.Address,
		},
		authRouter, // Pass the auth router as a route registrar
	)

	// Start the server
	if err := router.Start(); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
