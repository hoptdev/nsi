package psql

import (
	"context"
	models "nsi/internal/domain"
	join_models "nsi/internal/domain/join"
)

func (s *Storage) GetDashboardRightByData(ctx context.Context, userId int, dashboardId int) (*models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	var result models.AccessRight

	query := "SELECT ar.* FROM accessRights ar LEFT JOIN dashboardOnAccessRights d ON ar.id=d.accessRightId WHERE d.dashboardId=$1 AND ar.userId=$2;"
	row := conn.QueryRow(ctx, query, dashboardId, userId)
	if err := row.Scan(&result.Id, &result.UserId, &result.UserGroupId, &result.AccessToken, &result.Type); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Storage) GetWidgetRightByData(ctx context.Context, userId int, dashboardId int) (*models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	var result models.AccessRight

	query := "SELECT ar.* FROM accessRights ar JOIN widgetOnAccessRights d ON ar.id=d.accessRightId WHERE d.dashboardId=$1 AND ar.userId=$2;"
	row := conn.QueryRow(ctx, query, dashboardId, userId)
	if err := row.Scan(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Storage) GetWidgetsByDashboard(ctx context.Context, userId int, dashboardId int) (*[]join_models.WidgetWithRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT w.id, w.dashboardId, w.type, w.config, access.type FROM widgets w JOIN widgetOnAccessRights wr ON w.id=wr.widgetId JOIN accessRights access ON access.id=wr.accessRightId WHERE access.userId=$1;"

	count := 0
	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		count++
	}
	result := make([]join_models.WidgetWithRight, 0, count)

	rows, err = conn.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item join_models.WidgetWithRight
		if err := rows.Scan(&item.Id, &item.DashboardId, &item.WidgetType, &item.Config, &item.AccessType); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return &result, nil
}
