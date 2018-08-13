package config

import "github.com/JREAMLU/j-kit/go-micro/util"

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
}

// Load load config
func Load() (*GuardConfig, error) {
	// load redis mysql elastic client

	// load parent config
	config := &GuardConfig{}
	err := util.LoadCustomConfig("10.200.202.35:8500", name, version, config)

	return config, err
}
