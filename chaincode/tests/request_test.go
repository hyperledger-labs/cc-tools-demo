package tests
import (
	"github.com/cucumber/godog"
)
func iMakeARequestToWith(arg1, arg2 string, arg3 *godog.DocString) error {
	return godog.ErrPending
}

func theResponseCodeShouldBe(arg1 int) error {
	return godog.ErrPending
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