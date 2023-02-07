package chaincode

import (
	"fmt"
	"log"
	"os"

	"github.com/goledgerdev/ccapi/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func Event(channelName, ccName, eventName string, fn func(*fab.CCEvent)) {
	// create channel manager
	fabMngr, err := common.NewFabricChClient(channelName, os.Getenv("USER"), os.Getenv("ORG"))
	if err != nil {
		log.Println("error creating channel manager: ", err)
		return
	}

	for {
		// Register chaincode event
		registration, notifier, err := fabMngr.Client.RegisterChaincodeEvent(ccName, eventName)
		if err != nil {
			log.Println("error registering chaincode event: ", err)
			return
		}

		// Execute handler function on event notification
		ccEvent := <-notifier
		fmt.Printf("Received CC event: %v\n", ccEvent)
		fn(ccEvent)

		fabMngr.Client.UnregisterChaincodeEvent(registration)
	}
}
