package rightsController

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	grpcHandler "nsi/internal/auth"
	models "nsi/internal/domain"
	"strconv"
	"time"
)

type rightsHelper struct {
	log      *slog.Logger
	timeout  time.Duration
	handlers RightsHandlers
	rights   RightHandler
}

type RightsHandlers interface {
	Create(ctx context.Context, dashboardId *int, widgetdId *int, userId int, grantType models.GrantType) (id int, err error)
	Delete(ctx context.Context, dashboardId *int, widgetdId *int, rightId int) error
	Update(ctx context.Context, id int, dashboard models.Dashboard) error
	GetRights(ctx context.Context, id int, isDasboard bool) ([]models.AccessRight, error)
}

type RightHandler interface {
	CheckDashboardRight(ctx context.Context, userId int, dashboardId int, rightType models.GrantType) (err error)
	CheckWidgetRight(ctx context.Context, userId int, widgetId int, rightType models.GrantType) (err error)
}

func Register(logger *slog.Logger, mux *http.ServeMux, t time.Duration, grpc *grpcHandler.Handler, handlers RightsHandlers, right RightHandler) {
	helper := &rightsHelper{logger, t, handlers, right}

	mux.HandleFunc("POST /rights/create", grpc.ValidateHandler(helper.Create(models.Admin)))
	mux.HandleFunc("DELETE /rights/{rightId}", grpc.ValidateHandler(helper.Delete(models.Admin)))

	mux.HandleFunc("GET /rights/dashboard/{id}", grpc.ValidateHandler(helper.Get(models.Admin, true)))
	mux.HandleFunc("GET /rights/widget/{id}", grpc.ValidateHandler(helper.Get(models.Admin, false)))

	//mux.HandleFunc("PATCH /rights/{id}", grpc.ValidateHandler(helper.GetDashboard(models.ReadOnly)))
}

func (d *rightsHelper) validateRole(ctx context.Context, w http.ResponseWriter, r *http.Request, role models.GrantType, dashboardId *int, widgetId *int) (err error) {
	userId, _ := strconv.Atoi(r.Header.Get("UserId"))
	if dashboardId != nil {
		err = d.rights.CheckDashboardRight(ctx, userId, *dashboardId, role)
	} else if widgetId != nil {
		err = d.rights.CheckWidgetRight(ctx, userId, *widgetId, role)
	}
	return err
}

func (d *rightsHelper) Create(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		params := struct {
			UserId      int              `json:"userId"`
			DashboardId *int             `json:"dashboardId"`
			WidgetId    *int             `json:"widgetId"`
			Type        models.GrantType `json:"type"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		err = d.validateRole(ctx, w, r, role, params.DashboardId, params.WidgetId)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		id, err := d.handlers.Create(ctx, params.DashboardId, params.WidgetId, params.UserId, params.Type)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}

func (d *rightsHelper) Delete(role models.GrantType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		rightId, err1 := strconv.Atoi(r.PathValue("id"))

		params := struct {
			DashboardId *int `json:"dashboardId"`
			WidgetId    *int `json:"widgetId"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil || err1 != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		err = d.validateRole(ctx, w, r, role, params.DashboardId, nil)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		err = d.handlers.Delete(ctx, params.DashboardId, params.WidgetId, rightId)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, "Success")
	}
}

func (d *rightsHelper) Get(role models.GrantType, isDasboard bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		if isDasboard {
			err = d.validateRole(ctx, w, r, role, &id, nil)
		} else {
			err = d.validateRole(ctx, w, r, role, nil, &id)
		}

		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		data, err := d.handlers.GetRights(ctx, id, isDasboard)
		result, err1 := json.Marshal(data)

		if err != nil || err1 != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, result)
	}
}
