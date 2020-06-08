package http

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		panic(Error{http.StatusInternalServerError, err})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
