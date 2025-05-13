package grpc_client

import (
	"fmt"
	"log/slog"

	pb "github.com/hoptdev/sso_protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log        *slog.Logger
	gRPCconn   *grpc.ClientConn
	GRPCClient pb.AuthClient
	port       int
}

func New(log *slog.Logger, port int) *App {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", port), opts...)
	if err != nil {
		panic(err)
	}

	return &App{log, conn, nil, port}
}

func (app *App) Run() {
	const op = "grpc_client.Run"

	log := app.log.With(slog.String("op", op), slog.Int("port", app.port))

	log.Info("starting gRPC client")

	app.GRPCClient = pb.NewAuthClient(app.gRPCconn)
}
