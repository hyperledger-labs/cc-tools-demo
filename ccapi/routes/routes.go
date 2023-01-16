package routes

import (
	"github.com/gin-gonic/gin"
)

// Register routes and handlers used by engine
func AddRoutesToEngine(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// CHANNEL routes
	chaincodeRG := r.Group("/api")
	addCCRoutes(chaincodeRG)

	// Update SDK route
	sdkRG := r.Group("/sdk")
	addSDKRoutes(sdkRG)

}
