package assettypes

import (
	"fmt"

	"github.com/goledgerdev/cc-tools/assets"
)

var Person = assets.AssetType{
	Tag:         "person",
	Label:       "Person",
	Description: "",

	Props: []assets.AssetProp{
		{
			Tag:      "cpf",
			Label:    "CPF",
			DataType: "cpf",
		},
		{
			Tag:      "name",
			Label:    "Asset Name",
			Required: true,
			IsKey:    true,
			DataType: "string",
			Validate: func(name interface{}) error {
				nameStr := name.(string)
				if nameStr == "" {
					return fmt.Errorf("name must be non-empty")
				}
				return nil
			},
		},
		{
			Tag:          "readerScore",
			Label:        "Reader Score",
			DefaultValue: 0,
			DataType:     "number",
			Writers:      []string{`org1MSP`},
		},
	},
}
