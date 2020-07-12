package rpc

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPRPCClient - Make RPC calls via HTTP
type HTTPRPCClient struct {
	c      *http.Client
	rpcPwd string
}

type gzreadCloser struct {
	*gzip.Reader
	io.Closer
}

func (gz gzreadCloser) Close() error {
	return gz.Closer.Close()
}

// NewHTTPRPCClient - Initializes new client
func NewHTTPRPCClient(rpcPwd string) *HTTPRPCClient {
	return &HTTPRPCClient{
		&http.Client{
			Transport: &http.Transport{
				IdleConnTimeout:    60 * time.Second,
				DisableCompression: true,
			},
			Timeout: 60 * time.Second,
		},
		rpcPwd,
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
	req.Header.Set("Authorization", c.rpcPwd)

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
