package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	api "github.com/contact-tracker/apiService/pkg/api/lambda"
	"github.com/contact-tracker/apiService/pkg/auth"
	"github.com/contact-tracker/apiService/places"
	t "github.com/contact-tracker/apiService/places/types"
)

type handler struct {
	usecase places.PlaceService
	jwt     *auth.JWTService
}

func isAuthorized(ctx context.Context) error {
	claims := auth.ClaimsFromContext(ctx)
	if claims == nil || claims.Subject == "" {
		return errors.New("Unauthorized")
	}
	return nil
}

func (h handler) router() func(context.Context, api.Request) (api.Response, error) {
	return func(ctx context.Context, req api.Request) (api.Response, error) {

		// Add cancellation deadline to context
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		// Add auth
		var err error
		if ctx, err = h.jwt.IncludeLambdaAuth(ctx, &req); err != nil {
			return api.Fail(err, http.StatusInternalServerError)
		}

		// Routes
		if api.MatchesRoute("/places", "POST", &req) {
			return h.Create(ctx, &req)
		} else if api.MatchesRoute("/places/{id}", "GET", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Get(ctx, &req)
		} else if api.MatchesRoute("/places", "GET", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.GetAll(ctx, &req)
		} else if api.MatchesRoute("/places/{id}", "PUT", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Update(ctx, &req)
		} else if api.MatchesRoute("/places/{id}", "DELETE", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Delete(ctx, &req)
		} else if api.MatchesRoute("/places/login", "POST", &req) {
			return h.SignIn(ctx, &req)
		} else if api.MatchesRoute("/places/{id}/confirm", "GET", &req) {
			return h.Confirm(ctx, &req)
		} else {
			return api.Fail(errors.New("not found"), http.StatusNotFound)
		}
	}
}

// Get a single place
func (h *handler) Get(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var place *t.Place
	if place, err = h.usecase.Get(ctx, id); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(place, http.StatusOK)
}

// GetAll places
func (h *handler) GetAll(ctx context.Context, _ *api.Request) (resp api.Response, err error) {
	var places []*t.Place
	if places, err = h.usecase.GetAll(ctx); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(places, http.StatusOK)
}

// Update a single place
func (h *handler) Update(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	body := []byte(req.Body)
	var update t.UpdatePlace
	if err := json.Unmarshal(body, &update); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	update.ID = id
	var place *t.Place
	if place, err = h.usecase.Update(ctx, &update); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(place, http.StatusOK)
}

// Create a place
func (h *handler) Create(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	body := []byte(r.Body)
	var req t.CreatePlace
	if err := json.Unmarshal(body, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var respPlace *t.Place
	if respPlace, err = h.usecase.Create(ctx, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(respPlace, http.StatusCreated)
}

// Delete a place
func (h *handler) Delete(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	if err = h.usecase.Delete(ctx, id); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(map[string]interface{}{"success": true}, http.StatusNoContent)
}

// SignIn a place
func (h *handler) SignIn(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	body := []byte(r.Body)
	var req t.SignInReq
	if err := json.Unmarshal(body, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var place *t.Place
	if place, err = h.usecase.SignIn(ctx, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var accessToken string
	if accessToken, err = h.jwt.GenAccessToken(place); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:  auth.AccessTokenKey,
		Value: accessToken,
		// Expires: expirationTime,
	}
	return api.SuccessWithCookie(place, http.StatusOK, cookie.String())
}

// Confirm a place
func (h *handler) Confirm(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	if err := h.usecase.Confirm(ctx, id); err != nil {
		return api.Fail(err, http.StatusUnprocessableEntity)
	}
	return api.Success(map[string]interface{}{"success": true}, http.StatusNoContent)
}

func main() {
	fmt.Println("Starting place lambda main...")
	usecase, j, err := places.Init()
	if err != nil {
		log.Panic(err)
	}
	log.Println("After init usecase...")

	h := &handler{usecase, j}
	lambda.Start(h.router())
}
