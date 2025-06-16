package widget

import (
	"context"
	"log/slog"
	models "nsi/internal/domain"
	join_models "nsi/internal/domain/join"
	widgetController "nsi/internal/http/widget"
)

type Service struct {
	log            *slog.Logger
	widgetUpdater  WidgetUpdater
	widgetProvider WidgetProvider
	widgetRemover  WidgetRemover
	widgetCreator  WidgetCreator
}

type WidgetProvider interface {
	GetWidgetsByDashboard(ctx context.Context, userId int, dashboardId int) (*[]join_models.WidgetWithRight, error)
	GetAllWidgetsByDashboard(ctx context.Context, dashboardId int) (*[]join_models.WidgetWithRight, error)
}

type WidgetRemover interface {
	DeleteWidget(ctx context.Context, id int) error
}

type WidgetCreator interface {
	CreateWidget(ctx context.Context, model *models.Widget) error
}

type WidgetUpdater interface {
}

func New(log *slog.Logger, updater WidgetUpdater, provider WidgetProvider, widgetRemover WidgetRemover, widgetCreator WidgetCreator) *Service {
	return &Service{log, updater, provider, widgetRemover, widgetCreator}
}

func (service *Service) Create(ctx context.Context, name string, dashboardId int, widgetType models.WidgetType, config string, ownerId int, rightService widgetController.RightHandler) (id int, err error) {
	model := &models.Widget{Id: 0, Name: name, DashboardId: dashboardId, WidgetType: widgetType, Config: config}

	err = service.widgetCreator.CreateWidget(ctx, model)
	if err != nil {
		return 0, err
	}

	_, err = rightService.Create(ctx, nil, &model.Id, ownerId, models.Admin)
	if err != nil {
		return 0, err
	}

	return model.Id, nil
}

func (service *Service) Delete(ctx context.Context, id int) error {
	return service.widgetRemover.DeleteWidget(ctx, id)
}

func (service *Service) Update(ctx context.Context, id int, dashboard models.Widget) error {
	return nil
}

func (service *Service) GetByDashboard(ctx context.Context, userId int, dashboardId int) (*[]join_models.WidgetWithRight, error) {
	result, err := service.widgetProvider.GetWidgetsByDashboard(ctx, userId, dashboardId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *Service) GetAllByDashboard(ctx context.Context, dashboardId int) (*[]join_models.WidgetWithRight, error) {
	result, err := service.widgetProvider.GetAllWidgetsByDashboard(ctx, dashboardId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
