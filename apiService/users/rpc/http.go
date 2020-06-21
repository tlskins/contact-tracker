package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	api "github.com/contact-tracker/apiService/pkg/http"
	pT "github.com/contact-tracker/apiService/places/types"
)

// HTTPRPCClient - Make RPC calls via HTTP
type HTTPRPCClient struct {
	c              *http.Client
	placesHostName string
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

func (c *HTTPRPCClient) GetPlace(ctx context.Context, id string) (*pT.Place, error) {
	resp, err := c.c.Get(fmt.Sprintf("%s/places/%s", c.placesHostName, id))
	api.CheckHTTPError(resp.StatusCode, err)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var place pT.Place
	err = json.Unmarshal(b, &place)
	return &place, err
}
