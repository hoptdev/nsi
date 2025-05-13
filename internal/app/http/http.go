package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	grpcHandler "nsi/internal/grpc"
	dashboardController "nsi/internal/http/dashboard"

	widgetController "nsi/internal/http/widget"
	"nsi/internal/services/dashboard"
	"nsi/internal/services/widget"
	"time"

	"github.com/doganarif/govisual"
)

type App struct {
	log    *slog.Logger
	mux    *http.ServeMux
	server *http.Server
	port   int
}

func New(log *slog.Logger, port int, timeout time.Duration, grpc *grpcHandler.Handler, ds *dashboard.Service, ws *widget.Service) *App {
	mux := http.NewServeMux()
	dashboardController.Register(log, mux, timeout, grpc, ds)
	widgetController.Register(log, mux, timeout, grpc, ws)

	return &App{log, mux, nil, port}
}

func (app *App) Run() {
	handler := govisual.Wrap(app.mux, govisual.WithRequestBodyLogging(true), govisual.WithResponseBodyLogging(true))

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", app.port),
		Handler: handler,
	}
	app.server = server

	app.log.Info(fmt.Sprintf("[server] start %v", app.server.Addr))

	if err := server.ListenAndServe(); err != nil {
		app.log.Error(err.Error())
	}
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
