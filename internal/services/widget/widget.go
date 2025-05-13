package widget

import (
	"context"
	"log/slog"
	models "nsi/internal/domain"
)

type Service struct {
	log          *slog.Logger
	userUpdater  WidgetUpdater
	userProvider WidgetProvider
}

type WidgetProvider interface {
	CreateWidget(ctx context.Context, model *models.Widget) error
	DeleteWidget(ctx context.Context, id int) error
}

type WidgetUpdater interface {
}

func New(log *slog.Logger, updater WidgetUpdater, provider WidgetProvider) *Service {
	return &Service{log, updater, provider}
}

func (service *Service) Create(ctx context.Context, name string, dashboardId int, widgetType models.WidgetType, config string) (id int, err error) {
	model := &models.Widget{Id: 0, Name: name, DashboardId: dashboardId, WidgetType: widgetType, Config: config}

	err = service.userProvider.CreateWidget(ctx, model)
	if err != nil {
		return 0, err
	}

	return model.Id, nil
}

func (service *Service) Delete(ctx context.Context, id int) error {
	return service.userProvider.DeleteWidget(ctx, id)
}

func (service *Service) Update(ctx context.Context, id int, dashboard models.Widget) error {
	return nil
}
