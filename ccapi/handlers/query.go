package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/ccapi/chaincode"
	"github.com/goledgerdev/ccapi/common"
)

func Query(c *gin.Context) {
	// Get channel information from request
	req := make(map[string]interface{})
	c.ShouldBind(&req)

	channelName := c.Param("channelName")
	chaincodeName := c.Param("chaincodeName")
	txName := c.Param("txname")

	args, err := json.Marshal(req)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	res, status, err := chaincode.Query(channelName, chaincodeName, txName, [][]byte{args})
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	var payload interface{}
	err = json.Unmarshal(res.Payload, &payload)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	common.Respond(c, payload, status, err)
}
