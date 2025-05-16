package grpcService

import (
	"context"
	"log/slog"
	grpc_client "nsi/internal/app/grpc"

	ssov1 "github.com/hoptdev/sso_protos/gen/go/sso"
)

type Service struct {
	log *slog.Logger
	app *grpc_client.App
}

func New(logger *slog.Logger, app *grpc_client.App) *Service {
	return &Service{logger, app}
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*ssov1.ValidateTokenResponse, error) {
	request := &ssov1.ValidateTokenRequest{
		RefreshToken: token,
	}
	resp, err := s.app.GRPCClient.Validate(ctx, request)

	return resp, err
}

func (s *Service) SignIn(ctx context.Context, login string, password string) (refresh string, access string, err error) {
	request := &ssov1.LoginRequest{
		Login:    login,
		Password: password,
	}
	resp, err := s.app.GRPCClient.Login(ctx, request)
	if err != nil {
		return "", "", err
	}

	return resp.RefreshToken, resp.AccessToken, err
}

func (s *Service) SignUp(ctx context.Context, login string, password string) (bool, error) {
	request := &ssov1.RegisterRequest{
		Login:    login,
		Password: password,
	}
	resp, err := s.app.GRPCClient.Register(ctx, request)
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (s *Service) Refresh(ctx context.Context, token string) (string, error) {
	request := &ssov1.RefreshRequest{
		RefreshToken: token,
	}
	resp, err := s.app.GRPCClient.Refresh(ctx, request)

	if err != nil {
		return "", err
	}

	return resp.RefreshToken, err
}
