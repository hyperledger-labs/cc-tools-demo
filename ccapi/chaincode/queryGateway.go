package chaincode

import (
	"context"
	"os"

	"github.com/hyperledger-labs/cc-tools-demo/ccapi/common"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// QueryGateway evaluates the transaction and returns its result and the success status set by the chaincode.
func QueryGateway(channelName, chaincodeName, txName, user string, args []string) ([]byte, int, error) {
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

	// Query transaction
	proposal, err := contract.NewProposal(txName, client.WithArguments(args...))
	if err != nil {
		return nil, 0, err
	}

	sign, err := common.NewGatewaySign(user)
	if err != nil {
		return nil, 0, err
	}

	return evaluate(grpcConn, channelName, sign, proposal)
}

// evaluate sends the proposal straight to the Gateway service, since client.Proposal.Evaluate
// returns only the payload and discards the status set by the chaincode.
func evaluate(conn grpc.ClientConnInterface, channelName string, sign identity.Sign, proposal *client.Proposal) ([]byte, int, error) {
	signature, err := sign(proposal.Digest())
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to sign proposal")
	}

	proposalBytes, err := proposal.Bytes()
	if err != nil {
		return nil, 0, err
	}

	proposedTransaction := &gateway.ProposedTransaction{}
	if err := proto.Unmarshal(proposalBytes, proposedTransaction); err != nil {
		return nil, 0, errors.Wrap(err, "failed to unmarshal proposal")
	}
	proposedTransaction.GetProposal().Signature = signature

	ctx, cancel := context.WithTimeout(context.Background(), common.GatewayEvaluateTimeout)
	defer cancel()

	response, err := gateway.NewGatewayClient(conn).Evaluate(ctx, &gateway.EvaluateRequest{
		TransactionId:       proposedTransaction.GetTransactionId(),
		ChannelId:           channelName,
		ProposedTransaction: proposedTransaction.GetProposal(),
		TargetOrganizations: proposedTransaction.GetEndorsingOrganizations(),
	})
	if err != nil {
		return nil, 0, err
	}

	return response.GetResult().GetPayload(), successStatus(response.GetResult().GetStatus()), nil
}
