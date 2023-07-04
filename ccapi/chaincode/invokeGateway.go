package chaincode

import (
	"os"

	"github.com/goledgerdev/ccapi/common"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/pkg/errors"
)

func InvokeGateway(channelName, chaincodeName, txName, args string, transientArgs map[string][]byte, endorsingOrgs []string) ([]byte, error) {
	// Gateway endpoint
	endpoint := os.Getenv("FABRIC_GATEWAY_ENDPOINT")

	// Create client grpc connection
	grpcConn, err := common.CreateGrpcConnection(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc connection")
	}
	defer grpcConn.Close()

	// Create gateway connection
	gw, err := common.CreateGatewayConnection(grpcConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gateway connection")
	}
	defer gw.Close()

	// Obtain smart contract deployed on the network.
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Invoke transaction
	if transientArgs != nil {
		return contract.Submit(txName,
			client.WithArguments(args),
			client.WithTransient(transientArgs),
			client.WithEndorsingOrganizations(endorsingOrgs...),
		)
	}

	return contract.SubmitTransaction(txName, args)
}
