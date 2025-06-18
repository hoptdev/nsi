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

func (s *Storage) GetDashboardsWithRights(ctx context.Context, userId int) ([]join_models.DashboardWithRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
        SELECT d.id, d.name, ar.type 
        FROM dashboards d
        JOIN dashboardOnAccessRights dar ON d.id = dar.dashboardId
        JOIN accessRights ar ON dar.accessRightId = ar.id
        WHERE ar.userId = $1;
    `
	rows, err := conn.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []join_models.DashboardWithRight
	for rows.Next() {
		var item join_models.DashboardWithRight
		if err := rows.Scan(&item.Id, &item.Name, &item.AccessType); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *Storage) GetDashboardRightByData(ctx context.Context, userId int, dashboardId int) (*models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	var result models.AccessRight

	query := "SELECT ar.* FROM accessRights ar JOIN dashboardOnAccessRights d ON ar.id=d.accessRightId WHERE d.dashboardId=$1 AND ar.userId=$2;"
	row := conn.QueryRow(ctx, query, dashboardId, userId)
	if err := row.Scan(&result.Id, &result.UserId, &result.UserGroupId, &result.AccessToken, &result.Type); err != nil {
		return nil, err
	}

	return &result, nil
}
