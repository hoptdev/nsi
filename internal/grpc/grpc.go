package grpcHandler

import (
	"context"
	"net/http"
	grpc_client "nsi/internal/app/grpc"
	"time"

	ssov1 "github.com/hoptdev/sso_protos/gen/go/sso"
)

type Handler struct {
	app *grpc_client.App
}

func NewHandler(app *grpc_client.App) *Handler {
	return &Handler{app}
}

func (handler *Handler) ValidateHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		request := &ssov1.ValidateTokenRequest{
			RefreshToken: r.Header.Get("Authorization"),
		}

		resp, err := handler.app.GRPCClient.Validate(ctx, request)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if !resp.IsValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
