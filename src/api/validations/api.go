package api

import (
	"automart/api/routers"
	"automart/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitServer() {
	cfg := config.GetConfig()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		health := v1.Group("/health")
		testRouter := v1.Group("test")
		routers.Health(health)
		routers.TestRouter(testRouter)
	}

	v2 := api.Group("/v2")
	{
		health := v2.Group("/health")
		routers.Health(health)
	}

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
