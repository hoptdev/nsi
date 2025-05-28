package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	grpcHandler "nsi/internal/auth"
	dashboardController "nsi/internal/http/dashboard"
	userController "nsi/internal/http/user"
	widgetController "nsi/internal/http/widget"
	"nsi/internal/services/dashboard"
	grpcService "nsi/internal/services/grpc"
	"nsi/internal/services/rights"
	"nsi/internal/services/widget"
	"time"

	"github.com/doganarif/govisual"
	"github.com/rs/cors"
)

type App struct {
	log    *slog.Logger
	mux    *http.ServeMux
	server *http.Server
	port   int
}

func New(log *slog.Logger, port int, timeout time.Duration, rights *rights.Service, grpc *grpcHandler.Handler, gservice *grpcService.Service, ds *dashboard.Service, ws *widget.Service) *App {
	mux := http.NewServeMux()
	dashboardController.Register(log, mux, timeout, grpc, ds, rights)
	widgetController.Register(log, mux, timeout, grpc, ws, rights)
	userController.Register(log, mux, timeout, grpc, gservice)

	return &App{log, mux, nil, port}
}

func (app *App) Run() {
	handler := govisual.Wrap(app.mux, govisual.WithRequestBodyLogging(true), govisual.WithResponseBodyLogging(true))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //not safe
		AllowedHeaders:   []string{"*"}, //not safe
		AllowCredentials: true,

		Debug: true,
	})
	handler = c.Handler(handler)

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
