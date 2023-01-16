package tests

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

// Key Structs definitions
type bodyCtxKey struct{}
type statusCtxKey struct{}

func iMakeARequestToOnPortWith(ctx context.Context, method, endpoint string, port int, reqBody *godog.DocString) (context.Context, error) {
	var res *http.Response
	var req *http.Request
	var err error

	// Initialize http client
	client := &http.Client{}

	// Create request
	if method == "GET" {
		b64str := b64.StdEncoding.EncodeToString([]byte(reqBody.Content))
		reqParam := "?@request=" + b64str
		req, err = http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+endpoint+reqParam, nil)
		if err != nil {
			return ctx, err
		}
	} else {
		dataAsBytes := bytes.NewBuffer([]byte(reqBody.Content))

		req, err = http.NewRequest(method, "http://localhost:"+strconv.Itoa(port)+endpoint, dataAsBytes)
		if err != nil {
			return ctx, err
		}
	}

	// Set header and make request
	req.Header.Set("Content-Type", "application/json")
	res, err = client.Do(req)
	if err != nil {
		return ctx, err
	}

	// Get status code and response body
	statusCode := res.StatusCode
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return ctx, err
	}
	res.Body.Close()

	// Append status code and response body to context
	bodyCtx := context.WithValue(ctx, bodyCtxKey{}, resBody)
	statusCtx := context.WithValue(bodyCtx, statusCtxKey{}, statusCode)
	return statusCtx, nil
}

func theResponseCodeShouldBe(ctx context.Context, expectedCode int) (context.Context, error) {
	// Get Status Code from context
	statusCode, ok := ctx.Value(statusCtxKey{}).(int)
	if !ok {
		return ctx, errors.New("context unaveilable while retrieving status")
	}

	if statusCode != expectedCode {
		// Get Response body from context
		resBody, ok := ctx.Value(bodyCtxKey{}).([]byte)
		if !ok {
			return ctx, errors.New("unaveilable context while retrieving body")
		}

		// Test Failed
		return ctx, fmt.Errorf("received wrong status response. Got %d Expected: %d\nResponse body: %s", statusCode, expectedCode, string(resBody))
	}

	return ctx, nil
}

func theResponseShouldHave(ctx context.Context, body *godog.DocString) error {
	// Get 'ResponseBody' from context
	respBody, ok := ctx.Value(bodyCtxKey{}).([]byte)
	if !ok {
		return errors.New("unavailable context")
	}

	var expected map[string]interface{}
	var received map[string]interface{}

	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}
	if err := json.Unmarshal(respBody, &received); err != nil {
		return err
	}

	for key, value := range expected {
		if !reflect.DeepEqual(value, received[key]) {
			var expectedBytes []byte
			var receivedBytes []byte
			var err error
			if expectedBytes, err = json.MarshalIndent(value, "", "  "); err != nil {
				return err
			}
			if receivedBytes, err = json.MarshalIndent(received[key], "", "  "); err != nil {
				return err
			}

			return fmt.Errorf("Expected %s to be equal %s, but received %s", key, expectedBytes, receivedBytes)
		}
	}

	return nil
}

func theFieldShouldHaveSize(ctx context.Context, field string, expectedSize int) (context.Context, error) {
	// Get 'ResponseBody' from context
	respBody, ok := ctx.Value(bodyCtxKey{}).([]byte)
	if !ok {
		return ctx, errors.New("unavailable context")
	}

	var bodyMap map[string]interface{}

	if err := json.Unmarshal(respBody, &bodyMap); err != nil {
		return ctx, err
	}

	resultField, ok := bodyMap[field]
	if !ok {
		return ctx, errors.New("unavailable filed on body response")
	}

	fieldLen := len(resultField.([]interface{}))
	if fieldLen != expectedSize {
		// Test Failed
		return ctx, fmt.Errorf("received wrong filed size on response body. Got %d Expected: %d\nResponse body: %s", fieldLen, expectedSize, string(respBody))
	}

	return ctx, nil
}

