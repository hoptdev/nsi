package dashboardController

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

type dashboardHelper struct {
	log      *slog.Logger
	timeout  time.Duration
	handlers DashboardHandlers
	rights   RightHandler
}

type DashboardHandlers interface {
	Create(ctx context.Context, name string, parentId *int, ownerId int, rightService RightHandler) (id int, err error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, dashboard models.Dashboard) error

	GetDashboard(ctx context.Context, id int) (*models.Dashboard, error)
	GetDashboardsWithAccess(ctx context.Context, userId int) ([]join_models.DashboardWithRight, error)
}

type RightHandler interface {
	Create(ctx context.Context, dashboardId *int, widgetdId *int, userId int, grantType models.GrantType) (id int, err error)
	CheckDashboardRight(ctx context.Context, userId int, dashboardId int, rightType models.GrantType) (err error)
}

func Register(logger *slog.Logger, mux *http.ServeMux, t time.Duration, grpc *grpcHandler.Handler, handlers DashboardHandlers, right RightHandler) {
	helper := &dashboardHelper{logger, t, handlers, right}

	mux.HandleFunc("POST /dashboard/create", grpc.ValidateHandler(helper.Create()))
	mux.HandleFunc("DELETE /dashboard/{id}", grpc.ValidateHandler(helper.Delete(models.Admin)))
	mux.HandleFunc("GET /dashboard/{id}", grpc.ValidateHandler(helper.GetDashboard(models.ReadOnly)))
	mux.HandleFunc("GET /dashboards", grpc.ValidateHandler(helper.GetDashboards()))
}

func (d *dashboardHelper) validateRole(ctx context.Context, w http.ResponseWriter, r *http.Request, role models.GrantType, dashboardId int) error {
	userId, _ := strconv.Atoi(r.Header.Get("UserId"))
	err := d.rights.CheckDashboardRight(ctx, userId, dashboardId, role)
	return err
}

func (d *dashboardHelper) GetDashboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		userId, _ := strconv.Atoi(r.Header.Get("UserId"))
		model, err := d.handlers.GetDashboardsWithAccess(ctx, userId)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		result, err := json.Marshal(model)

		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, string(result))
	}
}

func (d *dashboardHelper) GetDashboard(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		err = d.validateRole(ctx, w, r, role, int(id))
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		model, err := d.handlers.GetDashboard(ctx, int(id))
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		result, err := json.Marshal(model)

		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, string(result))
	}
}

func (d *dashboardHelper) Delete(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		err = d.validateRole(ctx, w, r, role, int(id))
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

func (d *dashboardHelper) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		params := struct {
			Name     string `json:"name"`
			ParentId *int   `json:"parentId"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		userId, _ := strconv.Atoi(r.Header.Get("UserId"))
		id, err := d.handlers.Create(ctx, params.Name, params.ParentId, userId, d.rights)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}
