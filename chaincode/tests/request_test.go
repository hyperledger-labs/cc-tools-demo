package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"reflect"

	"github.com/cucumber/godog"
)

type bodyCtxKey struct{}

func iMakeARequestToWith(arg1, arg2 string, arg3 *godog.DocString) error {
	return godog.ErrPending
}

func theResponseCodeShouldBe(arg1 int) error {
	return godog.ErrPending
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
