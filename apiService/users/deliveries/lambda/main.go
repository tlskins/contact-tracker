package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/contact-tracker/apiService/pkg/auth"
	api "github.com/contact-tracker/apiService/pkg/lambda"
	"github.com/contact-tracker/apiService/users"
	t "github.com/contact-tracker/apiService/users/types"
)

type handler struct {
	usecase users.UserService
	jwt     *auth.JWTService
}

func (handler handler) router() func(context.Context, api.Request) (api.Response, error) {
	return func(ctx context.Context, req api.Request) (api.Response, error) {

		// Add cancellation deadline to context
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		log.Println("Before router auth...")
		// Add auth
		if cookieVal, ok := req.Headers["cookie"]; ok {
			if claims, err := handler.jwt.Decode(cookieVal); err != nil {
				log.Printf("decode cookie err: %+v\n", err)
				return api.Fail(err, http.StatusInternalServerError)
			} else {
				ctx = context.WithValue(ctx, auth.AccessTokenKey, claims)
			}
		}

		log.Println("Before MatchesRoute...")
		if api.MatchesRoute("/users/{id}", "GET", &req) {
			log.Println("Before handle get...")
			return handler.Get(ctx, &req)
		}
		log.Println("After MatchesRoute...")

		switch req.HTTPMethod {
		case "POST":
			return handler.Create(ctx, []byte(req.Body))

		case "PUT":
			id, ok := req.PathParameters["id"]
			if !ok {
				return api.Response{}, errors.New("id parameter missing")
			}
			if strings.HasSuffix(req.Path, "/check_in") {
				return handler.CheckIn(ctx, id, []byte(req.Body))
			}
			return handler.Update(ctx, id, []byte(req.Body))

		case "DELETE":
			id, ok := req.PathParameters["id"]
			if !ok {
				return api.Response{}, errors.New("id parameter missing")
			}
			return handler.Delete(ctx, id)

		default:
			return api.Response{}, errors.New("route not found")
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
func (h *handler) GetAll(ctx context.Context) (resp api.Response, err error) {
	var users []*t.User
	if users, err = h.usecase.GetAll(ctx); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(users, http.StatusOK)
}

// Update a single user
func (h *handler) Update(ctx context.Context, id string, body []byte) (resp api.Response, err error) {
	var update t.UpdateUser
	if err := json.Unmarshal(body, &update); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var user *t.User
	if user, err = h.usecase.Update(ctx, &update); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(user, http.StatusOK)
}

// CheckIn a single user
func (h *handler) CheckIn(ctx context.Context, id string, body []byte) (resp api.Response, err error) {
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

// Create a user
func (h *handler) Create(ctx context.Context, body []byte) (resp api.Response, err error) {
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
func (h *handler) Delete(ctx context.Context, id string) (resp api.Response, err error) {
	if err = h.usecase.Delete(ctx, id); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
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
