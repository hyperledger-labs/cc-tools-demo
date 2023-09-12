package main_test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	cc "github.com/goledgerdev/cc-tools-demo/chaincode"
	"github.com/goledgerdev/cc-tools/mock"
)

const clientUserOrg3Cert string = `-----BEGIN CERTIFICATE-----
MIICKzCCAdGgAwIBAgIRANRvqM5kq+BSl9lWvvdKxZowCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzMuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzMuZXhhbXBsZS5jb20wHhcNMjMwOTEyMjEwOTAwWhcNMzMwOTA5MjEwOTAw
WjBsMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MR8wHQYDVQQDDBZVc2VyMUBv
cmczLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEj37aqXOj
dND0KVyKTCt6knrVpFx9N/WX5y8jVNKStIBZ2YSFxbf6CdwVszFHA054wA3tr9TA
0oOlgB2z7OML5qNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYD
VR0jBCQwIoAg8VyvOpzJ9YtVRTx93VWDKbu+L4OpFOk4B91QvlfUlNIwCgYIKoZI
zj0EAwIDSAAwRQIhAN8qQuiX0lwdYDxtISLF77EEw2BQ5qnVwyuL49tNufAsAiAp
9sXexeuPIEjl7MU9pRSgWLqBPbyLSvtPq3ZQxIL/gg==
-----END CERTIFICATE-----`

func TestCreateNewLibrary(t *testing.T) {
	stub, err := mock.NewMockStubWithCert("org3MSP", new(cc.CCDemo), []byte(clientUserOrg3Cert))
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
