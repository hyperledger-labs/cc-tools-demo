package main

import (
	"log"
	"os"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)

	err := SetupCC()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	stub := shimtest.NewMockStub("org1MSP", new(CCDemo))
	res := stub.MockInit("testInit", [][]byte{[]byte("init")})
	if res.GetStatus() != 200 {
		log.Println(res.GetMessage())
		os.Exit(1)
	}

	os.Exit(m.Run())
}
