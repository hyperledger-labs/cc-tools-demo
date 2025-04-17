# Hyperledger Labs CC Tools Demo Chaincode 

## Directory Structure

- `/fabric`: Fabric network v2.5 used as a test environment
- `/chaincode`: chaincode-related files
- `/ccapi`: chaincode REST API in Golang project
- `/fabric-private-chaincode`: Explaining the integration project between CC-tools and FPC

## Development

The `cc-tools` library has been tested in Fabric v2.2, v2.4 and v2.5 networks.

Dependencies for chaincode and chaincode API:

- Go 1.21 or higher

Dependencies for test environment:

- Docker 20.10.5 or higher
- Docker Compose 1.28.5 or higher

Intallation of the Chaincode API Go:

```bash
$ cd chaincode; go mod vendor; cd ..
$ cd ccapi; go mod vendor; cd ..
```


## Deploying test environment

After installing, use the script `./startDev.sh` in the root folder to start the development environment. It will
start all components of the project with 3 organizations.

If you want to deploy with 1 organization, run the command `./startDev.sh -n 1`.

To apply chaincode changes, run `$ ./upgradeCC.sh <version> <sequence>` with a version higher than the current one (starts with 0.1). Append `-n 1` to the command if running with 1 organization.

To apply CC API changes, run `$ ./scripts/reloadCCAPI.sh`.

## Deploying Chaincode as a service

After installing, use the script `./startDev.sh -ccaas` in the root folder to start the development environment. It will
start all components of the project with 3 organizations.

If you want to deploy with 1 organization, run the command `./startDev.sh -ccaas -n 1`.

To apply chaincode changes, run `$ ./upgradeCC.sh -ccaas <version> <sequence>` with a version higher than the current one (starts with 0.1). Append `-n 1` to the command if running with 1 organization.

To apply CC API changes, run `$ ./scripts/reloadCCAPI.sh`.

## Automated tryout and test

To test transactions after starting all components, run `$ ./scripts/tryout.sh`. 

To test transactions using the godog tool, run `$ ./scripts/godog.sh`.


## Generate TAR archive for the chaincode

The `generatePackage.sh` script is available to generate a `tar.gz` archive of the chaincode. 

By running `$ ./generatePackage.sh` without any option, the script generates a `collections.json` file for the private data on the chaincode with all the organizations defined on the readers section of private asset types, and then archives the code without the CCAPI.

By using the `--org/-o` option along the script, it's possible to specify the organizations to be considered when generating the `collections.json` file. This option may be used multiple times to add all the organizations, ex: `$ ./generatePackage.sh -o org1MSP -o org2MSP`.

By standard the archive is created using the project name with *1.0* label, to change it the `--name/-n` and `--label/-l` flags may be used. Example: `$ ./generatePackage.sh -n my-project -l 2.0`

## Integration with Fabric Private Chaincode

If you want to execute your chaincode in a Trusted Execution Environment (TEE) using Fabric Private Chaincode (FPC), we've set up an integration guide to help you. Check out the instructions in the `./fabric-private-chaincode` directory to seamlessly integrate FPC with CC Tools for enhanced privacy and security.

## More

You can reach GoLedger developers and `cc-tools` maintainers at our Discord - [Join us!](https://discord.gg/GndkYHxNyQ)

More documentation and details on `cc-tools` can be found at [https://goledger-cc-tools.readthedocs.io/en/latest/](https://goledger-cc-tools.readthedocs.io/en/latest/)
