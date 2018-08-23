package main

import (
	"github.com/JREAMLU/j-guard/config"
	"github.com/JREAMLU/j-guard/middleware"
	"github.com/JREAMLU/j-guard/router"
	"github.com/JREAMLU/j-guard/service"
	"github.com/JREAMLU/j-kit/http"
)

func main() {
	// load config
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	RunHTTPService(conf)
}

// RunHTTPService run http service
func RunHTTPService(conf *config.GuardConfig) {
	ms, g, t := http.NewHTTPService(conf.Config)
	// @TODO InitGCache
	g.Use(middleware.Middle())

	// init http client
	service.InitHTTPClient(t)

	// init micro client
	service.InitMicroClient(ms)

	g = router.GetRouters(g, conf)
	g.Run(conf.Web.URL)
}
