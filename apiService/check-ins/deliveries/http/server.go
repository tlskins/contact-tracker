package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	chk "github.com/contact-tracker/apiService/check-ins"
	t "github.com/contact-tracker/apiService/check-ins/types"
	api "github.com/contact-tracker/apiService/pkg/api/http"
	"github.com/contact-tracker/apiService/pkg/auth"
)

const fiveSecondsTimeout = time.Second * 5

type handler struct {
	usecase chk.CheckInService
	jwt     *auth.JWTService
}

func (d *handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		user, err := d.usecase.Get(ctx, id)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, user)
	}
}

func (d *handler) GetHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		placeID := chi.URLParam(r, "placeId")
		resp, err := d.usecase.GetHistory(ctx, placeID)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func (d *handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.GetCheckIns{}
		api.ParseHTTPParams(r, req)

		resp, err := d.usecase.GetAll(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func (d *handler) CheckIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.CreateCheckIn{}
		api.ParseHTTPParams(r, req)

		resp, err := d.usecase.CheckIn(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func NewServer(port, mongoDBName, mongoHost, mongoCheckIn, mongoPwd, usersHost, placesHost, jwtKeyPath, jwtSecretPath, rpcPwd string) (server *api.Server, err error) {
	fmt.Println("Starting check ins http routes...")

	usecase, j, err := chk.Init(mongoDBName, mongoHost, mongoCheckIn, mongoPwd, usersHost, placesHost, jwtKeyPath, jwtSecretPath, rpcPwd)
	if err != nil {
		return nil, err
	}

	h := &handler{
		usecase: usecase,
		jwt:     j,
	}

	server = api.NewServer(port)
	r := server.Router
	r.Get("/check-ins/{id}", j.AuthorizeHandler(h.Get()))
	r.Get("/check-ins/history/{placeId}", j.AuthorizeHandler(h.GetHistory()))
	r.Get("/check-ins", j.AuthorizeHandler(h.GetAll()))
	r.Post("/check-ins", j.AuthorizeHandler(h.CheckIn()))

	return
}
