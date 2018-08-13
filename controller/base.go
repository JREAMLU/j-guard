package controller

import (
	"github.com/JREAMLU/j-guard/config"
	jsoniter "github.com/json-iterator/go"
)

// Controller base controller
type Controller struct {
	config *config.GuardConfig
	json   jsoniter.API
}

func init() {
}
