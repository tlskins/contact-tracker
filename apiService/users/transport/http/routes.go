package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/contact-tracker/apiService/users"
	t "github.com/contact-tracker/apiService/users/types"

	"github.com/gorilla/mux"
)

const fiveSecondsTimeout = time.Second * 5

type handler struct {
	usecase users.UserService
}

func writeErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (d *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	user, err := d.usecase.Get(ctx, id)
	if err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (d *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	users, err := d.usecase.GetAll(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (d *handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	user := &t.UpdateUser{}
	if err := decoder.Decode(&user); err != nil {
		writeErr(w, err)
		return
	}

	vars := mux.Vars(r)
	user.ID = vars["id"]

	var err error
	var usr *t.User
	if usr, err = d.usecase.Update(ctx, user); err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(usr)
	if err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (d *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	user := &t.User{}
	var err error
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&user); err != nil {
		writeErr(w, err)
		return
	}

	var resp *t.User
	if resp, err = d.usecase.Create(ctx, user); err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (d *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	if err := d.usecase.Delete(ctx, id); err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Deleted"))
}

// Routes -
func Routes() (*mux.Router, error) {
	fmt.Println("Starting user http routes...")
	usecase, err := users.InitMongoService()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	r := mux.NewRouter()
	r.HandleFunc("/users", h.Create).Methods("POST")
	r.HandleFunc("/users", h.GetAll).Methods("GET")
	r.HandleFunc("/users/{id}", h.Get).Methods("GET")
	r.HandleFunc("/users/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/users/{id}", h.Delete).Methods("DELETE")

	return r, nil
}
