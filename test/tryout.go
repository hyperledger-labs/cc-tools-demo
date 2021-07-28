package main

import (
	"fmt"
	"log"
)

func tryout() []error {
	var err error

	// Get Header
	fmt.Print("Get Header... ")
	err = GetAndVerify(
		"http://localhost:80/api/query/getHeader",
		200,
		map[string]interface{}{
			"ccToolsVersion": "v0.7.1",
			"colors": []interface{}{
				"#4267B2",
				"#34495E",
				"#ECF0F1",
			},
			"name":     "CC Tools Demo",
			"orgMSP":   "org1MSP",
			"orgTitle": "CC Tools Demo",
			"version":  "1.0.0",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Get Transactions
	fmt.Print("Get Transactions... ")
	err = GetAndVerify(
		"http://localhost:80/api/query/getTx",
		200,
		[]interface{}{
			map[string]interface{}{
				"description": "",
				"label":       "Create Asset",
				"tag":         "createAsset",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Update Asset",
				"tag":         "updateAsset",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Delete Asset",
				"tag":         "deleteAsset",
			},
			map[string]interface{}{
				"callers":     []interface{}{"$org3MSP"},
				"description": "Create a New Library",
				"label":       "Create New Library",
				"tag":         "createNewLibrary",
			},
			map[string]interface{}{
				"callers":     []interface{}{"$org2MSP"},
				"description": "Return the number of books of a library",
				"label":       "Get Number Of Books From Library",
				"tag":         "getNumberOfBooksFromLibrary",
			},
			map[string]interface{}{
				"callers":     []interface{}{"$org\\dMSP"},
				"description": "Change the tenant of a book",
				"label":       "Update Book Tenant",
				"tag":         "updateBookTenant",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Get Tx",
				"tag":         "getTx",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Get Header",
				"tag":         "getHeader",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Get Schema",
				"tag":         "getSchema",
			},
			map[string]interface{}{
				"description": "GetDataTypes returns the primary data type map",
				"label":       "Get DataTypes",
				"tag":         "getDataTypes",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Read Asset",
				"tag":         "readAsset",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Read Asset History",
				"tag":         "readAssetHistory",
			},
			map[string]interface{}{
				"description": "",
				"label":       "Search World State",
				"tag":         "search",
			},
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Get CreateAsset definition
	fmt.Print("Get CreateAsset definition... ")
	err = PostAndVerify(
		"http://localhost:80/api/query/getTx",
		map[string]interface{}{
			"txName": "createAsset",
		},
		200,
		map[string]interface{}{
			"args": []interface{}{
				map[string]interface{}{
					"dataType":    "[]@asset",
					"description": "List of assets to be created.",
					"label":       "",
					"private":     false,
					"required":    true,
					"tag":         "asset",
				},
			},
			"description": "",
			"label":       "Create Asset",
			"metaTx":      true,
			"method":      "POST",
			"readOnly":    false,
			"tag":         "createAsset",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Get Asset Types
	fmt.Print("Get Asset Types... ")
	err = GetAndVerify(
		"http://localhost:80/api/query/getSchema",
		200,
		[]interface{}{
			map[string]interface{}{
				"description": "Personal data of someone",
				"label":       "Person",
				"tag":         "person",
				"writers":     nil,
			},
			map[string]interface{}{
				"description": "Book",
				"label":       "Book",
				"tag":         "book",
				"writers":     nil,
			},
			map[string]interface{}{
				"description": "Library as a collection of books",
				"label":       "Library",
				"tag":         "library",
				"writers":     nil,
			},
			map[string]interface{}{
				"description": "Secret between Org2 and Org3",
				"label":       "Secret",
				"readers": []interface{}{
					"org2MSP",
					"org3MSP",
				},
				"tag":     "secret",
				"writers": nil,
			},
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Get Person asset type definition
	fmt.Print("Get Person asset type definition... ")
	err = PostAndVerify(
		"http://localhost:80/api/query/getSchema",
		map[string]interface{}{
			"assetType": "person",
		},
		200,
		map[string]interface{}{
			"tag":         "person",
			"label":       "Person",
			"description": "Personal data of someone",
			"props": []interface{}{
				map[string]interface{}{
					"dataType":    "cpf",
					"description": "",
					"isKey":       true,
					"label":       "CPF (Brazilian ID)",
					"readOnly":    false,
					"required":    true,
					"tag":         "id",
					"writers": []interface{}{
						"org1MSP",
					},
				},
				map[string]interface{}{
					"dataType":    "string",
					"description": "",
					"isKey":       false,
					"label":       "Name of the person",
					"readOnly":    false,
					"required":    true,
					"tag":         "name",
					"writers":     nil,
				},
				map[string]interface{}{
					"dataType":    "datetime",
					"description": "",
					"isKey":       false,
					"label":       "Date of Birth",
					"readOnly":    false,
					"required":    false,
					"tag":         "dateOfBirth",
					"writers":     nil,
				},
				map[string]interface{}{
					"dataType":     "number",
					"defaultValue": 0.0,
					"description":  "",
					"isKey":        false,
					"label":        "Person's height",
					"readOnly":     false,
					"required":     false,
					"tag":          "height",
					"writers":      nil,
				},
			},
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Create Person
	fmt.Print("Create Person... ")
	err = PostAndVerify(
		"http://localhost:80/api/invoke/createAsset",
		map[string]interface{}{
			"asset": []map[string]interface{}{
				{
					"@assetType": "person",
					"name":       "Maria",
					"id":         "318.207.920-48",
				},
			},
		},
		200,
		[]interface{}{
			map[string]interface{}{
				"@assetType":   "person",
				"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
				"@lastTouchBy": "org1MSP",
				"@lastTx":      "createAsset",
				"height":       0.0,
				"id":           "31820792048",
				"name":         "Maria",
			},
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Create Book
	fmt.Print("Create Book... ")
	err = PostAndVerify(
		"http://localhost:980/api/invoke/createAsset",
		map[string]interface{}{
			"asset": []map[string]interface{}{
				{
					"@assetType": "book",
					"title":      "Meu Nome é Maria",
					"author":     "Maria Viana",
					"currentTenant": map[string]interface{}{
						"id": "318.207.920-48",
					},
					"genres": []string{
						"biography",
						"non-fiction",
					},
					"published": "2019-05-06T22:12:41Z",
				},
			},
		},
		200,
		[]interface{}{
			map[string]interface{}{
				"@assetType":   "book",
				"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
				"@lastTouchBy": "org2MSP",
				"@lastTx":      "createAsset",
				"title":        "Meu Nome é Maria",
				"author":       "Maria Viana",
				"currentTenant": map[string]interface{}{
					"@assetType": "person",
					"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
				},
				"genres": []interface{}{
					"biography",
					"non-fiction",
				},
				"published": "2019-05-06T22:12:41Z",
			},
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Read Book
	fmt.Print("Read Book... ")
	err = PostAndVerify(
		"http://localhost:80/api/query/readAsset",
		map[string]interface{}{
			"key": map[string]interface{}{
				"@assetType": "book",
				"author":     "Maria Viana",
				"title":      "Meu Nome é Maria",
			},
			"resolve": true,
		},
		200,
		map[string]interface{}{
			"@assetType":   "book",
			"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
			"@lastTouchBy": "org2MSP",
			"@lastTx":      "createAsset",
			"title":        "Meu Nome é Maria",
			"author":       "Maria Viana",
			"currentTenant": map[string]interface{}{
				"@assetType":   "person",
				"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
				"@lastTouchBy": "org1MSP",
				"@lastTx":      "createAsset",
				"height":       0.0,
				"id":           "31820792048",
				"name":         "Maria",
			},
			"genres": []interface{}{
				"biography",
				"non-fiction",
			},
			"published": "2019-05-06T22:12:41Z",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Update Person
	fmt.Print("Update Person... ")
	err = PostAndVerify(
		"http://localhost:80/api/invoke/updateAsset",
		map[string]interface{}{
			"update": map[string]interface{}{
				"@assetType": "person",
				"id":         "318.207.920-48",
				"name":       "Maria",
				"height":     1.66,
			},
		},
		200,
		map[string]interface{}{
			"@assetType":   "person",
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "updateAsset",
			"height":       1.66,
			"id":           "31820792048",
			"name":         "Maria",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Check if person was updated
	fmt.Print("Check if person was updated... ")
	err = PostAndVerify(
		"http://localhost:80/api/query/readAsset",
		map[string]interface{}{
			"key": map[string]interface{}{
				"@assetType": "person",
				"id":         "318.207.920-48",
			},
		},
		200,
		map[string]interface{}{
			"@assetType":   "person",
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "updateAsset",
			"height":       1.66,
			"id":           "31820792048",
			"name":         "Maria",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Query all books using CouchDB
	fmt.Print("Query all books using CouchDB... ")
	err = PostAndVerify(
		"http://localhost:80/api/query/search",
		map[string]interface{}{
			"query": map[string]interface{}{
				"selector": map[string]interface{}{
					"@assetType": "book",
				},
			},
			"resolve": true,
		},
		200,
		map[string]interface{}{
			"metadata": map[string]interface{}{},
			"result": []interface{}{
				map[string]interface{}{
					"@assetType":   "book",
					"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
					"@lastTouchBy": "org2MSP",
					"@lastTx":      "createAsset",
					"author":       "Maria Viana",
					"currentTenant": map[string]interface{}{
						"@assetType":   "person",
						"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
						"@lastTouchBy": "org1MSP",
						"@lastTx":      "updateAsset",
						"height":       1.66,
						"id":           "31820792048",
						"name":         "Maria"},
					"genres": []interface{}{
						"biography",
						"non-fiction",
					},
					"published": "2019-05-06T22:12:41Z",
					"title":     "Meu Nome é Maria",
				},
			},
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Delete book
	fmt.Print("Delete book... ")
	err = PostAndVerify(
		"http://localhost:980/api/invoke/deleteAsset",
		map[string]interface{}{
			"key": map[string]interface{}{
				"@assetType": "book",
				"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
			},
		},
		200,
		map[string]interface{}{
			"@assetType":   "book",
			"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
			"@lastTouchBy": "org2MSP",
			"@lastTx":      "createAsset",
			"title":        "Meu Nome é Maria",
			"author":       "Maria Viana",
			"currentTenant": map[string]interface{}{
				"@assetType": "person",
				"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			},
			"genres": []interface{}{
				"biography",
				"non-fiction",
			},
			"published": "2019-05-06T22:12:41Z",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	// Delete person
	fmt.Print("Delete person... ")
	err = PostAndVerify(
		"http://localhost:80/api/invoke/deleteAsset",
		map[string]interface{}{
			"key": map[string]interface{}{
				"@assetType": "person",
				"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			},
		},
		200,
		map[string]interface{}{
			"@assetType":   "person",
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "updateAsset",
			"height":       1.66,
			"id":           "31820792048",
			"name":         "Maria",
		},
	)
	if err != nil {
		fail()
		log.Fatalln(err)
	}
	pass()

	return nil
}
