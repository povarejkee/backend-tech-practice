package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	core_logger "github.com/povarejkee/backend-tech-practice/internal/core/logger"
	core_http_middleware "github.com/povarejkee/backend-tech-practice/internal/core/transport/http/middleware"
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

	// создание первого слоя для работы с users
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(nil)
	// получение ручек users
	usersRoutes := usersTransportHTTP.Routes()

	// создание мукса с версией апи
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	// регистрируем все ручки в муксе аля GET /users userMethod
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	// создаем инстанс сервака с еще одним муксом,
	// который будет смотреть на входящие запросы
	// например /api/v1/users
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(), core_http_middleware.Trace(),
	)
	// регаем версионный роутер в корневом муксе под префиксом /api/v1
	// потом обрезаем /api/v1 и перенаправимся в мукс,
	// который мы прокинули внутрь него вместе с apiVersionRouter
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
