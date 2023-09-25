package main

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/hyperledger-labs/cc-tools/mock"
)

const clientUserOrg2Cert string = `-----BEGIN CERTIFICATE-----
MIICKzCCAdGgAwIBAgIRANlrNyW+FJdC1n3b2uqAH+gwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzIuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzIuZXhhbXBsZS5jb20wHhcNMjMwOTEyMjEwOTAwWhcNMzMwOTA5MjEwOTAw
WjBsMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MR8wHQYDVQQDDBZVc2VyMUBv
cmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfQ8ISNJB
MhJqUjjtKkkiDhmOElLMabjhw7K/4D5tKfwHilY7rOWV+XRUDKFAxw2f2ImW3AAs
cAo6shS2jPyLI6NNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYD
VR0jBCQwIoAgMFP0bQ6QE7la/v0AXYw7boHMnjisxfjhxYak1g1DhBwwCgYIKoZI
zj0EAwIDSAAwRQIhAJOeweSEKwdUZckzTv31n6Sfjsl4tXF2eyqA0tsL/voHAiBF
HYqjQ2S+f++Bjv/kFLvhdY7/acdJYsWH2xyO1XseeA==
-----END CERTIFICATE-----`

func TestGetNumberOfBooksFromLibrary(t *testing.T) {
	stub, err := mock.NewMockStubWithCert("org2MSP", new(CCDemo), []byte(clientUserOrg2Cert))
	if err != nil {
		t.FailNow()
	}

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
