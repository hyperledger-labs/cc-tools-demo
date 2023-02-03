package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

// DynamicAssetTypes contains the configuration for the Dynamic AssetTypes feature.
var DynamicAssetTypes = assets.DynamicAssetType{
	Enabled:     true,
	AssetAdmins: []string{`org1MSP`, "orgMSP"},
}
