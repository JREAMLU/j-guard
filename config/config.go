package config

import (
	"github.com/JREAMLU/j-kit/go-micro/util"
)

const (
	name    = "guard"
	version = "v1"
)

// GuardConfig guard config
type GuardConfig struct {
	*util.Config

	Guard struct {
		Timeout int64
	}

	Cache struct {
		Expire int64
		Size   int
	}
}

// Load load config
func Load() (*GuardConfig, error) {
	// load redis mysql elastic client

	// load parent config
	config := &GuardConfig{}
	err := util.LoadCustomConfig("10.200.202.35:8500", name, version, config)
	if err != nil {
		return nil, err
	}

	// set default
	if config.Guard.Timeout == 0 {
		config.Guard.Timeout = 3
	}

	return config, err
}
