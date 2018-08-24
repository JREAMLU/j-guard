package controller

import (
	"github.com/JREAMLU/j-guard/config"
	"github.com/bluele/gcache"
	jsoniter "github.com/json-iterator/go"
)

// Controller base controller
type Controller struct {
	config *config.GuardConfig
	json   jsoniter.API
	cache  gcache.Cache
}

func init() {
}
