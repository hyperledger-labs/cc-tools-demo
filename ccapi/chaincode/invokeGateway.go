package chaincode

import (
	"os"

	"github.com/hyperledger-labs/cc-tools-demo/ccapi/common"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	cb "github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// InvokeGateway submits the transaction and returns its result and the success status set by the chaincode.
func InvokeGateway(channelName, chaincodeName, txName, user string, args []string, transientArgs []byte, endorsingOrgs []string) ([]byte, int, error) {
	// Gateway endpoint
	endpoint := os.Getenv("FABRIC_GATEWAY_ENDPOINT")

	// Create client grpc connection
	grpcConn, err := common.CreateGrpcConnection(endpoint)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create grpc connection")
	}
	defer grpcConn.Close()

	// Create gateway connection
	gw, err := common.CreateGatewayConnection(grpcConn, user)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create gateway connection")
	}
	defer gw.Close()

	// Obtain smart contract deployed on the network.
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Make transient request
	transientMap := make(map[string][]byte)
	transientMap["@request"] = transientArgs

	options := []client.ProposalOption{client.WithArguments(args...)}
	if transientArgs != nil {
		options = append(options, client.WithTransient(transientMap))
	}
	if len(endorsingOrgs) > 0 {
		options = append(options, client.WithEndorsingOrganizations(endorsingOrgs...))
	}

	// Invoke transaction
	return submit(contract, txName, options...)
}

// submit follows the same steps as client.Contract.Submit, also reading the chaincode status
// from the endorsed transaction.
func submit(contract *client.Contract, txName string, options ...client.ProposalOption) ([]byte, int, error) {
	proposal, err := contract.NewProposal(txName, options...)
	if err != nil {
		return nil, 0, err
	}

	transaction, err := proposal.Endorse()
	if err != nil {
		return nil, 0, err
	}

	transactionBytes, err := transaction.Bytes()
	if err != nil {
		return nil, 0, err
	}

	chaincodeStatus, err := chaincodeStatusFromPreparedTransaction(transactionBytes)
	if err != nil {
		return nil, 0, err
	}

	commit, err := transaction.Submit()
	if err != nil {
		return nil, 0, err
	}

	commitStatus, err := commit.Status()
	if err != nil {
		return nil, 0, err
	}

	if !commitStatus.Successful {
		return nil, 0, &client.CommitError{
			TransactionID: commitStatus.TransactionID,
			Code:          commitStatus.Code,
		}
	}

	return transaction.Result(), successStatus(chaincodeStatus), nil
}

// chaincodeStatusFromPreparedTransaction reads the chaincode response status from a serialized
// gateway.PreparedTransaction, as returned by client.Transaction.Bytes.
func chaincodeStatusFromPreparedTransaction(preparedTransactionBytes []byte) (int32, error) {
	preparedTransaction := &gateway.PreparedTransaction{}
	if err := proto.Unmarshal(preparedTransactionBytes, preparedTransaction); err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal prepared transaction")
	}

	payload := &cb.Payload{}
	if err := proto.Unmarshal(preparedTransaction.GetEnvelope().GetPayload(), payload); err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal payload")
	}

	transaction := &peer.Transaction{}
	if err := proto.Unmarshal(payload.GetData(), transaction); err != nil {
		return 0, errors.Wrap(err, "failed to unmarshal transaction")
	}

	for _, action := range transaction.GetActions() {
		actionPayload := &peer.ChaincodeActionPayload{}
		if err := proto.Unmarshal(action.GetPayload(), actionPayload); err != nil {
			return 0, errors.Wrap(err, "failed to unmarshal chaincode action payload")
		}

		responsePayload := &peer.ProposalResponsePayload{}
		if err := proto.Unmarshal(actionPayload.GetAction().GetProposalResponsePayload(), responsePayload); err != nil {
			return 0, errors.Wrap(err, "failed to unmarshal proposal response payload")
		}

		chaincodeAction := &peer.ChaincodeAction{}
		if err := proto.Unmarshal(responsePayload.GetExtension(), chaincodeAction); err != nil {
			return 0, errors.Wrap(err, "failed to unmarshal chaincode action")
		}

		return chaincodeAction.GetResponse().GetStatus(), nil
	}

	return 0, errors.New("no chaincode action found in transaction")
}
