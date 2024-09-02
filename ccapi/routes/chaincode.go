package routes

import (
	"os"

	"github.com/hyperledger-labs/ccapi/handlers"

	"github.com/gin-gonic/gin"
)

func addCCRoutes(rg *gin.RouterGroup) {
	if os.Getenv("FPC_MODE") == "true" {
		//Use FPC Handlers
		rg.POST("/:channelName/:chaincodeName/invoke/:txname", handlers.InvokeFpc)
		rg.PUT("/:channelName/:chaincodeName/invoke/:txname", handlers.InvokeFpc)
		rg.DELETE("/:channelName/:chaincodeName/invoke/:txname", handlers.InvokeFpc)
		rg.POST("/:channelName/:chaincodeName/query/:txname", handlers.QueryFpc)
		rg.GET("/:channelName/:chaincodeName/query/:txname", handlers.QueryFpc)

		rg.POST("/invoke/:txname/", handlers.InvokeFpc)
		rg.POST("/invoke/:txname", handlers.InvokeFpc)
		rg.PUT("/invoke/:txname/", handlers.InvokeFpc)
		rg.PUT("/invoke/:txname", handlers.InvokeFpc)
		rg.DELETE("/invoke/:txname/", handlers.InvokeFpc)
		rg.DELETE("/invoke/:txname", handlers.InvokeFpc)
		rg.POST("/query/:txname/", handlers.QueryFpc)
		rg.POST("/query/:txname", handlers.QueryFpc)
		rg.GET("/query/:txname/", handlers.QueryFpc)
		rg.GET("/query/:txname", handlers.QueryFpc)

		rg.GET("/:channelName/qscc/:txname", handlers.QueryQSCC)

	} else {
		//Use Fabric Handlers
		// Gateway routes
		rg.POST("/gateway/:channelName/:chaincodeName/invoke/:txname", handlers.InvokeGatewayCustom)
		rg.PUT("/gateway/:channelName/:chaincodeName/invoke/:txname", handlers.InvokeGatewayCustom)
		rg.DELETE("/gateway/:channelName/:chaincodeName/invoke/:txname", handlers.InvokeGatewayCustom)
		rg.POST("/gateway/:channelName/:chaincodeName/query/:txname", handlers.QueryGatewayCustom)
		rg.GET("/gateway/:channelName/:chaincodeName/query/:txname", handlers.QueryGatewayCustom)

		rg.POST("/gateway/invoke/:txname", handlers.InvokeGatewayDefault)
		rg.PUT("/gateway/invoke/:txname", handlers.InvokeGatewayDefault)
		rg.DELETE("/gateway/invoke/:txname", handlers.InvokeGatewayDefault)
		rg.POST("/gateway/query/:txname", handlers.QueryGatewayDefault)
		rg.GET("/gateway/query/:txname", handlers.QueryGatewayDefault)

		// Other
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

		rg.GET("/:channelName/qscc/:txname", handlers.QueryQSCC)
	}
}
