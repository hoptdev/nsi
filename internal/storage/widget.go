package psql

import (
	"context"
	models "nsi/internal/domain"
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
