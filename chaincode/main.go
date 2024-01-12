package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger-labs/cc-tools-demo/chaincode/assettypes"
	"github.com/hyperledger-labs/cc-tools-demo/chaincode/datatypes"
	"github.com/hyperledger-labs/cc-tools-demo/chaincode/header"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

var startupCheckExecuted = false

func SetupCC() error {
	tx.InitHeader(tx.Header{
		Name:    header.Name,
		Version: header.Version,
		Colors:  header.Colors,
		Title:   header.Title,
	})

	assets.InitDynamicAssetTypeConfig(assettypes.DynamicAssetTypes)

	tx.InitTxList(txList)

	err := assets.CustomDataTypes(datatypes.CustomDataTypes)
	if err != nil {
		fmt.Printf("Error injecting custom data types: %s", err)
		return err
	}
	assets.InitAssetList(append(assetTypeList, assettypes.CustomAssets...))

	events.InitEventList(eventTypeList)

	return nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	// Generate collection json
	genFlag := flag.Bool("g", false, "Enable collection generation")
	flag.Bool("orgs", false, "List of orgs to generate collection for")
	flag.Parse()
	if *genFlag {
		listOrgs := flag.Args()
		generateCollection(listOrgs)
		return
	}

	log.Printf("Starting chaincode %s version %s\n", header.Name, header.Version)

	err := SetupCC()
	if err != nil {
		return
	}
	if err = shim.Start(new(CCDemo)); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// CCDemo implements the shim.Chaincode interface
type CCDemo struct{}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *CCDemo) Init(stub shim.ChaincodeStubInterface) (response pb.Response) {

	res := InitFunc(stub)
	startupCheckExecuted = true
	if res.Status != 200 {
		return res
	}

	response = shim.Success(nil)
	return
}

func InitFunc(stub shim.ChaincodeStubInterface) (response pb.Response) {
	// Defer logging function
	defer logTx(stub, time.Now(), &response)

	if assettypes.DynamicAssetTypes.Enabled {
		sw := &sw.StubWrapper{
			Stub: stub,
		}
		err := assets.RestoreAssetList(sw, true)
		if err != nil {
			response = err.GetErrorResponse()
			return
		}
	}

	err := assets.StartupCheck()
	if err != nil {
		response = err.GetErrorResponse()
		return
	}

	err = tx.StartupCheck()
	if err != nil {
		response = err.GetErrorResponse()
		return
	}

	response = shim.Success(nil)
	return
}

// Invoke is called per transaction on the chaincode.
func (t *CCDemo) Invoke(stub shim.ChaincodeStubInterface) (response pb.Response) {
	// Defer logging function
	defer logTx(stub, time.Now(), &response)

	if !startupCheckExecuted {
		fmt.Println("Running startup check...")
		res := InitFunc(stub)
		if res.Status != 200 {
			return res
		}
		startupCheckExecuted = true
	}

	var result []byte

	result, err := tx.Run(stub)

	if err != nil {
		response = err.GetErrorResponse()
		return
	}
	response = shim.Success([]byte(result))
	return
}

func logTx(stub shim.ChaincodeStubInterface, beginTime time.Time, response *pb.Response) {
	fn, _ := stub.GetFunctionAndParameters()
	log.Printf("%d %s %s %s\n", response.Status, fn, time.Since(beginTime), response.Message)
}
