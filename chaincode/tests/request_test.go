package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"

	"github.com/cucumber/godog"
)

// Key Structs definitions
type bodyCtxKey struct{}
type statusCtxKey struct{}

func iMakeARequestToWith(ctx context.Context, method, endpoint string, reqBody *godog.DocString) (context.Context, error) {
	var res *http.Response
	var err error

	// Initialize http client
	client := &http.Client{}

	dataAsBytes := bytes.NewBuffer([]byte(reqBody.Content))

	// Create request
	req, err := http.NewRequest(method, "http://localhost:1080"+endpoint, dataAsBytes)
	if err != nil {
		return ctx, err
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

func theResponseShouldMatchJson(ctx context.Context, body *godog.DocString) error {
	// Get 'ResponseBody' from context
	respBody, ok := ctx.Value(bodyCtxKey{}).([]byte)
	if !ok {
		return errors.New("unavailable context")
	}

	var expected interface{}
	var received interface{}

	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}
	if err := json.Unmarshal(respBody, &received); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, received) {
		var expectedBytes []byte
		var receivedBytes []byte
		var err error
		if expectedBytes, err = json.MarshalIndent(expected, "", "  "); err != nil {
			return err
		}
		if receivedBytes, err = json.MarshalIndent(received, "", "  "); err != nil {
			return err
		}

		return fmt.Errorf("RECEIVED json:\n%s\ndoes not match EXPECTED:\n%s", string(receivedBytes), string(expectedBytes))
	}

	return nil
}

func thereIsARunningTestNetwork(arg1 string) error {
	// Start test network
	cmd := exec.Command("../../startDev.sh")

	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I make a "([^"]*)" request to "([^"]*)" with:$`, iMakeARequestToWith)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, theResponseShouldMatchJson)
	ctx.Step(`^there is a running "([^"]*)" test network$`, thereIsARunningTestNetwork)
}
