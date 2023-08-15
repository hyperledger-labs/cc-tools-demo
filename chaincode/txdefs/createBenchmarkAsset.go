package txdefs

import (
	"encoding/json"
	"time"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Create a new Library on channel
// POST Method
var CreateBenchmarkAsset = tx.Transaction{
	Tag:         "createBenchmarkAsset",
	Label:       "Create Benchmark Asset",
	Description: "Create Benchmark Asset",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Tag:         "id",
			Label:       "Id",
			Description: "Id",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		id, _ := req["id"].(string)

		assetMap := make(map[string]interface{})
		assetMap["@assetType"] = "benchmarkAsset"
		assetMap["id"] = id

		txTime, _ := stub.Stub.GetTxTimestamp()
		assetMap["timestamp"] = txTime.AsTime().Format(time.RFC3339)

		asset, err := assets.NewAsset(assetMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		// Save the new library on channel
		_, err = asset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		// Marshal asset back to JSON format
		libraryJSON, nerr := json.Marshal(asset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return libraryJSON, nil
	},
}
