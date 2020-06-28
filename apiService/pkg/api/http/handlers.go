package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func ParseHTTPParams(r *http.Request, out interface{}) {
	if r.Method == "GET" {
		if err := decoder.Decode(out, r.URL.Query()); err != nil {
			Abort(http.StatusUnprocessableEntity, err)
		}
	} else {
		var b bytes.Buffer
		bodyCpy := io.TeeReader(r.Body, &b)
		if err := json.NewDecoder(bodyCpy).Decode(out); err != nil {
			Abort(http.StatusUnprocessableEntity, err)
		}
		r.Body = ioutil.NopCloser(&b)
	}
}

func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		panic(Error{http.StatusInternalServerError, err})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RealIP,
		middleware.RequestID,
		middleware.Timeout(60*time.Second),
		Recoverer,
	)

	return r
}

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer HandleError(w)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
