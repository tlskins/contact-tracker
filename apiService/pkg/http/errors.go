package http

import (
	"encoding/json"
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
