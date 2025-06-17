package psql

import (
	"context"
	models "nsi/internal/domain"
)

func (s *Storage) UpdateAccessRight(ctx context.Context, id int, update models.AccessRight) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `
        UPDATE accessRights 
        SET userId = $1, userGroupId = $2, accessToken = $3, type = $4 
        WHERE id = $5;
    `
	_, err = conn.Exec(
		ctx,
		query,
		update.UserId,
		update.UserGroupId,
		update.AccessToken,
		update.Type,
		id,
	)
	return err
}

func (s *Storage) DeleteDashboardAccessRight(ctx context.Context, dashboardId int, rightId int) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `DELETE FROM accessRights
	USING widgetOnAccessRights a
	WHERE a.widgetId=$1 AND id=$2;`
	_, err = conn.Exec(ctx, query, dashboardId, rightId)
	return err
}

func (s *Storage) DeleteWidgetAccessRight(ctx context.Context, widgetId int, rightId int) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `DELETE FROM accessRights
	USING dashboardOnAccessRights a
	WHERE a.dashboardId=$1 AND id=$2;`
	_, err = conn.Exec(ctx, query, widgetId, rightId)
	return err
}

func (s *Storage) CreateAccessRight(ctx context.Context, right *models.AccessRight) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `
        INSERT INTO accessRights (userId, userGroupId, accessToken, type) 
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `
	return conn.QueryRow(
		ctx,
		query,
		right.UserId,
		right.UserGroupId,
		right.AccessToken,
		right.Type,
	).Scan(&right.Id)
}

func (s *Storage) CreateDashboardAccessRight(ctx context.Context, dashboardId int, accessId int) (int, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return accessId, err
	}
	defer conn.Release()

	query := `
        INSERT INTO dashboardOnAccessRights (accessRightId, dashboardId) 
        VALUES ($1, $2)
    `
	_, err = conn.Query(
		ctx,
		query,
		accessId,
		dashboardId,
	)

	return accessId, err
}

func (s *Storage) CreateWidgetAccessRight(ctx context.Context, widgetId int, accessId int) (int, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return accessId, err
	}
	defer conn.Release()

	query := `
        INSERT INTO widgetOnAccessRights (accessRightId, widgetId) 
        VALUES ($1, $2)
    `
	_, err = conn.Query(
		ctx,
		query,
		accessId,
		widgetId,
	)

	return accessId, err
}

func (s *Storage) GetDashboardRights(ctx context.Context, dashboardId int) ([]models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
        SELECT a.id, a.userId, a.userGroupId, a.accessToken, a.type
        FROM accessRights a
        JOIN dashboardOnAccessRights d ON d.accessRightId = a.id

        WHERE d.id = $1;
    `
	rows, err := conn.Query(ctx, query, dashboardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.AccessRight
	for rows.Next() {
		var item models.AccessRight
		if err := rows.Scan(&item.Id, &item.UserId, &item.UserGroupId, &item.AccessToken, &item.Type); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *Storage) GetWidgetRights(ctx context.Context, widgetdId int) ([]models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
        SELECT a.id, a.userId, a.userGroupId, a.accessToken, a.type
        FROM accessRights a
        JOIN widgetOnAccessRights d ON d.accessRightId = a.id

        WHERE d.widgetId = $1;
    `
	rows, err := conn.Query(ctx, query, widgetdId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.AccessRight
	for rows.Next() {
		var item models.AccessRight
		if err := rows.Scan(&item.Id, &item.UserId, &item.UserGroupId, &item.AccessToken, &item.Type); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (s *Storage) GetWidgetRightByData(ctx context.Context, userId int, widgetId int) (*models.AccessRight, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	var result models.AccessRight

	query := `WITH widget_dash AS (
		SELECT dashboardId
		FROM widgets 
		WHERE id = $1
		LIMIT 1
	)
	SELECT 
		ar.id, 
		ar.userId, 
		ar.usergroupId, 
		ar.accesstoken, 
		ar.type
	FROM accessRights ar
	LEFT JOIN widgetOnAccessRights wor 
		ON wor.accessRightId = ar.id 
		AND wor.widgetId = $1
	LEFT JOIN dashboardOnAccessRights dor 
		ON dor.accessRightId = ar.id
		AND dor.dashboardId = (SELECT dashboardId FROM widget_dash)
		AND ar.type = 'admin'
	WHERE ar.userId = $2
	AND (wor.widgetId IS NOT NULL OR dor.dashboardId IS NOT NULL)
	ORDER BY 
		CASE WHEN wor.widgetId IS NOT NULL THEN 0 ELSE 1 END
	LIMIT 1;`
	row := conn.QueryRow(ctx, query, widgetId, userId)
	if err := row.Scan(&result.Id, &result.UserId, &result.UserGroupId, &result.AccessToken, &result.Type); err != nil {
		return nil, err
	}

	return &result, nil
}
