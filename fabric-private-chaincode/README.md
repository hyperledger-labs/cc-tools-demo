# Integration with Fabric Private Chaincode

## Motivation

Fabric Private Chaincode (FPC) enhances the privacy and security of smart contracts on Hyperledger Fabric by leveraging Trusted Execution Environments (TEEs), such as Intel SGX. With FPC, chaincode execution is shielded, ensuring that sensitive data and business logic remain confidential—even from the hosting peers. This makes FPC an ideal choice for organizations handling susceptible information, such as financial transactions or medical records, while still benefiting from the transparency and immutability of blockchain.

By integrating FPC with CC Tools, developers can now enjoy a streamlined workflow for deploying and managing private chaincode. This integration combines the usability of CC Tools with the robust privacy guarantees of FPC, enabling faster development cycles without compromising data security. To learn more about the design and technical details of this integration project, refer to the [design document](https://github.com/hyperledger/fabric-private-chaincode/tree/main/docs/design/integrate-with-cc-tools) for a comprehensive explanation.

## Prerequisites

To be able to use FPC, you are required to have some software components as dependencies to enable the trusted execution of chaincodes. Follow the [getting started guide](https://github.com/hyperledger/fabric-private-chaincode?tab=readme-ov-file#getting-started) on FPC until you can set up the development environment.

## Integrating your chaincode with FPC

As part of the integration project, there is now a tutorial explaining step-by-step how to run your chaincode using cc-tools framework within a trusted execution environment with FPC. Follow the [tutorial](https://github.com/hyperledger/fabric-private-chaincode/tree/main/samples/chaincode/cc-tools-demo) for more.
