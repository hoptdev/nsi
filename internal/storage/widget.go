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

	query := "DELETE FROM widgets WHERE id=$1;"
	_, err = conn.Exec(ctx, query, id)

	return err
}

func (s *Storage) GetWidgetsByDashboard(ctx context.Context, userId int, dashboardId int) (*[]join_models.WidgetWithRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT w.id, w.dashboardId, w.type, w.config, access.type FROM widgets w JOIN widgetOnAccessRights wr ON w.id=wr.widgetId JOIN accessRights access ON access.id=wr.accessRightId WHERE w.dashboardId=$1 AND access.userId=$2;"

	count := 0
	rows, err := conn.Query(ctx, query, dashboardId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		count++
	}
	result := make([]join_models.WidgetWithRight, 0, count)

	rows, err = conn.Query(ctx, query, dashboardId, userId)
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

func (s *Storage) GetAllWidgetsByDashboard(ctx context.Context, dashboardId int) (*[]join_models.WidgetWithRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT w.id, w.dashboardId, w.type, w.config FROM widgets w WHERE w.dashboardId=$1;"

	rows, err := conn.Query(ctx, query, dashboardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []join_models.WidgetWithRight
	for rows.Next() {
		var item join_models.WidgetWithRight
		if err := rows.Scan(&item.Id, &item.DashboardId, &item.WidgetType, &item.Config); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return &result, nil
}

func (s *Storage) UpdatePosition(id int, x, y float64) error {
	query := `
        UPDATE widgets
		SET config = jsonb_set(
			jsonb_set(
				COALESCE(config, '{}')::jsonb, 
				'{position, x}',
				to_jsonb($1::float),
				true 
			),
			'{position, y}',
			to_jsonb($2::float),
			true
		)
		WHERE id = $3;
    `

	_, err := s.dbPool.Exec(context.Background(), query, x, y, id)
	return err
}

func (s *Storage) UpdateConfig(id int, config string) error {
	query := `
        UPDATE widgets 
        SET config = $1 
        WHERE id = $2;
    `

	_, err := s.dbPool.Exec(context.Background(), query, config, id)
	return err
}
