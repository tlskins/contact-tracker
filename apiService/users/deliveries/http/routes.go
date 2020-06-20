package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/contact-tracker/apiService/pkg/auth"
	api "github.com/contact-tracker/apiService/pkg/http"
	"github.com/contact-tracker/apiService/users"
	t "github.com/contact-tracker/apiService/users/types"

	"github.com/go-chi/chi"
)

const fiveSecondsTimeout = time.Second * 5

type handler struct {
	usecase users.UserService
}

func (d *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	id := chi.URLParam(r, "id")
	user, err := d.usecase.Get(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, user)
}

func (d *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	// AUTH //

	c, err := r.Cookie(auth.AccessTokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			api.CheckError(http.StatusUnauthorized, fmt.Errorf("Unauthorized access"))
		}
		api.CheckError(http.StatusUnauthorized, err)
	}

	tokenStr := c.Value
	claims, err := d.usecase.Decode(ctx, tokenStr)
	ctx = context.WithValue(ctx, auth.AccessTokenKey, claims)

	// AUTH //

	users, err := d.usecase.GetAll(ctx)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, users)
}

func (d *handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	req := &t.UpdateUser{}
	api.ParseHTTPParams(r, req)

	req.ID = chi.URLParam(r, "id")
	resp, err := d.usecase.Update(ctx, req)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, resp)
}

func (d *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	req := &t.SignInReq{}
	api.ParseHTTPParams(r, req)

	user, accessToken, err := d.usecase.SignIn(ctx, req)
	api.CheckError(http.StatusInternalServerError, err)

	http.SetCookie(w, &http.Cookie{
		Name:  auth.AccessTokenKey,
		Value: accessToken,
		// Expires: expirationTime,
	})

	api.WriteJSON(w, http.StatusOK, user)
}

func (d *handler) CheckIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	req := &t.CheckInReq{}
	api.ParseHTTPParams(r, req)

	id := chi.URLParam(r, "id")
	usr, err := d.usecase.CheckIn(ctx, id, req)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, usr)
}

func (d *handler) CheckOut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	req := &t.CheckOutReq{}
	api.ParseHTTPParams(r, req)
	id := chi.URLParam(r, "id")

	usr, err := d.usecase.CheckOut(ctx, id, req)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, usr)
}

func (d *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	req := &t.CreateUser{}
	api.ParseHTTPParams(r, req)

	resp, err := d.usecase.Create(ctx, req)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, resp)
}

func (d *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	id := chi.URLParam(r, "id")

	err := d.usecase.Delete(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	api.WriteJSON(w, http.StatusOK, nil)
}

// Routes -
func Routes() (*chi.Mux, error) {
	fmt.Println("Starting user http routes...")
	usecase, err := users.Init()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	r := api.NewRouter()

	r.Post("/users", h.Create)
	r.Get("/users", h.GetAll)
	r.Get("/users/{id}", h.Get)
	r.Put("/users/{id}", h.Update)
	r.Put("/users/{id}/check_in", h.CheckIn)
	r.Put("/users/{id}/check_out", h.CheckOut)
	r.Delete("/users/{id}", h.Delete)
	r.Post("/users/login", h.SignIn)

	return r, nil
}
