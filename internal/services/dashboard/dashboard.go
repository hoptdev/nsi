package dashboard

import (
	"context"
	"errors"
	"log/slog"
	models "nsi/internal/domain"
)

var (
	ErrDashboardNotFound = errors.New("dashboard not found")
	ErrDashboardInvalid  = errors.New("token invalid")
	ErrorUpdateFailed    = errors.New("update token failed")
)

type Service struct {
	log          *slog.Logger
	userUpdater  DashboardUpdater
	userProvider DashboardProvider
}

type DashboardProvider interface {
	CreateDashboard(ctx context.Context, model *models.Dashboard) error
	DeleteDashboard(ctx context.Context, id int) error
}

type DashboardUpdater interface {
}

func New(log *slog.Logger, updater DashboardUpdater, provider DashboardProvider) *Service {
	return &Service{log, updater, provider}
}

func (service *Service) Create(ctx context.Context, name string, parentId *int) (id int, err error) {
	model := &models.Dashboard{Id: 0, Name: name, ParentId: parentId}

	err = service.userProvider.CreateDashboard(ctx, model)
	if err != nil {
		return 0, err
	}

	return model.Id, nil
}

func (service *Service) Delete(ctx context.Context, id int) error {
	return service.userProvider.DeleteDashboard(ctx, id)
}

func (service *Service) Update(ctx context.Context, id int, dashboard models.Dashboard) error {
	return nil
}
