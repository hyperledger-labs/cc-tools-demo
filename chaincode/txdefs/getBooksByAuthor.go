package txdefs

import (
	"encoding/json"
	"fmt"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Return the all books from an specific author
// GET method
var GetBooksByAuthor = tx.Transaction{
	Tag:         "getBooksByAuthor",
	Label:       "Get Books by the Author Name",
	Description: "Return all the books from an author",
	Method:      "GET",
	Callers:     []string{"$org1MSP", "$org2MSP"}, // Only org1 and org2 can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "authorName",
			Label:       "Author Name",
			Description: "Author Name",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "limit",
			Label:       "Limit",
			Description: "Limit",
			DataType:    "number",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		authorName, _ := req["authorName"].(string)
		limit, hasLimit := req["limit"].(float64)

		if hasLimit && limit <= 0 {
			return nil, errors.NewCCError("limit must be greater than 0", 400)
		}

		//prepare couchdb query
		queryString := fmt.Sprintf(`{
			"selector": {
				"@assetType": "book",
				"author": "%s"
			}
		}`, authorName)

		query := make(map[string]interface{})
		err := json.Unmarshal(([]byte)(queryString), &query)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error unmarshaling query string", 500)
		}

		if hasLimit {
			query["limit"] = limit
		}

		response, err := assets.Search(stub, query, "", true)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error searching for book's author", 500)
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", 500)
		}

		return responseJSON, nil
	},
}
