package join_models

import models "nsi/internal/domain"

type WidgetWithRight struct {
	models.Widget
	AccessType models.GrantType
}
