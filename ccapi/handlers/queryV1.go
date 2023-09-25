package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger-labs/ccapi/chaincode"
	"github.com/hyperledger-labs/ccapi/common"
)

func QueryV1(c *gin.Context) {
	var args []byte
	var err error

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

	channelName := os.Getenv("CHANNEL")
	chaincodeName := os.Getenv("CCNAME")
	txName := c.Param("txname")

	argList := [][]byte{}
	if args != nil {
		argList = append(argList, args)
	}

	res, status, err := chaincode.Query(channelName, chaincodeName, txName, argList)
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
