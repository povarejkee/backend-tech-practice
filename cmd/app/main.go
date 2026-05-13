package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	core_logger "github.com/povarejkee/backend-tech-practice/internal/core/logger"
	core_postgres_pool "github.com/povarejkee/backend-tech-practice/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/povarejkee/backend-tech-practice/internal/core/transport/http/middleware"
	core_http_server "github.com/povarejkee/backend-tech-practice/internal/core/transport/http/server"
	users_postgres_repository "github.com/povarejkee/backend-tech-practice/internal/features/users/repository/postgres"
	users_service "github.com/povarejkee/backend-tech-practice/internal/features/users/service"
	users_transport_http "github.com/povarejkee/backend-tech-practice/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init app logger:", err)
	}
	defer logger.Close()

	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewCfgMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature..", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing http server..")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(), core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
