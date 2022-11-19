package routes

import (
	"github.com/goledgerdev/ccapi/handlers"

	"github.com/gin-gonic/gin"
)

func addCCRoutes(rg *gin.RouterGroup) {
	rg.POST(":channelName/:chaincodeName/invoke/:txname", handlers.Invoke)
	rg.POST(":channelName/:chaincodeName/query/:txname", handlers.Query)
}
