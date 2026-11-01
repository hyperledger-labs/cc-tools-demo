package chaincode

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/multi"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/pkg/errors"
)

func extractStatusCode(msg string) int {
	re := regexp.MustCompile(`Code:\s*\((\d+)\)`)

	matches := re.FindStringSubmatch(msg)
	if len(matches) == 0 {
		fmt.Println("No status code found in message")
		return http.StatusInternalServerError
	}

	statusCode, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println("Failed to parse string to int when extracting status code")
		return http.StatusInternalServerError
	}

	if statusCode < 100 || statusCode > 599 {
		return http.StatusInternalServerError
	}

	return statusCode
}

// parseSDKError returns the status and the chaincode message of a fabric-sdk-go error.
// When no endorser returned a chaincode status, the original error is kept.
func parseSDKError(err error) (int, error) {
	errs := []error{err}
	if m, ok := errors.Cause(err).(multi.Errors); ok {
		errs = m
	}

	for _, e := range errs {
		s, ok := errors.Cause(e).(*status.Status)
		if !ok || s.Code < 100 || s.Code > 599 {
			continue
		}
		// Chaincode responses are reported as ChaincodeStatus by the endorser and as
		// EndorserServerStatus by the endorsement validation
		if s.Group == status.ChaincodeStatus || s.Group == status.EndorserServerStatus {
			// Only the chaincode message is returned to the client, so keep the full error in the logs
			fmt.Println("chaincode error: ", err)
			return int(s.Code), errors.New(s.Message)
		}
	}

	return extractStatusCode(err.Error()), err
}
