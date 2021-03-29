package assettypes

import "github.com/goledgerdev/cc-tools/assets"

var Book = assets.AssetType{
	Tag:         "book",
	Label:       "Book",
	Description: "",

	Props: []assets.AssetProp{
		{
			Tag:      "title",
			Label:    "Book Title",
			Required: true,
			IsKey:    true,
			DataType: "string",
		},
		{
			Tag:      "author",
			Label:    "Book Author",
			Required: true,
			IsKey:    true,
			DataType: "string",
		},
		{
			Tag:      "currentTenant",
			Label:    "Current Tenant",
			DataType: "->person",
		},
		{
			Tag:      "genres",
			Label:    "Genres",
			DataType: "[]string",
		},
		{
			Tag:      "published",
			Label:    "Publishment Date",
			DataType: "datetime",
		},
	},
}
