package dashboard

import (
	"context"
	"errors"
	"log/slog"
	models "nsi/internal/domain"
	join_models "nsi/internal/domain/join"
	dashboardController "nsi/internal/http/dashboard"
)

var (
	ErrDashboardNotFound = errors.New("dashboard not found")
	ErrDashboardInvalid  = errors.New("token invalid")
	ErrorUpdateFailed    = errors.New("update token failed")
)

type Service struct {
	log               *slog.Logger
	dashboardUpdater  DashboardUpdater
	dashboardProvider DashboardProvider
	dashboardCreator  DashboardCreator
	dashboardRemover  DashboardRemover
}

type DashboardProvider interface {
	GetDashboard(ctx context.Context, model *models.Dashboard) error
	GetDashboardsWithRights(ctx context.Context, userId int) ([]join_models.DashboardWithRight, error)
}

type DashboardCreator interface {
	CreateDashboard(ctx context.Context, model *models.Dashboard) error
}

type DashboardUpdater interface {
}

type DashboardRemover interface {
	DeleteDashboard(ctx context.Context, id int) error
}

func New(log *slog.Logger, updater DashboardUpdater, provider DashboardProvider, creator DashboardCreator, remover DashboardRemover) *Service {
	return &Service{log, updater, provider, creator, remover}
}

func (service *Service) Create(ctx context.Context, name string, parentId *int, ownerId int, rightService dashboardController.RightHandler) (id int, err error) {
	model := &models.Dashboard{Id: 0, Name: name, ParentId: parentId}

	//todo: transaction
	err = service.dashboardCreator.CreateDashboard(ctx, model)
	if err != nil {
		return 0, err
	}

	_, err = rightService.Create(ctx, &id, nil, ownerId, models.Admin)
	if err != nil {
		return 0, err
	}

	return model.Id, nil
}

func (service *Service) GetDashboard(ctx context.Context, id int) (*models.Dashboard, error) {
	model := &models.Dashboard{Id: id}

	err := service.dashboardProvider.GetDashboard(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (service *Service) GetDashboardsWithAccess(ctx context.Context, userId int) ([]join_models.DashboardWithRight, error) {
	result, err := service.dashboardProvider.GetDashboardsWithRights(ctx, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *Service) Delete(ctx context.Context, id int) error {
	return service.dashboardRemover.DeleteDashboard(ctx, id)
}

func (service *Service) Update(ctx context.Context, id int, dashboard models.Dashboard) error {
	return nil
}
