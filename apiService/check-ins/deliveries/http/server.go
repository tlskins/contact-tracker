package http

import (
	"fmt"
	"log"
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
	Usecase chk.CheckInService
	jwt     *auth.JWTService
}

func (d *handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		user, err := d.Usecase.Get(ctx, id)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, user)
	}
}

func (d *handler) GetHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		placeID := chi.URLParam(r, "placeId")
		now := time.Now()
		start := now.Add(time.Hour * -24)
		resp, err := d.Usecase.GetHistory(ctx, placeID, &start, &now)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func (d *handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.GetCheckIns{}
		api.ParseHTTPParams(r, req)

		resp, err := d.Usecase.GetAll(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func (d *handler) CheckIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.CreateCheckIn{}
		api.ParseHTTPParams(r, req)

		resp, err := d.Usecase.CheckIn(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func NewServer(port, mongoDBName, mongoHost, mongoCheckIn, mongoPwd, usersHost, placesHost, jwtKeyPath, jwtSecretPath, rpcPwd string) (server *api.Server, service *chk.CheckInService, err error) {
	fmt.Printf("Listening for check-ins on %s...\n", port)

	svc, j, err := chk.Init(mongoDBName, mongoHost, mongoCheckIn, mongoPwd, usersHost, placesHost, jwtKeyPath, jwtSecretPath, rpcPwd)
	if err != nil {
		log.Panic(err)
		return nil, nil, err
	}
	service = &svc

	h := &handler{
		Usecase: svc,
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
