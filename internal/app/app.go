package app

import (
	"log/slog"
	grpc_client "nsi/internal/app/grpc"
	httpapp "nsi/internal/app/http"
	"nsi/internal/config"
	grpcHandler "nsi/internal/grpc"
	"nsi/internal/services/dashboard"
	"nsi/internal/services/widget"
	psql "nsi/internal/storage"
)

type App struct {
	HttpServer *httpapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	storage, err := psql.New(log, cfg.PSQL_Connect)
	if err != nil {
		panic(err)
	}

	grpcClient := grpc_client.New(log, cfg.Client.Port)
	grpcClient.Run()
	grpcHandler := grpcHandler.NewHandler(grpcClient)

	dashboardService := dashboard.New(log, storage, storage)
	widgetService := widget.New(log, storage, storage)

	server := httpapp.New(log, cfg.Server.Port, cfg.Server.Timeout, grpcHandler, dashboardService, widgetService)

	return &App{
		HttpServer: server,
	}
}
