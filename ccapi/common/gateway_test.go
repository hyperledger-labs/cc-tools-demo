package common

import (
	"testing"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func gatewayError(code codes.Code, msg string, details ...string) error {
	st := status.New(code, msg)
	for _, d := range details {
		st, _ = st.WithDetails(&gateway.ErrorDetail{Address: "peer0:7051", MspId: "orgMSP", Message: d})
	}
	return st.Err()
}

func TestParseError(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{
			"evaluate chaincode 400",
			gatewayError(codes.Unknown, "evaluate call to endorser returned error: chaincode response 400, limit must be greater than 0", "chaincode response 400, limit must be greater than 0"),
			400, "limit must be greater than 0",
		},
		{
			"chaincode 500 without prefix",
			gatewayError(codes.Aborted, "failed to evaluate transaction, see attached details for more info", "error searching for book's author: boom", "error searching for book's author: boom"),
			500, "error searching for book's author: boom",
		},
		{
			"chaincode status after unreachable peer",
			gatewayError(codes.Aborted, "failed to collect enough transaction endorsements, see attached details for more info", "connection refused", "chaincode response 409, conflict"),
			409, "conflict",
		},
		{
			"no chaincode status",
			gatewayError(codes.Aborted, "failed to endorse transaction, see attached details for more info", "connection refused"),
			500, "connection refused",
		},
		{
			"no details",
			gatewayError(codes.DeadlineExceeded, "evaluate timeout expired"),
			500, "unexpected error type:rpc error: code = DeadlineExceeded desc = evaluate timeout expired",
		},
		{
			"commit error",
			&client.CommitError{TransactionID: "tx1", Code: peer.TxValidationCode_MVCC_READ_CONFLICT},
			500, "transaction failed to commit with status code 11 (MVCC_READ_CONFLICT)",
		},
	}

	for _, c := range cases {
		err, st := ParseError(c.err)
		if st != c.status || err.Error() != c.message {
			t.Errorf("%s: expected %d %q, got %d %q", c.name, c.status, c.message, st, err.Error())
		}
	}
}
