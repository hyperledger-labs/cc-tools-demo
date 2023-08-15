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
var UpdateBenchmarkAsset = tx.Transaction{
	Tag:         "updateBenchmarkAsset",
	Label:       "Update Benchmark Asset",
	Description: "Update Benchmark Asset",
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

		asset, err := assets.NewAsset(assetMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		assetValue, _ := asset.Get(stub)

		txTime, _ := stub.Stub.GetTxTimestamp()
		assetMap["timestamp"] = txTime.AsTime().Format(time.RFC3339)

		res, _ := assetValue.Update(stub, assetMap)

		// Marshal asset back to JSON format
		JSON, nerr := json.Marshal(res)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return JSON, nil
	},
}
