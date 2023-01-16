package routes

import (
	"github.com/goledgerdev/ccapi/handlers"

	"github.com/gin-gonic/gin"
)

func addCCRoutes(rg *gin.RouterGroup) {
	rg.POST("/:channelName/:chaincodeName/invoke/:txname", handlers.Invoke)
	rg.PUT("/:channelName/:chaincodeName/invoke/:txname", handlers.Invoke)
	rg.DELETE("/:channelName/:chaincodeName/invoke/:txname", handlers.Invoke)
	rg.POST("/:channelName/:chaincodeName/query/:txname", handlers.Query)
	rg.GET("/:channelName/:chaincodeName/query/:txname", handlers.Query)

	rg.POST("/invoke/:txname/", handlers.InvokeV1)
	rg.POST("/invoke/:txname", handlers.InvokeV1)
	rg.PUT("/invoke/:txname/", handlers.InvokeV1)
	rg.PUT("/invoke/:txname", handlers.InvokeV1)
	rg.DELETE("/invoke/:txname/", handlers.InvokeV1)
	rg.DELETE("/invoke/:txname", handlers.InvokeV1)
	rg.POST("/query/:txname/", handlers.QueryV1)
	rg.POST("/query/:txname", handlers.QueryV1)
	rg.GET("/query/:txname/", handlers.QueryV1)
	rg.GET("/query/:txname", handlers.QueryV1)
}
