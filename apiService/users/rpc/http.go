package rpc

import (
	"context"
	"fmt"

	apiHttp "github.com/contact-tracker/apiService/pkg/api/http"
	apiRpc "github.com/contact-tracker/apiService/pkg/api/rpc"
	pT "github.com/contact-tracker/apiService/places/types"
)

// RPCClient - Make RPC calls via HTTP
type RPCClient struct {
	client         *apiRpc.HTTPRPCClient
	placesHostName string
}

// NewHTTPRPCClient - Initializes new client
func NewRPCClient(placesHostName, rpcPwd string) *RPCClient {
	return &RPCClient{
		client:         apiRpc.NewHTTPRPCClient(rpcPwd),
		placesHostName: placesHostName,
	}
}

func (c *RPCClient) GetPlace(ctx context.Context, id string) (*pT.Place, error) {
	var place pT.Place
	if code, err := c.client.HttpRequest("GET", fmt.Sprintf("%s/places/%s", c.placesHostName, id), nil, &place); err != nil {
		apiHttp.CheckHTTPError(code, err)
	}

	return &place, nil
}
