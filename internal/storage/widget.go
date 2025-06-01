package psql

import (
	"context"
	models "nsi/internal/domain"
	join_models "nsi/internal/domain/join"
)

func (s *Storage) CreateWidget(ctx context.Context, model *models.Widget) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "INSERT INTO widgets (name, dashboardId, type, config) VALUES ($1, $2, $3, $4) RETURNING id;"
	row := conn.QueryRow(ctx, query, model.Name, model.DashboardId, model.WidgetType, model.Config)
	if err := row.Scan(&model.Id); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteWidget(ctx context.Context, id int) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "DELETE FROM dashboards WHERE id=$1;"
	_, err = conn.Exec(ctx, query, id)

	return err
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

func (s *Storage) GetWidgetRightByData(ctx context.Context, userId int, widgetId int) (*models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	var result models.AccessRight

	query := "SELECT ar.* FROM accessRights ar JOIN widgetOnAccessRights d ON ar.id=d.accessRightId WHERE d.widgetId=$1 AND ar.userId=$2;"
	row := conn.QueryRow(ctx, query, widgetId, userId)
	if err := row.Scan(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
