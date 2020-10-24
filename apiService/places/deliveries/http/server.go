package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	api "github.com/contact-tracker/apiService/pkg/api/http"
	"github.com/contact-tracker/apiService/pkg/auth"
	"github.com/contact-tracker/apiService/places"
	t "github.com/contact-tracker/apiService/places/types"
)

const fiveSecondsTimeout = time.Second * 5

type handler struct {
	usecase places.PlaceService
	jwt     *auth.JWTService
}

func (d *handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		place, err := d.usecase.Get(ctx, id)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, place)
	}
}

func (d *handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		places, err := d.usecase.GetAll(ctx)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, places)
	}
}

func (d *handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.UpdatePlace{}
		api.ParseHTTPParams(r, req)

		req.ID = chi.URLParam(r, "id")
		resp, err := d.usecase.Update(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func (d *handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.CreatePlace{}
		api.ParseHTTPParams(r, req)
		place, err := d.usecase.Create(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		place.AuthToken, err = d.jwt.GenAccessToken(place)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, place)
	}
}

func (d *handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		err := d.usecase.Delete(ctx, id)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, nil)
	}
}

func (d *handler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.SignInReq{}
		api.ParseHTTPParams(r, req)
		place, err := d.usecase.SignIn(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		place.AuthToken, err = d.jwt.GenAccessToken(place)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, place)
	}
}

func (d *handler) Confirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		err := d.usecase.Confirm(ctx, id)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, nil)
	}
}

func NewServer(port, mongoDBName, mongoHost, mongoPlace, mongoPwd, placesHost, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd string) (server *api.Server, err error) {
	fmt.Println("Starting place http routes...")

	usecase, j, err := places.Init(mongoDBName, mongoHost, mongoPlace, mongoPwd, placesHost, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd)
	if err != nil {
		log.Panic(err)
	}

	h := &handler{
		usecase: usecase,
		jwt:     j,
	}

	server = api.NewServer(port)
	r := server.Router
	r.Post("/places", h.Create())
	r.Get("/places", j.AuthorizeHandler(h.GetAll()))
	r.Get("/places/{id}", j.AuthorizeHandler(h.Get()))
	r.Put("/places/{id}", j.AuthorizeHandler(h.Update()))
	r.Delete("/places/{id}", j.AuthorizeHandler(h.Delete()))
	r.Post("/places/login", h.SignIn())
	r.Get("/places/{id}/confirm", h.Confirm())

	return
}
