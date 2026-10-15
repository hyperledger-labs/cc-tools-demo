package main

import (
	"github.com/hyperledger-labs/cc-tools-demo/chaincode/assettypes"
	"github.com/hyperledger-labs/cc-tools/v2/assets"
)

var assetTypeList = []assets.AssetType{
	assettypes.Person,
	assettypes.Book,
	assettypes.Library,
	assettypes.Secret,
}
