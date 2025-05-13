package widgetController

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	models "nsi/internal/domain"
	grpcHandler "nsi/internal/grpc"
	"strconv"
	"time"
)

type widgetHelper struct {
	log      *slog.Logger
	timeout  time.Duration
	handlers WidgetHandlers
}

type WidgetHandlers interface {
	Create(ctx context.Context, name string, dashboardId int, widgetType models.WidgetType, config string) (id int, err error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, dashboard models.Widget) error
}

func Register(logger *slog.Logger, mux *http.ServeMux, t time.Duration, grpc *grpcHandler.Handler, handlers WidgetHandlers) {
	helper := &widgetHelper{logger, t, handlers}

	mux.HandleFunc("POST /widget/create", helper.Create())
	mux.HandleFunc("DELETE /widget/{id}", helper.Delete())
}

func (d *widgetHelper) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
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

func (d *widgetHelper) Create() http.HandlerFunc {
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

		id, err := d.handlers.Create(ctx, params.Name, params.DashboardId, models.WidgetType(params.WidgetType), params.Config)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}
