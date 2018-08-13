package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JREAMLU/j-guard/config"
	"github.com/JREAMLU/j-guard/service"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// GuardController guard controller
type GuardController struct {
	Controller
}

// GrpcReq grpc request
type GrpcReq struct {
	Req []struct {
		Logo    string      `json:"logo"`
		Service string      `json:"service"`
		Address string      `json:"address"`
		Method  string      `json:"method"`
		Request interface{} `json:"request"`
	} `json:"req"`
}

// NewGuardController new hello
func NewGuardController(conf *config.GuardConfig) *GuardController {
	return &GuardController{
		Controller{
			config: conf,
			json:   jsoniter.ConfigCompatibleWithStandardLibrary,
		},
	}
}

// Grpc gateway grpc
func (g *GuardController) Grpc(c *gin.Context) {
	var reqs GrpcReq
	raw, err := c.GetRawData()
	if err != nil {
		fmt.Println("++++++++++++: err", err)
		return
	}

	err = g.json.Unmarshal(raw, &reqs)
	if err != nil {
		fmt.Println("++++++++++++: err", err)
		return
	}

	grpcRequest(c.Request.Context(), reqs)

	var resp struct {
		Name string
	}
	resp.Name = "LUj"

	c.JSON(http.StatusOK, resp)
}

// Respone respone
type Respone struct {
	Logo  string
	Resp  string
	Error error
}

func grpcRequest(ctx context.Context, reqs GrpcReq) {
	respChan := make(chan *Respone, len(reqs.Req))

	// service.Request(ctx, req.Logo, req.Service, req.Method, req.Address, req.Request)
	for _, req := range reqs.Req {
		go func(ctx context.Context, logo, serviceName, method, address string, request interface{}) {
			resp, err := service.Request(ctx, logo, serviceName, method, address, request)
			respChan <- &Respone{
				Logo:  logo,
				Resp:  resp,
				Error: err,
			}
		}(ctx, req.Logo, req.Service, req.Method, req.Address, req.Request)
	}

	for i := 0; i < len(reqs.Req); i++ {
		select {
		case resp := <-respChan:
			fmt.Println("++++++++++++:resp ", resp)
		case <-time.After(1000 * time.Millisecond):
			fmt.Println("++++++++++++: 超时")
			return
		}
	}
}
