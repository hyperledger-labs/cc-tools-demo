package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/ccapi/chaincode"
	"github.com/goledgerdev/ccapi/common"
	"github.com/pkg/errors"
)

func InvokeGatewayDefault(c *gin.Context) {
	channelName := os.Getenv("CHANNEL")
	chaincodeName := os.Getenv("CCNAME")

	invokeGateway(c, channelName, chaincodeName)
}

func InvokeGatewayCustom(c *gin.Context) {
	channelName := c.Param("channelName")
	chaincodeName := c.Param("chaincodeName")

	invokeGateway(c, channelName, chaincodeName)
}

func invokeGateway(c *gin.Context, channelName, chaincodeName string) {
	// Get channel information from request
	req := make(map[string]interface{})
	err := c.BindJSON(&req)
	if err != nil {
		common.Abort(c, http.StatusBadRequest, err)
		return
	}

	txName := c.Param("txname")

	var collections []string
	collectionsQuery := c.Query("@collections")
	if collectionsQuery != "" {
		collectionsByte, err := base64.StdEncoding.DecodeString(collectionsQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the @collections query parameter must be a base64-encoded JSON array of strings",
			})
			return
		}

		err = json.Unmarshal(collectionsByte, &collections)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the @collections query parameter must be a base64-encoded JSON array of strings",
			})
			return
		}
	} else {
		collectionsQuery := c.QueryArray("collections")
		if len(collectionsQuery) > 0 {
			collections = collectionsQuery
		} else {
			collections = []string{c.Query("collections")}
		}
	}

	// Make args
	reqBytes, err := json.Marshal(req)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, errors.Wrap(err, "failed to marshal req body"))
		return
	}

	// Invoke
	result, err := chaincode.InvokeGateway(channelName, chaincodeName, txName, string(reqBytes), nil, nil)
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
