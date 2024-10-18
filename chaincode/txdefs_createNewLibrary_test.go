package main_test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	cc "github.com/hyperledger-labs/cc-tools-demo/chaincode"
	"github.com/hyperledger-labs/cc-tools/mock"
)

const clientAdminOrg3Cert string = `-----BEGIN CERTIFICATE-----
MIICJjCCAcygAwIBAgIQHv152Ul3TG/REl3mHfYyUjAKBggqhkjOPQQDAjBxMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEYMBYGA1UEChMPb3JnLmV4YW1wbGUuY29tMRswGQYDVQQDExJjYS5v
cmcuZXhhbXBsZS5jb20wHhcNMjQwNTA5MjEwOTAwWhcNMzQwNTA3MjEwOTAwWjBq
MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2Fu
IEZyYW5jaXNjbzEOMAwGA1UECxMFYWRtaW4xHjAcBgNVBAMMFUFkbWluQG9yZy5l
eGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABAoAt6mlUBMB0Ab1
paR0ILegN6qKmNfOYR0WV0kGOQwkO4lYcN76lSA2wSlWNTgxtGQDzja1708Ezdr5
vJ5KFhmjTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQk
MCKAID7vB1ct0j2yeNTm45AlCyj9TW22dYtjmPOGq+SVlMKQMAoGCCqGSM49BAMC
A0gAMEUCIQDYol2ylLCcz8qrGJmAFEG/cIG2Kxv8BD5t7Gv/28y8kgIgTz0Y75p6
3kbL5VN/PCiG2SbX72AVPSiEqj6PSiZJMz4=
-----END CERTIFICATE-----`

func TestCreateNewLibrary(t *testing.T) {
	stub, err := mock.NewMockStubWithCert("org3MSP", new(cc.CCDemo), []byte(clientAdminOrg3Cert))
	if err != nil {
		t.FailNow()
	}

	expectedResponse := map[string]interface{}{
		"@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
		"@lastTouchBy": "org3MSP",
		"@lastTx":      "createNewLibrary",
		"@assetType":   "library",
		"name":         "Maria's Library",
	}
	req := map[string]interface{}{
		"name": "Maria's Library",
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createNewLibrary", [][]byte{
		[]byte("createNewLibrary"),
		reqBytes,
	})

	expectedResponse["@lastUpdated"] = stub.TxTimestamp.AsTime().Format(time.RFC3339)

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

	var state map[string]interface{}
	stateBytes := stub.State["library:3cab201f-9e2b-579d-b7b2-72297ed17f49"]
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
