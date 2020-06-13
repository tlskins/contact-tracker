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

	l "github.com/contact-tracker/apiService/pkg/lambda"
	"github.com/contact-tracker/apiService/places"
	t "github.com/contact-tracker/apiService/places/types"
)

type handler struct {
	usecase places.PlaceService
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

// Get a single place
func (h *handler) Get(ctx context.Context, id string) (l.Response, error) {
	place, err := h.usecase.Get(ctx, id)
	if err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(place, http.StatusOK)
}

// GetAll places
func (h *handler) GetAll(ctx context.Context) (l.Response, error) {
	places, err := h.usecase.GetAll(ctx)
	if err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(places, http.StatusOK)
}

// Update a single place
func (h *handler) Update(ctx context.Context, id string, body []byte) (l.Response, error) {
	updatePlace := &t.UpdatePlace{}
	if err := json.Unmarshal(body, &updatePlace); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	var place *t.Place
	var err error
	if place, err = h.usecase.Update(ctx, updatePlace); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(place, http.StatusOK)
}

// Create a place
func (h *handler) Create(ctx context.Context, body []byte) (l.Response, error) {
	place := &t.Place{}
	if err := json.Unmarshal(body, &place); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	var resp *t.Place
	var err error
	if resp, err = h.usecase.Create(ctx, place); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(resp, http.StatusCreated)
}

// Delete a place
func (h *handler) Delete(ctx context.Context, id string) (l.Response, error) {
	if err := h.usecase.Delete(ctx, id); err != nil {
		return l.Fail(err, http.StatusInternalServerError)
	}

	return l.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func main() {
	fmt.Println("Starting place lambda main...")
	usecase, err := places.InitMongoService()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	lambda.Start(h.router())
}
