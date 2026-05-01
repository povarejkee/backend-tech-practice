package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	core_logger "github.com/povarejkee/backend-tech-practice/internal/core/logger"
	core_http_server "github.com/povarejkee/backend-tech-practice/internal/core/transport/http/server"
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
	logger.Debug("Starting app!")

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
	)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
