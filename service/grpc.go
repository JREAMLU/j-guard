package service

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/client"
)

// Request grpc request
func Request(ctx context.Context, logo, service, method, address string, request interface{}) ([]byte, error) {
	req := microClient.NewRequest(service, method, request, client.WithContentType("application/json"))

	var response json.RawMessage
	var err error

	if len(address) > 0 {
		err = microClient.Call(ctx, req, &response, client.WithAddress(address))
	} else {
		err = microClient.Call(ctx, req, &response)
	}
	if err != nil {
		return nil, err
	}

	resp, err := response.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