func thereAreBooksWithPrefixByAuthor(ctx context.Context, nBooks int, prefix string, author string) (context.Context, error) {
	var res *http.Response
	var err error

	for i := 1; i <= nBooks; i++ {
		// Verify if book already exists
		requestJSON := map[string]interface{}{
			"query": map[string]interface{}{
				"selector": map[string]interface{}{
					"author":     author,
					"title":      prefix + strconv.Itoa(i),
					"@assetType": "book",
				},
			},
			"resolve": true,
		}
		jsonStr, e := json.Marshal(requestJSON)
		if e != nil {
			return ctx, err
		}
		dataAsBytes := bytes.NewBuffer([]byte(jsonStr))

		if res, err = http.Post("http://localhost/api/query/search", "application/json", dataAsBytes); err != nil {
			return ctx, err
		}
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return ctx, err
		}
		res.Body.Close()

		var received map[string]interface{}
		if err = json.Unmarshal(resBody, &received); err != nil {
			return ctx, err
		}

		// Create book if it doesnt exists
		if len(received["result"].([]interface{})) == 0 {
			requestJSON := map[string]interface{}{
				"asset": []interface{}{
					map[string]interface{}{
						"author":     author,
						"title":      prefix + strconv.Itoa(i),
						"@assetType": "book",
					},
				},
			}
			jsonStr, e := json.Marshal(requestJSON)
			if e != nil {
				return ctx, err
			}
			dataAsBytes := bytes.NewBuffer([]byte(jsonStr))

			if res, err = http.Post("http://localhost:880/api/invoke/createAsset", "application/json", dataAsBytes); err != nil {
				return ctx, err
			}

			if res.StatusCode != 200 {
				return ctx, errors.New("Failed to create book asset")
			}
		}
	}

	return ctx, nil
}

func thereIsALibraryWithName(ctx context.Context, name string) (context.Context, error) {
	var res *http.Response
	var err error

	// Verify if library already exists
	requestJSON := map[string]interface{}{
		"query": map[string]interface{}{
			"selector": map[string]interface{}{
				"name":       name,
				"@assetType": "library",
			},
		},
		"resolve": true,
	}
	jsonStr, e := json.Marshal(requestJSON)
	if e != nil {
		return ctx, err
	}
	dataAsBytes := bytes.NewBuffer([]byte(jsonStr))

	if res, err = http.Post("http://localhost/api/query/search", "application/json", dataAsBytes); err != nil {
		return ctx, err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ctx, err
	}
	res.Body.Close()

	var received map[string]interface{}
	if err = json.Unmarshal(resBody, &received); err != nil {
		return ctx, err
	}

	// Create library if it doesnt exists
	if len(received["result"].([]interface{})) == 0 {
		requestJSON = map[string]interface{}{
			"asset": []interface{}{
				map[string]interface{}{
					"name":       name,
					"@assetType": "library",
				},
			},
		}
		jsonStr, e = json.Marshal(requestJSON)
		if e != nil {
			return ctx, err
		}
		dataAsBytes = bytes.NewBuffer([]byte(jsonStr))

		if res, err = http.Post("http://localhost:880/api/invoke/createAsset", "application/json", dataAsBytes); err != nil {
			return ctx, err
		}

		if res.StatusCode != 200 {
			return ctx, errors.New("Failed to create library asset")
		}
	}

	return ctx, nil
}

func thereIsARunningTestNetworkFromScratch(arg1 string) error {
	// Start test network with 1 org only
	cmd := exec.Command("../../startDev.sh", "-n", "1")

	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Wait for ccapi
	err = waitForNetwork("880")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func thereIsARunningTestNetwork(arg1 string) error {
	if !verifyContainer("api.org.example.com", "3000") {
		// Start test network with 1 org only
		cmd := exec.Command("../../startDev.sh", "-n", "1")

		_, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// Wait for ccapi
		err = waitForNetwork("880")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I make a "([^"]*)" request to "([^"]*)" on port (\d+) with:$`, iMakeARequestToOnPortWith)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should have:$`, theResponseShouldHave)
	ctx.Step(`^there is a running "([^"]*)" test network$`, thereIsARunningTestNetwork)
	ctx.Step(`^there is a running "([^"]*)" test network from scratch$`, thereIsARunningTestNetworkFromScratch)
	ctx.Step(`^there are (\d+) books with prefix "([^"]*)" by author "([^"]*)"$`, thereAreBooksWithPrefixByAuthor)
	ctx.Step(`^the "([^"]*)" field should have size (\d+)$`, theFieldShouldHaveSize)
	ctx.Step(`^there is a library with name "([^"]*)"$`, thereIsALibraryWithName)
}

func waitForNetwork(port string) error {
	channel := make(chan error, 1)
	t := time.NewTimer(3 * time.Minute)

	defer t.Stop()

	go func() {
		for {
			_, err := http.Post("http://localhost:"+port+"/api", "application/json", nil)

			if err == nil {
				break
			}

			time.Sleep(1 * time.Second)
		}

		channel <- nil
	}()

	select {
	case err := <-channel:
		return err
	case <-t.C:
		return errors.New("Timed out waiting for network")
	}
}

func verifyContainer(container string, port string) bool {
	strCmd := "docker inspect --format=\"{{json .Config.ExposedPorts }}\" " + container

	cmd := exec.Command("bash", "-c", strCmd)
	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		return false
	}

	// If ccapi is the expected container, continue
	return outb.String() == "{\""+port+"/tcp\":{}}\n"
}
