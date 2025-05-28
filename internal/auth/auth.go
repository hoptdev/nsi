package authHandler

import (
	"context"
	"net/http"
	grpcService "nsi/internal/services/grpc"
	"strconv"
	"time"
)

type Handler struct {
	service *grpcService.Service
}

func NewHandler(service *grpcService.Service) *Handler {
	return &Handler{service}
}

func (handler *Handler) ValidateHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		resp, err := handler.service.ValidateToken(ctx, r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if !resp.IsValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Add("UserId", strconv.Itoa(int(resp.UserId)))

		next(w, r)
	}
}
