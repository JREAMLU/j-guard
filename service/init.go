package service

import (
	"github.com/JREAMLU/j-kit/http"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	httpClient   *http.Requests
	microService micro.Service
	microClient  client.Client
)

// InitHTTPClient init http client
func InitHTTPClient(tracer opentracing.Tracer) {
	httpClient = http.NewRequests(tracer)
}

// InitMicroClient init micro service
func InitMicroClient(service micro.Service) {
	microClient = service.Client()
}
