package models

type WidgetType string

const (
	Square WidgetType = "square"
)

type Widget struct {
	Id          int
	Name        string
	DashboardId int
	WidgetType  WidgetType
	Config      string
}
