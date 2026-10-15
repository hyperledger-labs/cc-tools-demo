package main

import (
	"github.com/hyperledger-labs/cc-tools-demo/chaincode/eventtypes"
	"github.com/hyperledger-labs/cc-tools/v2/events"
)

var eventTypeList = []events.Event{
	eventtypes.CreateLibraryLog,
}
