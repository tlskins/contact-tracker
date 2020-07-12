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
	"github.com/gorilla/schema"

	chk "github.com/contact-tracker/apiService/check-ins"
	t "github.com/contact-tracker/apiService/check-ins/types"
	api "github.com/contact-tracker/apiService/pkg/api/lambda"
	"github.com/contact-tracker/apiService/pkg/auth"
)

var decoder = schema.NewDecoder()

type handler struct {
	usecase chk.CheckInService
	jwt     *auth.JWTService
}

func isAuthorized(ctx context.Context) error {
	if authorized, _ := auth.ClaimsFromContext(ctx); !authorized {
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
		if api.MatchesRoute("/check-ins/{id}", "GET", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.Get(ctx, &req)
		} else if api.MatchesRoute("/check-ins", "GET", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.GetAll(ctx, &req)
		} else if api.MatchesRoute("/check-ins", "POST", &req) {
			if err := isAuthorized(ctx); err != nil {
				return api.Fail(err, http.StatusUnauthorized)
			}
			return h.CheckIn(ctx, &req)
		} else {
			return api.Fail(errors.New("not found"), http.StatusNotFound)
		}
	}
}

// Get a check in
func (h *handler) Get(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var id string
	if id, err = api.GetPathParam("id", req); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var checkIn *t.CheckIn
	if checkIn, err = h.usecase.Get(ctx, id); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(checkIn, http.StatusOK)
}

// GetAll check ins
func (h *handler) GetAll(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	var query t.GetCheckIns
	parsed := make(map[string][]string)
	for k, v := range req.QueryStringParameters {
		parsed[k] = []string{v}
	}
	if err := decoder.Decode(query, parsed); err != nil {
		api.Fail(err, http.StatusInternalServerError)
	}
	var checkIns []*t.CheckIn
	if checkIns, err = h.usecase.GetAll(ctx, &query); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(checkIns, http.StatusOK)
}

// Check in user
func (h *handler) CheckIn(ctx context.Context, req *api.Request) (resp api.Response, err error) {
	body := []byte(req.Body)
	var query t.CreateCheckIn
	if err := json.Unmarshal(body, &query); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	var checkIn *t.CheckIn
	if checkIn, err = h.usecase.CheckIn(ctx, &query); err != nil {
		return api.Fail(err, http.StatusInternalServerError)
	}
	return api.Success(checkIn, http.StatusOK)
}

func main() {
	fmt.Println("Starting chekc ins lambda main...")
	usecase, j, err := chk.Init()
	if err != nil {
		log.Panic(err)
	}
	log.Println("After init usecase...")

	h := &handler{usecase, j}
	lambda.Start(h.router())
}
