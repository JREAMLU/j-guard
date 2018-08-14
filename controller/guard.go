package controller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
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

// Respone respone
type Respone struct {
	Logo  string
	Resp  []byte
	Error error
}

// Respones respones
type Respones struct {
	Resp map[string]interface{}
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

	res := g.grpcRequest(c.Request.Context(), reqs)

	resps := &Respones{
		Resp: res,
	}

	c.JSON(http.StatusOK, resps)
}

func (g *GuardController) grpcRequest(ctx context.Context, reqs GrpcReq) map[string]interface{} {
	respChan := make(chan *Respone, len(reqs.Req))
	resps := make(map[string]interface{})
	var rsp interface{}
	var mutex sync.Mutex

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
			if resp.Error != nil {
				// @TODO log trace
				return resps
			}
			err := g.json.Unmarshal(resp.Resp, &rsp)
			if err != nil {
				// @TODO log trace
				return resps
			}
			mutex.Lock()
			resps[resp.Logo] = rsp
			mutex.Unlock()

		case <-time.After(time.Duration(g.config.Guard.Timeout) * time.Millisecond):
			// @TODO log trace
			return resps
		}
	}

	return resps
}
