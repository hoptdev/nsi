package widgetController

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	grpcHandler "nsi/internal/auth"
	models "nsi/internal/domain"
	join_models "nsi/internal/domain/join"
	"strconv"
	"time"
)

type widgetHelper struct {
	log      *slog.Logger
	timeout  time.Duration
	handlers WidgetHandlers
	rights   RightHandler
}
type RightHandler interface {
	Create(ctx context.Context, dashboardId *int, widgetdId *int, userId int, grantType models.GrantType) (id int, err error)

	CheckDashboardRight(ctx context.Context, userId int, dashboardId int, rightType models.GrantType) (right *models.AccessRight, err error)
	CheckWidgetRight(ctx context.Context, userId int, dashboardId int, rightType models.GrantType) (right *models.AccessRight, err error)
}

type WidgetHandlers interface {
	Create(ctx context.Context, name string, dashboardId int, widgetType models.WidgetType, config string, ownerId int, rightService RightHandler) (id int, err error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, dashboard models.Widget) error

	GetByDashboard(ctx context.Context, userId int, dashboardId int) (*[]join_models.WidgetWithRight, error)
	GetAllByDashboard(ctx context.Context, dashboardId int) (*[]join_models.WidgetWithRight, error)
}

func Register(logger *slog.Logger, mux *http.ServeMux, t time.Duration, grpc *grpcHandler.Handler, handlers WidgetHandlers, rights RightHandler) {
	helper := &widgetHelper{logger, t, handlers, rights}

	mux.HandleFunc("POST /widget/create", grpc.ValidateHandler(helper.Create(models.Update)))
	mux.HandleFunc("DELETE /widget/{id}", grpc.ValidateHandler(helper.Delete(models.Admin)))
	mux.HandleFunc("GET /widgets", grpc.ValidateHandler(helper.GetWidgets(models.ReadOnly)))
}

func (d *widgetHelper) validateRoleWidget(ctx context.Context, w http.ResponseWriter, r *http.Request, role models.GrantType, dashboardId int) error {
	userId, _ := strconv.Atoi(r.Header.Get("UserId"))
	_, err := d.rights.CheckWidgetRight(ctx, userId, dashboardId, role)
	return err
}
func (d *widgetHelper) validateRoleDashboard(ctx context.Context, w http.ResponseWriter, r *http.Request, role models.GrantType, dashboardId int) (models.GrantType, error) {
	userId, _ := strconv.Atoi(r.Header.Get("UserId"))
	result, err := d.rights.CheckDashboardRight(ctx, userId, dashboardId, role)
	return result.Type, err
}

func (d *widgetHelper) Delete(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		err = d.validateRoleWidget(ctx, w, r, role, int(id))
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		err = d.handlers.Delete(ctx, int(id))
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}
	}
}

func (d *widgetHelper) GetWidgets(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		value := r.URL.Query().Get("dashboardId")
		userId, _ := strconv.Atoi(r.Header.Get("UserId"))

		dashboardId, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		grant, err := d.validateRoleDashboard(ctx, w, r, role, dashboardId)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		var widgets *[]join_models.WidgetWithRight

		if grant == models.Admin {
			widgets, err = d.handlers.GetAllByDashboard(ctx, dashboardId)
		} else {
			widgets, err = d.handlers.GetByDashboard(ctx, userId, dashboardId)
		}

		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		result, err := json.Marshal(widgets)

		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, string(result))
	}
}

func (d *widgetHelper) Create(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		params := struct {
			Name        string `json:"name"`
			DashboardId int    `json:"dashboardId"`
			WidgetType  string `json:"type"`
			Config      string `json:"config"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		_, err = d.validateRoleDashboard(ctx, w, r, role, params.DashboardId)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		userId, _ := strconv.Atoi(r.Header.Get("UserId"))

		id, err := d.handlers.Create(ctx, params.Name, params.DashboardId, models.WidgetType(params.WidgetType), params.Config, userId, d.rights)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}
