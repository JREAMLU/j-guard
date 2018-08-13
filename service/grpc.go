package service

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/client"
	grpcClient "github.com/micro/go-plugins/client/grpc"
)

// Request grpc request
func Request(ctx context.Context, logo, service, method, address string, request interface{}) ([]byte, error) {
	grpc := grpcClient.NewClient()
	req := grpc.NewRequest(service, method, request, client.WithContentType("application/json"))

	var response json.RawMessage
	var err error

	if len(address) > 0 {
		err = grpc.Call(ctx, req, &response, client.WithAddress(address))
	} else {
		err = grpc.Call(ctx, req, &response)
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
