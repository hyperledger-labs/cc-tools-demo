package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func tryout() []error {
	var expectedResponse interface{}
	var receivedResponse interface{}
	var err error

	// Get Header
	resp, err := http.Get("http://localhost:80/api/query/getHeader")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	expectedResponse = map[string]interface{}{
		"ccToolsVersion": "v0.7.0-rc.3",
		"colors": []interface{}{
			"#4267B2",
			"#34495E",
			"#ECF0F1",
		},
		"name":     "CC Tools Demo",
		"orgMSP":   "org1MSP",
		"orgTitle": "CC Tools Demo",
		"version":  "1.0.0",
	}
	err = json.Unmarshal(body, &receivedResponse)
	if err != nil {
		log.Fatalln(err)
	}
	if !reflect.DeepEqual(expectedResponse, receivedResponse) {
		log.Println("getHeader call returned unexpected response")
		log.Printf("expected: %#v\n", expectedResponse)
		log.Fatalf("received: %#v\n", receivedResponse)
	}

	// Get Transactions
	resp, err = http.Get("http://localhost:80/api/query/getTx")
	if err != nil {
		log.Fatalln(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	expectedResponse = []interface{}{
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
	}
	err = json.Unmarshal(body, &receivedResponse)
	if err != nil {
		log.Fatalln(err)
	}
	if !reflect.DeepEqual(expectedResponse, receivedResponse) {
		log.Println("getTx call returned unexpected response")
		log.Printf("expected: %#v\n", expectedResponse)
		log.Fatalf("received: %#v\n", receivedResponse)
	}

	return nil
}
