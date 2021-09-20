package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Updates the address of a Library
// PUT Method
var UpdateLibraryAddress = tx.Transaction{
	Tag:         "updateLibraryAddress",
	Label:       "Update Library Name",
	Description: "Change the address of a library",
	Method:      "PUT",
	Callers:     []string{`$org3MSP`}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "library",
			Label:       "Library",
			Description: "Library",
			DataType:    "->library",
			Required:    true,
		},
		{
			Tag:         "address",
			Label:       "Library address",
			Description: "New library address",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		libraryKey, _ := req["library"].(assets.Key)
		address, _ := req["address"].(string)

		// Returns Libraryfrom channel
		var libraryAsset *assets.Asset
		libraryAsset, err := libraryKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}
		libraryMap := (map[string]interface{})(*libraryAsset)

		// Update data
		libraryMap["address"] = address

		libraryMap, err = libraryAsset.Update(stub, libraryMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}

		// Marshal asset back to JSON format
		libraryJSON, nerr := json.Marshal(libraryMap)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return libraryJSON, nil
	},
}
