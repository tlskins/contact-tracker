package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	api "github.com/contact-tracker/apiService/pkg/api/http"
	"github.com/contact-tracker/apiService/pkg/auth"
	"github.com/contact-tracker/apiService/users"
	t "github.com/contact-tracker/apiService/users/types"
)

const fiveSecondsTimeout = time.Second * 5

type handler struct {
	usecase users.UserService
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

func (d *handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users, err := d.usecase.GetAll(ctx)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, users)
	}
}

func (d *handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.UpdateUser{}
		api.ParseHTTPParams(r, req)

		req.ID = chi.URLParam(r, "id")
		resp, err := d.usecase.Update(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, resp)
	}
}

func (d *handler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.SignInReq{}
		api.ParseHTTPParams(r, req)
		user, err := d.usecase.SignIn(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		user.AuthToken, err = d.jwt.GenAccessToken(user)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, user)
	}
}

func (d *handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := &t.CreateUser{}
		api.ParseHTTPParams(r, req)
		user, err := d.usecase.Create(ctx, req)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		user.AuthToken, err = d.jwt.GenAccessToken(user)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, user)
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

func (d *handler) Confirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		err := d.usecase.Confirm(ctx, id)
		api.CheckHTTPError(http.StatusInternalServerError, err)
		api.WriteJSON(w, http.StatusOK, nil)
	}
}

func NewServer(port, mongoDBName, mongoHost, mongoPlace, mongoPwd, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd string) (server *api.Server, err error) {
	fmt.Println("Starting user http routes...")

	usecase, j, err := users.Init(mongoDBName, mongoHost, mongoPlace, mongoPwd, port, jwtKeyPath, jwtSecretPath, sesAccessKey, sesAccessSecret, sesRegion, senderEmail, rpcPwd)
	if err != nil {
		log.Panic(err)
	}

	h := &handler{
		usecase: usecase,
		jwt:     j,
	}

	server = api.NewServer(port)
	r := server.Router
	r.Post("/users", h.Create())
	r.Get("/users", j.AuthorizeHandler(h.GetAll()))
	r.Get("/users/{id}", j.AuthorizeHandler(h.Get()))
	r.Put("/users/{id}", j.AuthorizeHandler(h.Update()))
	r.Delete("/users/{id}", j.AuthorizeHandler(h.Delete()))
	r.Post("/users/login", h.SignIn())
	r.Get("/users/{id}/confirm", h.Confirm())

	return
}