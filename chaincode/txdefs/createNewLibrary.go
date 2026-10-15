package txdefs

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger-labs/cc-tools/v2/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/v2/assets"
	"github.com/hyperledger-labs/cc-tools/v2/errors"
	"github.com/hyperledger-labs/cc-tools/v2/events"
	sw "github.com/hyperledger-labs/cc-tools/v2/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/v2/transactions"
)

// Create a new Library on channel
// POST Method
// Uses ResponseRoutine to respond with status 201 (Created) instead of the default 200
var CreateNewLibrary = tx.Transaction{
	Tag:         "createNewLibrary",
	Label:       "Create New Library",
	Description: "Create a New Library",
	Method:      "POST",
	Callers: []accesscontrol.Caller{ // Only org3 admin can call this transaction
		{
			MSP: "org3MSP",
			OU:  "admin",
		},
		{
			MSP: "orgMSP",
			OU:  "admin",
		},
	},

	Args: []tx.Argument{
		{
			Tag:         "name",
			Label:       "Name",
			Description: "Name of the library",
			DataType:    "string",
			Required:    true,
		},
	},
	ResponseRoutine: func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
		name, _ := req["name"].(string)

		libraryMap := make(map[string]interface{})
		libraryMap["@assetType"] = "library"
		libraryMap["name"] = name

		libraryAsset, err := assets.NewAsset(libraryMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		// Save the new library on channel
		_, err = libraryAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "Error saving asset on blockchain", err.Status())
		}

		// Marshal asset back to JSON format
		libraryJSON, nerr := json.Marshal(libraryAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		// Marshall message to be logged
		logMsg, ok := json.Marshal(fmt.Sprintf("New library name: %s", name))
		if ok != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		// Call event to log the message
		events.CallEvent(stub, "createLibraryLog", logMsg)

		return &tx.Response{Status: 201, Payload: libraryJSON}, nil
	},
}
