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
	ErrRightExists     = errors.New("Right exists")
)

type Service struct {
	log            *slog.Logger
	rightsUpdater  RightsUpdater
	rightsProvider RightsProvider
	rightsRemover  RightsRemover
	rightsCreator  RightsCreator
}

type RightsCreator interface {
	CreateAccessRight(ctx context.Context, right *models.AccessRight) error
	CreateDashboardAccessRight(ctx context.Context, dashboardId int, accessId int) (int, error)
	CreateWidgetAccessRight(ctx context.Context, widgetId int, accessId int) (int, error)
}

type RightsRemover interface {
	DeleteDashboardAccessRight(ctx context.Context, dashboardId int, rightId int) error
	DeleteWidgetAccessRight(ctx context.Context, widgetId int, rightId int) error
}

type RightsProvider interface {
	GetDashboardRightByData(ctx context.Context, userId int, dashboardId int) (*models.AccessRight, error) //названия конечно очень отражают суть)))
	GetWidgetRightByData(ctx context.Context, userId int, widgetIdId int) (*models.AccessRight, error)

	GetDashboardRights(ctx context.Context, dashboardId int) ([]models.AccessRight, error)
	GetWidgetRights(ctx context.Context, widgetdId int) ([]models.AccessRight, error)

	GetAccessRightByData(ctx context.Context, userId int, id int) (*models.AccessRight, error)
}

type RightsUpdater interface {
	UpdateAccessRight(ctx context.Context, id int, update models.AccessRight) error
	UpdateAccessRightType(ctx context.Context, id int, grant models.GrantType) error
}

func New(log *slog.Logger, updater RightsUpdater, provider RightsProvider, remover RightsRemover, creator RightsCreator) *Service {
	return &Service{log, updater, provider, remover, creator}
}

func (service *Service) CheckDashboardRight(ctx context.Context, userId int, dashboardId int, rightType models.GrantType) (right *models.AccessRight, err error) {
	return service.checkRight(ctx, userId, rightType, &dashboardId, nil, nil)
}
func (service *Service) CheckWidgetRight(ctx context.Context, userId int, widgetId int, rightType models.GrantType) (right *models.AccessRight, err error) {
	return service.checkRight(ctx, userId, rightType, nil, &widgetId, nil)
}
func (service *Service) CheckAccessRight(ctx context.Context, userId int, accessId int, rightType models.GrantType) (right *models.AccessRight, err error) {
	return service.checkRight(ctx, userId, rightType, nil, nil, &accessId)
}

func (service *Service) checkRight(ctx context.Context, userId int, rightType models.GrantType, dashboardId, widgetId, accessId *int) (*models.AccessRight, error) {
	var right *models.AccessRight
	var err error

	if dashboardId != nil {
		right, err = service.rightsProvider.GetDashboardRightByData(ctx, userId, *dashboardId)
	} else if widgetId != nil {
		right, err = service.rightsProvider.GetWidgetRightByData(ctx, userId, *widgetId)
	} else if accessId != nil {
		right, err = service.rightsProvider.GetAccessRightByData(ctx, userId, *accessId)
	}

	if err != nil || right == nil {
		return nil, ErrRightNotFound
	}

	if right.Type.ToInt() >= rightType.ToInt() {
		return right, nil
	}

	return nil, ErrNotEnoughRights
}

func (service *Service) Create(ctx context.Context, dashboardId *int, widgetdId *int, userId int, grantType models.GrantType) (id int, err error) {
	access := models.AccessRight{
		Id:     0,
		UserId: &userId,
		Type:   grantType,
	}

	//todo : transaction
	_, err = service.checkRight(ctx, userId, grantType, dashboardId, widgetdId, nil)

	err = service.rightsCreator.CreateAccessRight(ctx, &access)
	if err != nil {
		return 0, err
	}

	id = access.Id

	if dashboardId != nil {
		id, err = service.rightsCreator.CreateDashboardAccessRight(ctx, *dashboardId, access.Id)
	} else if widgetdId != nil {
		id, err = service.rightsCreator.CreateWidgetAccessRight(ctx, *widgetdId, access.Id)
	}

	return id, err
}

func (service *Service) Delete(ctx context.Context, dashboardId *int, widgetdId *int, rightId int) error {
	var err error
	if dashboardId != nil {
		err = service.rightsRemover.DeleteDashboardAccessRight(ctx, *dashboardId, rightId)
	} else if widgetdId != nil {
		err = service.rightsRemover.DeleteWidgetAccessRight(ctx, *widgetdId, rightId)
	}

	return err
}

func (service *Service) Update(ctx context.Context, userId int, id int, grant models.GrantType) (int, error) {
	return id, service.rightsUpdater.UpdateAccessRightType(ctx, id, grant)
}

func (service *Service) GetRights(ctx context.Context, id int, isDasboard bool) ([]models.AccessRight, error) {
	var result []models.AccessRight
	var err error

	if isDasboard {
		result, err = service.rightsProvider.GetDashboardRights(ctx, id)
	} else {
		result, err = service.rightsProvider.GetWidgetRights(ctx, id)
	}

	return result, err
}
