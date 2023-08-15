package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

var BenchmarkAsset = assets.AssetType{
	Tag:         "benchmarkAsset",
	Label:       "Benchmark Asset",
	Description: "Benchmark Asset",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "ID",
			DataType: "string", // Datatypes are identified at datatypes folder
		},
		{
			Tag:      "timestamp",
			Label:    "Timestamp",
			DataType: "datetime",
		},
		{
			Tag:      "data",
			Label:    "Data",
			DataType: "string",
		},
	},
}
