package main

import (
	"github.com/JREAMLU/j-guard/config"
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
	_, g, t := http.NewHTTPService(conf.Config)

	// init http client
	service.InitHTTPClient(t)

	g = router.GetRouters(g, conf)
	g.Run(conf.Web.URL)
}
