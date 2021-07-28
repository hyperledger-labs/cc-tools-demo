package main

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestGetNumberOfBooksFromLibrary(t *testing.T) {
	stub := mock.NewMockStub("org2MSP", new(CCDemo))

	// Setup state
	setupBook := map[string]interface{}{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"genres":       []interface{}{"biography", "non-fiction"},
		"published":    "2019-05-06T22:12:41Z",
	}
	setupLibrary := map[string]interface{}{
		"@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
		"@lastTouchBy": "org3MSP",
		"@lastTx":      "createNewLibrary",
		"@assetType":   "library",
		"name":         "Maria's Library",
		"books": []map[string]interface{}{
			{
				"@assetType": "book",
				"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
			},
		},
	}
	setupBookJSON, _ := json.Marshal(setupBook)
	setupLibraryJSON, _ := json.Marshal(setupLibrary)

	stub.MockTransactionStart("setupGetNumberOfBooksFromLibrary")
	stub.PutState("book:a36a2920-c405-51c3-b584-dcd758338cb5", setupBookJSON)
	stub.PutState("library:3cab201f-9e2b-579d-b7b2-72297ed17f49", setupLibraryJSON)
	refIdx, err := stub.CreateCompositeKey("book:a36a2920-c405-51c3-b584-dcd758338cb5", []string{"library:3cab201f-9e2b-579d-b7b2-72297ed17f49"})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.PutState(refIdx, []byte{0x00})
	stub.MockTransactionEnd("setupGetNumberOfBooksFromLibrary")

	expectedResponse := map[string]interface{}{
		"numberOfBooks": 1.0,
	}
	req := map[string]interface{}{
		"library": map[string]interface{}{
			"name": "Maria's Library",
		},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("getNumberOfBooksFromLibrary", [][]byte{
		[]byte("getNumberOfBooksFromLibrary"),
		reqBytes,
	})

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resPayload map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(resPayload, expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}
