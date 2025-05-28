package rights

import (
	"context"
	"errors"
	"log/slog"
	models "nsi/internal/domain"
)

var (
	ErrRightNotFound   = errors.New("Right not found")
	ErrNotEnoughRights = errors.New("Not enough rights")
)

type Service struct {
	log            *slog.Logger
	rightsUpdater  RightsUpdater
	rightsProvider RightsProvider
	rightsRemover  RightsRemover
}
type RightsRemover interface {
	//DeleteDashboardRight(ctx context.Context, id int) error
	//DeleteWidgetRight(ctx context.Context, id int) error
}

type RightsProvider interface {
	GetDashboardRightByData(ctx context.Context, userId int, dashboardId int) (*models.AccessRight, error) //названия конечно очень отражают суть)))
	GetWidgetRightByData(ctx context.Context, userId int, widgetIdId int) (*models.AccessRight, error)
}

type RightsUpdater interface {
}

func New(log *slog.Logger, updater RightsUpdater, provider RightsProvider, remover RightsRemover) *Service {
	return &Service{log, updater, provider, remover}
}

func (service *Service) CheckDashboardRight(ctx context.Context, userId int, dashboardId int, rightType models.GrantType) (err error) {
	return service.checkRight(ctx, userId, rightType, &dashboardId, nil)
}
func (service *Service) CheckWidgetRight(ctx context.Context, userId int, widgetId int, rightType models.GrantType) (err error) {
	return service.checkRight(ctx, userId, rightType, nil, &widgetId)
}

func (service *Service) checkRight(ctx context.Context, userId int, rightType models.GrantType, dashboardId *int, widgetId *int) error {
	var right *models.AccessRight
	var err error

	if dashboardId != nil {
		right, err = service.rightsProvider.GetDashboardRightByData(ctx, userId, *dashboardId)
	} else if widgetId != nil {
		right, err = service.rightsProvider.GetWidgetRightByData(ctx, userId, *widgetId)
	}

	if err != nil || right == nil {
		return ErrRightNotFound
	}

	if right.Type > rightType {
		return ErrNotEnoughRights
	}

	return nil
}
