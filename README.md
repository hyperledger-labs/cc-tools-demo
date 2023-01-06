# Goledger CC Tools Demo Chaincode 

## Directory Structure

- `/fabric2`: Fabric network v2.2 used as a test environment
- `/chaincode`: chaincode-related files
- `/ccapi`: chaincode REST API in Golang project

The following directories are deprecated:
- `/fabric`: Fabric network v1.4 used as a test environment (deprecated)
- `/rest-server`: chaincode REST API in Node.js project (deprecated)

## Development

The `cc-tools` library has been tested in Fabric v1.4, v2.2 and v2.4 networks.

Dependencies for chaincode and chaincode API:

- Go 1.14 or higher

Dependencies for test environment:

- Docker 20.10.5 or higher
- Docker Compose 1.28.5 or higher
- Node v8.17.0
- Node Package Manager (npm) 6.14.11 or higher

Installation if using Chaincode API Node:

```bash
$ cd chaincode; go mod vendor; cd ..
$ cd rest-server; npm install; cd ..
```

Intallation if using Chaincode API Go:

```bash
$ cd chaincode; go mod vendor; cd ..
$ cd ccapi; go mod vendor; cd ..
```

## Deploying test env in v1.4

After installing, use the script `./startDev.sh` in the root folder to start the development environment. It will
start all components of the project with 3 organizations.
To start the development network with only 1 organization, run `$ ./startDev.sh -n 1`.

To apply chaincode changes, run `$ ./upgradeCC.sh <version>` with a version higher than the current one (starts with 0.1).
To upgrade a chaincode with only one organization, run `$ ./upgradeCC.sh <version> -n 1`

To apply CC API changes, run `$ ./reloadCCAPI.sh`.

## Deploying test env in v2.2

After installing, use the script `./startDev2.sh` in the root folder to start the development environment. It will
start all components of the project with 3 organizations. Deploying with only 1 org will come in a future version.

To apply chaincode changes, run `$ ./upgradeCC2.sh <version> <sequence>` with a version higher than the current one (starts with 0.1).

To apply CC API changes, run `$ ./reloadCCAPI.sh`.

## Observations for Chaincode API

If you want to deploy the Node.js Chaincode API with the scripts that deploy v2.2 you can uncomment `startDev2.sh:L15` and comment `startDev2.sh:L18`. 

```sh
## This brings up API in Node
# cd ./rest-server; ./startDev2.sh; cd ..

## This brings up API in Go
cd ./ccapi; docker-compose up -d; cd ..
```

## Automated tryout and test

To test transactions after starting all components, run `$ ./tryout.sh`.

To test transactions using the godog tool, run `$ ./godog.sh`.

## More

You can reach GoLedger developers and `cc-tools` maintainers at our Discord - [Join us!](https://discord.gg/GndkYHxNyQ)

More documentation and details on `cc-tools` can be found at [https://goledger-cc-tools.readthedocs.io/en/latest/](https://goledger-cc-tools.readthedocs.io/en/latest/)

For production deployment please consider using GoFabric - [https://gofabric.io](https://gofabric.io)