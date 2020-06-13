package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/contact-tracker/apiService/places"
	t "github.com/contact-tracker/apiService/places/types"

	"github.com/gorilla/mux"
)

const fiveSecondsTimeout = time.Second * 5

type handler struct {
	usecase places.PlaceService
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

	place, err := d.usecase.Get(ctx, id)
	if err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(place)
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

	places, err := d.usecase.GetAll(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}

	data, err := json.Marshal(places)
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
	place := &t.UpdatePlace{}
	if err := decoder.Decode(&place); err != nil {
		writeErr(w, err)
		return
	}

	vars := mux.Vars(r)
	place.ID = vars["id"]

	var err error
	var usr *t.Place
	if usr, err = d.usecase.Update(ctx, place); err != nil {
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

	place := &t.Place{}
	var err error
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&place); err != nil {
		writeErr(w, err)
		return
	}

	var resp *t.Place
	if resp, err = d.usecase.Create(ctx, place); err != nil {
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
	fmt.Println("Starting place http routes...")
	usecase, err := places.InitMongoService()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	r := mux.NewRouter()
	r.HandleFunc("/places", h.Create).Methods("POST")
	r.HandleFunc("/places", h.GetAll).Methods("GET")
	r.HandleFunc("/places/{id}", h.Get).Methods("GET")
	r.HandleFunc("/places/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/places/{id}", h.Delete).Methods("DELETE")

	return r, nil
}
