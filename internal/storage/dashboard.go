package psql

import (
	"context"
	models "nsi/internal/domain"
	join_models "nsi/internal/domain/join"
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

func (s *Storage) GetDashboard(ctx context.Context, model *models.Dashboard) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	query := "SELECT d.* FROM dashboards d WHERE d.id=$1;"

	row := conn.QueryRow(ctx, query, model.Id)
	if err := row.Scan(&model.Id, &model.Name, &model.ParentId); err != nil {
		return err
	}

	return err
}

func (s *Storage) GetDashboardsWithAccess(ctx context.Context, userId int) ([]join_models.DashboardWithRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	query := "SELECT d.id, d.name, d.parentId, ar.type FROM dashboards d JOIN dashboardOnAccessRights access ON d.id=access.dashboardId JOIN accessRights ar ON ar.userId=$1;"

	count := 0
	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		count++
	}
	result := make([]join_models.DashboardWithRight, 0, count)

	rows, err = conn.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item join_models.DashboardWithRight
		if err := rows.Scan(&item.Id, &item.Name, &item.ParentId, &item.AccessType); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
