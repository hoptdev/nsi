package userController

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	grpcHandler "nsi/internal/auth"
	"time"
)

type userHelper struct {
	log      *slog.Logger
	timeout  time.Duration
	handlers UserHandlers
}

type UserHandlers interface {
	SignIn(ctx context.Context, login string, password string) (refresh string, access string, err error)
	SignUp(ctx context.Context, login string, password string) (bool, error)
	Refresh(ctx context.Context, token string) (string, error)
}

func Register(logger *slog.Logger, mux *http.ServeMux, t time.Duration, grpc *grpcHandler.Handler, handlers UserHandlers) {
	helper := &userHelper{logger, t, handlers}

	mux.HandleFunc("POST /user/signin", helper.SignIn())
	mux.HandleFunc("POST /user/signup", helper.SignUp())
	mux.HandleFunc("POST /token/refresh", helper.Refresh())
}

func (d *userHelper) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		params := struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		rtoken, atoken, err := d.handlers.SignIn(ctx, params.Login, params.Password)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}
		// json encode better
		res := fmt.Sprintf(`{ 
			"RefreshToken": "%v",
			"AccessToken": "%v"
		}`, rtoken, atoken)

		fmt.Fprint(w, res)
	}
}

// todo validator
func (d *userHelper) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		params := struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		id, err := d.handlers.SignUp(ctx, params.Login, params.Password)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}

func (d *userHelper) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.log.Info(fmt.Sprintf("[%v] [%v] request", r.Method, r.URL.Path))

		ctx, cancel := context.WithTimeout(r.Context(), d.timeout)
		defer cancel()

		params := struct {
			Token string `json:"refreshToken"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&params)

		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		id, err := d.handlers.Refresh(ctx, params.Token)
		if err != nil {
			d.log.Error(err.Error())

			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, id)
	}
}
