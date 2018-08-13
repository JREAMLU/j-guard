package service

import (
	"context"
	"encoding/json"

	"github.com/JREAMLU/j-kit/constant"
	"github.com/micro/go-micro/client"
	grpcClient "github.com/micro/go-plugins/client/grpc"
)

// Request grpc request
func Request(ctx context.Context, logo, service, method, address string, request interface{}) (string, error) {
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
		return constant.EmptyStr, err
	}

	resp, err := response.MarshalJSON()
	if err != nil {
		return constant.EmptyStr, err
	}

	return string(resp), nil
}
