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
	"github.com/contact-tracker/apiService/users"
	t "github.com/contact-tracker/apiService/users/types"
)

type handler struct {
	usecase users.UserService
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
		if api.MatchesRoute("/users", "POST", &req) {
			return h.Create(ctx, &req)
		} else if api.MatchesRoute("/users/{id}", "GET", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Get(ctx, &req)
		} else if api.MatchesRoute("/users", "GET", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.GetAll(ctx, &req)
		} else if api.MatchesRoute("/users/{id}", "PUT", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Update(ctx, &req)
		} else if api.MatchesRoute("/users/{id}", "DELETE", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Delete(ctx, &req)
		} else if api.MatchesRoute("/users/login", "POST", &req) {
			return h.SignIn(ctx, &req)
		} else if api.MatchesRoute("/users/{id}/confirm", "GET", &req) {
			return h.Confirm(ctx, &req)
		} else if api.MatchesRoute("/users/{id}/check_in", "PUT", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.CheckIn(ctx, &req)
		} else if api.MatchesRoute("/users/{id}/check_out", "PUT", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.CheckOut(ctx, &req)
		} else {
			return api.Fail(errors.New("not found"), http.StatusNotFound)
		}
	}
}

// Get a single user
func (h *handler) Get(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var user *t.User
	if user, err = h.usecase.Get(ctx, id); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(user, http.StatusOK)
}

// GetAll users
func (h *handler) GetAll(ctx context.Context, _ *api.Request) (resp api.Response, err error) {
	var users []*t.User
	if users, err = h.usecase.GetAll(ctx); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(users, http.StatusOK)
}

// Update a single user
func (h *handler) Update(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	body := []byte(req.Body)
	var update t.UpdateUser
	if err := json.Unmarshal(body, &update); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	update.ID = id
	var user *t.User
	if user, err = h.usecase.Update(ctx, &update); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(user, http.StatusOK)
}

// CheckIn a single user
func (h *handler) CheckIn(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", r); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	body := []byte(r.Body)
	var req t.CheckInReq
	if err := json.Unmarshal(body, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var user *t.User
	if user, err = h.usecase.CheckIn(ctx, id, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(user, http.StatusOK)
}

// CheckOut a single user
func (h *handler) CheckOut(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", r); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	body := []byte(r.Body)
	var req t.CheckOutReq
	if err := json.Unmarshal(body, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var user *t.User
	if user, err = h.usecase.CheckOut(ctx, id, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(user, http.StatusOK)
}

// Create a user
func (h *handler) Create(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	body := []byte(r.Body)
	var req t.CreateUser
	if err := json.Unmarshal(body, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var respUser *t.User
	if respUser, err = h.usecase.Create(ctx, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(respUser, http.StatusCreated)
}

// Delete a user
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

// SignIn a user
func (h *handler) SignIn(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	body := []byte(r.Body)
	var req t.SignInReq
	if err := json.Unmarshal(body, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var user *t.User
	if user, err = h.usecase.SignIn(ctx, &req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var accessToken string
	if accessToken, err = h.jwt.GenAccessToken(user); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:  auth.AccessTokenKey,
		Value: accessToken,
		// Expires: expirationTime,
	}
	return api.SuccessWithCookie(user, http.StatusOK, cookie.String())
}

// Confirm a user
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
	fmt.Println("Starting user lambda main...")
	usecase, j, err := users.Init()
	if err != nil {
		log.Panic(err)
	}
	log.Println("After init usecase...")

	h := &handler{usecase, j}
	lambda.Start(h.router())
}
