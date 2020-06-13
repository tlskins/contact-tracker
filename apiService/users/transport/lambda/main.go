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

	l "github.com/contact-tracker/apiService/pkg/lambda"
	"github.com/contact-tracker/apiService/users"
	t "github.com/contact-tracker/apiService/users/types"
)

type handler struct {
	usecase users.UserService
}

func (handler handler) router() func(context.Context, l.Request) (l.Response, error) {
	return func(ctx context.Context, req l.Request) (l.Response, error) {

		// Add cancellation deadline to context
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		switch req.HTTPMethod {
		case "GET":
			id, ok := req.PathParameters["id"]
			if !ok {
				return handler.GetAll(ctx)
			}
			return handler.Get(ctx, id)

		case "POST":
			return handler.Create(ctx, []byte(req.Body))

		case "PUT":
			id, ok := req.PathParameters["id"]
			if !ok {
				return l.Response{}, errors.New("id parameter missing")
			}
			if strings.HasSuffix(req.Path, "/check_in") {
				return handler.CheckIn(ctx, id, []byte(req.Body))
			}
			return handler.Update(ctx, id, []byte(req.Body))

		case "DELETE":
			id, ok := req.PathParameters["id"]
			if !ok {
				return l.Response{}, errors.New("id parameter missing")
			}
			return handler.Delete(ctx, id)

		default:
			return l.Response{}, errors.New("invalid method")
		}
	}
}

// Get a single user
func (h *handler) Get(ctx context.Context, id string) (l.Response, error) {
	user, err := h.usecase.Get(ctx, id)
	if err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(user, http.StatusOK)
}

// GetAll users
func (h *handler) GetAll(ctx context.Context) (l.Response, error) {
	users, err := h.usecase.GetAll(ctx)
	if err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(users, http.StatusOK)
}

// Update a single user
func (h *handler) Update(ctx context.Context, id string, body []byte) (l.Response, error) {
	updateUser := &t.UpdateUser{}
	if err := json.Unmarshal(body, &updateUser); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	user, err := h.usecase.Update(ctx, updateUser)
	if err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(user, http.StatusOK)
}

// CheckIn a single user
func (h *handler) CheckIn(ctx context.Context, id string, body []byte) (l.Response, error) {
	req := &t.CheckInReq{}
	if err := json.Unmarshal(body, &req); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	user, err := h.usecase.CheckIn(ctx, id, req)
	if err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(user, http.StatusOK)
}

// Create a user
func (h *handler) Create(ctx context.Context, body []byte) (l.Response, error) {
	user := &t.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	var resp *t.User
	var err error
	if resp, err = h.usecase.Create(ctx, user); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(resp, http.StatusCreated)
}

// Delete a user
func (h *handler) Delete(ctx context.Context, id string) (l.Response, error) {
	if err := h.usecase.Delete(ctx, id); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func main() {
	fmt.Println("Starting user lambda main...")
	usecase, err := users.Init()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	lambda.Start(h.router())
}
