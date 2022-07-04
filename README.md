# Goledger CC Tools Demo Chaincode 

## Directory Structure

- `/fabric`: Fabric network used as a test environment
- `/chaincode`: chaincode-related files
- `/rest-server`: chaincode REST API project

## Development

Dependencies for chaincode:

- Go 1.14 or higher

Dependencies for test environment:

- Docker 20.10.5 or higher
- Docker Compose 1.28.5 or higher
- Node v8.17.0
- Node Package Manager (npm) 6.14.11 or higher

Installation:

```bash
$ cd chaincode; go mod vendor; cd ..
$ cd rest-server; npm install; cd ..
```

After installing, use the script `./startDev.sh` in the root folder to start the development environment. It will
start all components of the project with 3 organizations.
To start the development network with only 1 organization, run `$ ./startDev.sh -n 1`.

To apply chaincode changes, run `$ ./upgradeCC.sh <version>` with a version higher than the current one (starts with 0.1).
To upgrade a chaincode with only one organization, run `$ ./upgradeCC.sh <version> -n 1`

To apply CC API changes, run `$ ./reloadCCAPI.sh`.

To test transactions after starting all components, run `$ ./tryout.sh`.

To test transactions using the godog tool, run `$ ./godog.sh`.

More documentation and details can be found at [https://goledger-cc-tools.readthedocs.io/en/latest/](https://goledger-cc-tools.readthedocs.io/en/latest/)

For production deployment please consider using GoFabric - [https://gofabric.io](https://gofabric.io)