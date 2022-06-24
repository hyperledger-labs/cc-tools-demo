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
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

// Key Structs definitions
type bodyCtxKey struct{}
type statusCtxKey struct{}

func iMakeARequestToOnPortWith(ctx context.Context, method, endpoint string, port int, reqBody *godog.DocString) (context.Context, error) {
	var res *http.Response
	var err error

	// Initialize http client
	client := &http.Client{}

	dataAsBytes := bytes.NewBuffer([]byte(reqBody.Content))

	// Create request
	req, err := http.NewRequest(method, "http://localhost:"+strconv.Itoa(port)+endpoint, dataAsBytes)
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

func thereIsARunningTestNetwork(arg1 string) error {
	// Start test network
	cmd := exec.Command("../../startDev.sh")

	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Wait for ccapi of all orgs
	for i := 1; i <= 3; i++ {
		err = waitForNetwork("org" + strconv.Itoa(i))
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
}

func waitForNetwork(org string) error {
	// Read last line of ccapi log
	strCmd := "docker logs ccapi." + org + ".example.com | tail -n 1"

	wait := true

	for wait {
		// Execute log command
		cmd := exec.Command("bash", "-c", strCmd)
		var outb bytes.Buffer
		cmd.Stdout = &outb

		err := cmd.Run()
		if err != nil {
			return err
		}

		// If ccapi is listening, finalize execution
		if outb.String() == "Listening on port 80\n" {
			wait = false
		} else {
			time.Sleep(time.Second)
		}
	}

	return nil
}
