package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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
	req, err := http.NewRequest(method, "http://localhost:8080"+endpoint, dataAsBytes)
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

func theResponseShouldMatchJson(arg1 *godog.DocString) error {
	return godog.ErrPending
}

func thereIsARunningTestNetwork(arg1 string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I make a "([^"]*)" request to "([^"]*)" with:$`, iMakeARequestToWith)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, theResponseShouldMatchJson)
	ctx.Step(`^there is a running "([^"]*)" test network$`, thereIsARunningTestNetwork)
}
