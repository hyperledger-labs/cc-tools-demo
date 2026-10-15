package common

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Abort(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"status": status,
		"error":  err.Error(),
	})
	c.Error(err)
}

// Respond writes the response. On success, the status is kept only when it is a 2xx status,
// so chaincodes that always return 200 keep the previous behavior.
func Respond(c *gin.Context, res interface{}, status int, err error) {
	if err != nil {
		c.JSON(status, gin.H{
			"response": res,
			"status":   status,
			"error":    err.Error(),
		})
		c.Error(err)
		return
	}

	c.JSON(SuccessStatus(status), res)
}

// SuccessStatus returns status when it is a 2xx status, and 200 otherwise.
func SuccessStatus(status int) int {
	if status < 200 || status > 299 {
		return http.StatusOK
	}

	return status
}

// ParsePayload unmarshals a chaincode JSON payload. An empty payload, such as the
// one returned with status 204, results in a nil response instead of an error.
func ParsePayload(payload []byte) (interface{}, error) {
	if len(payload) == 0 {
		return nil, nil
	}

	var res interface{}
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
