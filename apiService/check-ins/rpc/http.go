package rpc

import (
	"context"
	"fmt"

	apiHttp "github.com/contact-tracker/apiService/pkg/api/http"
	apiRpc "github.com/contact-tracker/apiService/pkg/api/rpc"
	pT "github.com/contact-tracker/apiService/places/types"
	uT "github.com/contact-tracker/apiService/users/types"
)

// RPCClient - Make RPC calls via HTTP
type RPCClient struct {
	client         *apiRpc.HTTPRPCClient
	placesHostName string
	usersHostName  string
}

// NewRPCClient - Initializes new client
func NewRPCClient(placesHostName, usersHostName, rpcPwd string) *RPCClient {
	return &RPCClient{
		client:         apiRpc.NewHTTPRPCClient(rpcPwd),
		placesHostName: placesHostName,
		usersHostName:  usersHostName,
	}
}

func (c *RPCClient) GetPlace(ctx context.Context, id string) (*pT.Place, error) {
	var place pT.Place
	fmt.Printf("%s\n", fmt.Sprintf("%s/places/%s", c.placesHostName, id))
	if code, err := c.client.HttpRequest("GET", fmt.Sprintf("%s/places/%s", c.placesHostName, id), nil, &place); err != nil {
		apiHttp.CheckHTTPError(code, err)
	}

	return &place, nil
}

func (c *RPCClient) GetUser(ctx context.Context, id string) (*uT.User, error) {
	var user uT.User
	fmt.Printf("%s\n", fmt.Sprintf("%s/users/%s", c.usersHostName, id))
	if code, err := c.client.HttpRequest("GET", fmt.Sprintf("%s/users/%s", c.usersHostName, id), nil, &user); err != nil {
		apiHttp.CheckHTTPError(code, err)
	}

	return &user, nil
}
