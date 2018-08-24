package controller

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/JREAMLU/j-guard/config"
	"github.com/JREAMLU/j-guard/constant"
	"github.com/JREAMLU/j-guard/service"
	"github.com/bluele/gcache"

	"github.com/JREAMLU/j-kit/crypto"
	"github.com/JREAMLU/j-kit/go-micro/util"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

const gcacheExpire = "Gcache-Expire"

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
	Data       map[string]interface{}
	StatusCode int64
	Message    string
}

// NewGuardController new hello
func NewGuardController(conf *config.GuardConfig) *GuardController {
	return &GuardController{
		Controller{
			config: conf,
			json:   jsoniter.ConfigCompatibleWithStandardLibrary,
			cache:  gcache.New(conf.Cache.Size).LRU().Build(),
		},
	}
}

// Grpc gateway grpc
func (g *GuardController) Grpc(c *gin.Context) {
	var reqs GrpcReq
	resps := &Respones{}

	raw, err := c.GetRawData()
	if err != nil {
		resps.Message = err.Error()
		resps.StatusCode = constant.SystemErrorCode
		c.JSON(http.StatusBadRequest, resps)
		return
	}

	// GetRawData buffed
	if len(raw) == 0 {
		raw = c.MustGet("raw").([]byte)
	}

	cacheRes, err := g.getGCache(c.Request.Context(), string(raw))
	if err == nil {
		c.JSON(http.StatusOK, cacheRes.(*Respones))
		return
	}

	err = g.json.Unmarshal(raw, &reqs)
	if err != nil {
		resps.Message = err.Error()
		resps.StatusCode = constant.SystemErrorCode
		c.JSON(http.StatusBadRequest, resps)
		return
	}

	res := g.grpcRequest(c.Request.Context(), reqs)
	resps.Data = res

	g.setGCache(c.Request.Context(), string(raw), resps, c.Request.Header.Get(gcacheExpire))

	c.JSON(http.StatusOK, resps)
}

func (g *GuardController) getGCache(ctx context.Context, raw string) (interface{}, error) {
	hashKey, err := crypto.MD5(raw, true)
	if err != nil {
		return nil, err
	}

	val, err := g.cache.Get(hashKey)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// @TODO singleflight
func (g *GuardController) setGCache(ctx context.Context, raw string, resp *Respones, expire string) {
	if expire == "" {
		return
	}

	expireInt64, err := strconv.ParseInt(expire, 10, 64)
	if err != nil || expireInt64 == 0 {
		return
	}

	hashKey, err := crypto.MD5(raw, true)
	if err != nil {
		return
	}

	g.cache.SetWithExpire(hashKey, resp, time.Duration(expireInt64)*time.Second)
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
				util.TraceLog(ctx, resp.Error.Error())
				return resps
			}
			err := g.json.Unmarshal(resp.Resp, &rsp)
			if err != nil {
				util.TraceLog(ctx, err.Error())
				return resps
			}
			mutex.Lock()
			resps[resp.Logo] = rsp
			mutex.Unlock()

		case <-time.After(time.Duration(g.config.Guard.Timeout) * time.Millisecond):
			util.TraceLog(ctx, constant.GrpcConcurrentTimeout)
			return resps
		}
	}

	return resps
}
