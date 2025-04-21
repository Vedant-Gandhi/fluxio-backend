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

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	authController := controller.NewAuthController(userService)
	authRoute := routes.NewAuthRoute(authController)

	http.NewRouter(http.RouterConfig{
		Port:    cfg.Server.Port,
		Address: cfg.Server.Address,
	}, authRoute)

}
