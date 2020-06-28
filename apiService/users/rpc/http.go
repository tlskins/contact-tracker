package rpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/contact-tracker/apiService/pkg/auth"
	api "github.com/contact-tracker/apiService/pkg/http"
	pT "github.com/contact-tracker/apiService/places/types"
)

// HTTPRPCClient - Make RPC calls via HTTP
type HTTPRPCClient struct {
	c              *http.Client
	placesHostName string
}

type gzreadCloser struct {
	*gzip.Reader
	io.Closer
}

func (gz gzreadCloser) Close() error {
	return gz.Closer.Close()
}

// NewHTTPRPCClient - Initializes new client
func NewHTTPRPCClient(placesHostName string) *HTTPRPCClient {
	return &HTTPRPCClient{
		&http.Client{
			Transport: &http.Transport{
				IdleConnTimeout:    60 * time.Second,
				DisableCompression: true,
			},
			Timeout: 60 * time.Second,
		},
		placesHostName,
	}
}

func (c *HTTPRPCClient) HttpRequest(method, path string, data, out interface{}) (int, error) {
	var req *http.Request
	var err error
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		req, err = http.NewRequest(method, path, bytes.NewReader(b))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  auth.RPCAccessTokenKey,
		Value: "ABC",
	})

	resp, err := c.c.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status >= 400 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return status, err
		}
		return status, errors.New(string(b))
	}
	if out != nil {
		var b []byte
		if resp.Header.Get("Content-Encoding") == "gzip" {
			resp.Header.Del("Content-Length")
			zr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			resp.Body = gzreadCloser{zr, resp.Body}
		}
		if b, err = ioutil.ReadAll(resp.Body); err != nil {
			return status, err
		}
		if err := json.Unmarshal(b, out); err != nil {
			return status, err
		}
	}
	return status, nil
}

func (c *HTTPRPCClient) GetPlace(ctx context.Context, id string) (*pT.Place, error) {
	var place pT.Place
	if code, err := c.HttpRequest("GET", fmt.Sprintf("%s/places/%s", c.placesHostName, id), nil, &place); err != nil {
		api.CheckHTTPError(code, err)
	}

	return &place, nil

	// resp, err := c.c.Get(fmt.Sprintf("%s/places/%s", c.placesHostName, id))
	// api.CheckHTTPError(resp.StatusCode, err)
	// b, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// &
	// if api.IsErrorStatusCode(resp.StatusCode) {
	// 	return nil, fmt.Errorf(string(b))
	// }
	// var place pT.Place
	// err = json.Unmarshal(b, &place)
	// return &place, err
}
