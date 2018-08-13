package router

import (
	"github.com/JREAMLU/j-guard/config"
	"github.com/JREAMLU/j-guard/controller"

	"github.com/gin-gonic/gin"
)

// GetRouters init router
func GetRouters(router *gin.Engine, conf *config.GuardConfig) *gin.Engine {
	guard := controller.NewGuardController(conf)

	// grpc
	router.POST("/grpc", guard.Grpc)

	return router
}
