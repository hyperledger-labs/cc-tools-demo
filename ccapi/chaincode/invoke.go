package chaincode

import (
	"net/http"
	"os"

	"github.com/goledgerdev/ccapi/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

func Invoke(channelName, ccName, txName string, txArgs [][]byte) (*channel.Response, int, error) {
	// create channel manager
	fabMngr, err := common.NewFabricChClient(channelName, os.Getenv("USER"), os.Getenv("ORG"))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// TODO: Check support for private collections
	// TODO: Check support for transient map

	// Execute chaincode with channel's client
	rq := channel.Request{ChaincodeID: ccName, Fcn: txName, Args: txArgs}
	res, err := fabMngr.Client.Execute(rq, channel.WithRetry(retry.DefaultChannelOpts))

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &res, http.StatusInternalServerError, nil
}
