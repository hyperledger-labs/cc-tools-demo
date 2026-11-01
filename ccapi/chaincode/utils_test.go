package chaincode

import (
	"errors"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/multi"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	pkgerrors "github.com/pkg/errors"
)

func TestParseSDKError(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{
			"chaincode status",
			status.New(status.ChaincodeStatus, 400, "limit must be greater than 0", nil),
			400, "limit must be greater than 0",
		},
		{
			"multiple endorsers",
			multi.New(
				pkgerrors.Wrapf(errors.New("connection refused"), "Transaction processing for endorser [%s]", "peer0:7051"),
				pkgerrors.Wrapf(status.New(status.ChaincodeStatus, 403, "current caller not allowed", nil), "Transaction processing for endorser [%s]", "peer1:7051"),
			),
			403, "current caller not allowed",
		},
		{
			"wrapped chaincode status",
			pkgerrors.Wrap(status.New(status.ChaincodeStatus, 404, "asset not found", nil), "query failed"),
			404, "asset not found",
		},
		{
			"endorsement validation status",
			status.New(status.EndorserServerStatus, 409, "asset already exists", nil),
			409, "asset already exists",
		},
		{
			"non chaincode status",
			status.New(status.EndorserClientStatus, int32(status.ConnectionFailed), "connection failed", nil),
			500, "Endorser Client Status Code: (2) CONNECTION_FAILED. Description: connection failed",
		},
		{
			"no status",
			errors.New("connection refused"),
			500, "connection refused",
		},
	}

	for _, c := range cases {
		st, err := parseSDKError(c.err)
		if st != c.status || err.Error() != c.message {
			t.Errorf("%s: expected %d %q, got %d %q", c.name, c.status, c.message, st, err.Error())
		}
	}
}
