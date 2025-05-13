package dashboardController

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

type dashboardHelper struct {
	log      *slog.Logger
	timeout  time.Duration
	handlers DashboardHandlers
}

type DashboardHandlers interface {
	Create(ctx context.Context, name string, parentId *int) (id int, err error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, dashboard models.Dashboard) error
}

func Register(logger *slog.Logger, mux *http.ServeMux, t time.Duration, grpc *grpcHandler.Handler, handlers DashboardHandlers) {
	helper := &dashboardHelper{logger, t, handlers}

	mux.HandleFunc("POST /dashboard/create", grpc.ValidateHandler(helper.Create()))
	mux.HandleFunc("DELETE /dashboard/{id}", grpc.ValidateHandler(helper.Delete()))
}

func (d *dashboardHelper) Delete() http.HandlerFunc {
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

		id, err := d.handlers.Create(ctx, params.Name, params.ParentId)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}
