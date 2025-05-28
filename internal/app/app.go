package app

import (
	"log/slog"
	grpc_client "nsi/internal/app/grpc"
	httpapp "nsi/internal/app/http"
	grpcHandler "nsi/internal/auth"
	"nsi/internal/config"
	"nsi/internal/services/dashboard"
	grpcService "nsi/internal/services/grpc"
	"nsi/internal/services/rights"
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
	grpcservice := grpcService.New(log, grpcClient)

	grpcHandler := grpcHandler.NewHandler(grpcservice)

	dashboardService := dashboard.New(log, storage, storage, storage, storage)
	widgetService := widget.New(log, storage, storage, storage, storage)
	rightsService := rights.New(log, storage, storage, storage)

	server := httpapp.New(log, cfg.Server.Port, cfg.Server.Timeout, rightsService, grpcHandler, grpcservice, dashboardService, widgetService)

	return &App{
		HttpServer: server,
	}
}
