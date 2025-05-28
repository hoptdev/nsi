package join_models

import models "nsi/internal/domain"

type DashboardWithRight struct {
	models.Dashboard
	AccessType models.GrantType
}
