package main

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/hyperledger-labs/cc-tools/mock"
)

// const clientUserOrg1Cert string = `-----BEGIN CERTIFICATE-----
// MIICKjCCAdGgAwIBAgIRAND8OJV+aiUcBgAf79V0Z0gwCgYIKoZIzj0EAwIwczEL
// MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
// cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
// Lm9yZzEuZXhhbXBsZS5jb20wHhcNMjMwOTEyMjEwOTAwWhcNMzMwOTA5MjEwOTAw
// WjBsMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
// U2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MR8wHQYDVQQDDBZVc2VyMUBv
// cmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8F737D7W
// IlTT+z19O+UbxKqa9b92j41+zCU7G4rjL/eBHyxBGf77pHXSLxH7ul9ANfmm/eFs
// PURIzoIjz3CynaNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYD
// VR0jBCQwIoAghb/LOx4HT2gRIO9/TTfiLucofrmzqF5n1iqtfjKXxVwwCgYIKoZI
// zj0EAwIDRwAwRAIgG3HIEXUXDvnj0Ce7ApBDsW7Rw8ijKrH2T9HKEEd74rkCIBCw
// Gs8rKy8blPYK6H3I64nrylOTvr5Qmkq1Ag2Tqnwb
// -----END CERTIFICATE-----`

func TestUpdateBookTenant(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(CCDemo))

	// State setup
	setupPerson := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	setupBook := map[string]interface{}{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		// "currentTenant": map[string]interface{}{
		// 	"@assetType": "person",
		// 	"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		// },
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}
	setupPersonJSON, _ := json.Marshal(setupPerson)
	setupBookJSON, _ := json.Marshal(setupBook)

	stub.MockTransactionStart("setupUpdateBookTenant")
	stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupPersonJSON)
	stub.PutState("book:a36a2920-c405-51c3-b584-dcd758338cb5", setupBookJSON)
	stub.MockTransactionEnd("setupUpdateBookTenant")

	req := map[string]interface{}{
		"book": map[string]interface{}{
			"@key": "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		},
		"tenant": map[string]interface{}{
			"@key": "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
	}
	reqBytes, _ := json.Marshal(req)

	res := stub.MockInvoke("updateBookTenant", [][]byte{
		[]byte("updateBookTenant"),
		reqBytes,
	})

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resPayload map[string]interface{}
	err := json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedResponse := map[string]interface{}{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "updateBookTenant",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	expectedResponse["@lastUpdated"] = stub.TxTimestamp.AsTime().Format(time.RFC3339)

	if !reflect.DeepEqual(resPayload, expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}

	var state map[string]interface{}
	stateBytes := stub.State["book:a36a2920-c405-51c3-b584-dcd758338cb5"]
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(state, expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", state)
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}
