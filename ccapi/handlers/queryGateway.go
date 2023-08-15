package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/ccapi/chaincode"
	"github.com/goledgerdev/ccapi/common"
)

func QueryGatewayDefault(c *gin.Context) {
	channelName := os.Getenv("CHANNEL")
	chaincodeName := os.Getenv("CCNAME")

	queryGateway(c, channelName, chaincodeName)
}

func QueryGatewayCustom(c *gin.Context) {
	channelName := c.Param("channelName")
	chaincodeName := c.Param("chaincodeName")

	queryGateway(c, channelName, chaincodeName)
}

func queryGateway(c *gin.Context, channelName, chaincodeName string) {
	var args []byte
	var err error

	// Get request data
	if c.Request.Method == "GET" {
		request := c.Query("@request")
		if request != "" {
			args, _ = base64.StdEncoding.DecodeString(request)
		}
	} else if c.Request.Method == "POST" {
		req := make(map[string]interface{})
		c.ShouldBind(&req)
		args, err = json.Marshal(req)
		if err != nil {
			common.Abort(c, http.StatusInternalServerError, err)
			return
		}
	}

	txName := c.Param("txname")

	// Query
	result, err := chaincode.QueryGateway(channelName, chaincodeName, txName, string(args))
	if err != nil {
		err, status := common.ParseError(err)
		common.Abort(c, status, err)
		return
	}

	// Parse response
	var payload interface{}
	err = json.Unmarshal(result, &payload)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	common.Respond(c, payload, http.StatusOK, nil)
}
