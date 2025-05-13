package psql

import (
	"context"
	models "nsi/internal/domain"
)

func (s *Storage) CreateDashboard(ctx context.Context, model *models.Dashboard) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "INSERT INTO dashboards (name, parentId) VALUES ($1, $2) RETURNING id;"
	row := conn.QueryRow(ctx, query, model.Name, model.ParentId)
	if err := row.Scan(&model.Id); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteDashboard(ctx context.Context, id int) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "DELETE FROM dashboards WHERE id=$1;"
	_, err = conn.Exec(ctx, query, id)

	return err
}
