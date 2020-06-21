package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

type Error struct {
	Code  int
	Error interface{}
}

func (e Error) String() string {
	err := e.Error
	if bErr, ok := err.([]byte); ok {
		return string(bErr)
	} else if sErr, ok := err.(string); ok {
		return sErr
	} else if eErr, ok := err.(error); ok {
		return eErr.Error()
	} else if bErr, err := json.Marshal(err); err != nil {
		return string(bErr)
	}
	return ""
}

// @param err of type []byte, string, error, or json serializable object
func Abort(code int, err interface{}) {
	panic(Error{code, err})
}

func CheckHTTPError(statusCode int, err error) {
	if err != nil {
		Abort(statusCode, err)
	}
}

func HandleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		fmt.Println(r)
		debug.PrintStack()
		if err, ok := r.(Error); ok {
			WriteJSON(w, err.Code, map[string]interface{}{"message": err.String()})
		} else if err, ok := r.(error); ok {
			WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
		} else {
			WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "unknown error"})
		}
	}
}
