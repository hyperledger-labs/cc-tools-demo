package main

import (
	"github.com/goledgerdev/cc-tools-demo/chaincode/eventtypes"
	"github.com/goledgerdev/cc-tools/events"
)

var eventTypeList = []events.Event{
	eventtypes.TestEvent,
}
