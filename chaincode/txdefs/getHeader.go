package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools-demo/chaincode/header"

	"github.com/goledgerdev/cc-tools/errors"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetHeader returns data in CCHeader
var GetHeader = tx.Transaction{
	Tag:         "getHeader",
	Label:       "Get Header",
	Description: "",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args:     []tx.Argument{},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		colorMap := header.Colors
		nameMap := header.Title
		orgMSP, err := cid.GetMSPID(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get MSP ID")
		}

		var colors []string
		colors, orgExists := colorMap[orgMSP]
		if !orgExists {
			colors = colorMap["@default"]
		}

		var orgTitle string
		orgTitle, orgExists = nameMap[orgMSP]
		if !orgExists {
			orgTitle = nameMap["@default"]
		}

		header := map[string]interface{}{
			"name":     header.Name,
			"version":  header.Version,
			"colors":   colors,
			"orgMSP":   orgMSP,
			"orgTitle": orgTitle,
		}
		headerBytes, err := json.Marshal(header)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal header")
		}

		return headerBytes, nil
	},
}
