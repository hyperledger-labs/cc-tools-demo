#!/bin/bash

source scripts/utils.sh

CHANNEL_NAME=${1:-"mychannel"}
CC_NAME=${2}
CC_SRC_PATH=${3}
CCAAS_DOCKER_RUN=${4:-"true"}
CC_VERSION=${5:-"1.0"}
CC_SEQUENCE=${6:-"1"}
CC_INIT_FCN=${7:-"NA"}
CC_END_POLICY=${8:-"NA"}
CC_COLL_CONFIG=${9:-"NA"}
DELAY=${10:-"3"}
MAX_RETRY=${11:-"5"}
VERBOSE=${12:-"false"}
ORG_QNTY=${13:-3}
CCAAS_TLS_ENABLED=${14:-"false"}

CCAAS_SERVER_PORT=9999

: ${CONTAINER_CLI:="docker"}
: ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- CCAAS_DOCKER_RUN: ${C_GREEN}${CCAAS_DOCKER_RUN}${C_RESET}"
println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"
println "- ORG_QUANTITY: ${C_GREEN}${ORG_QNTY}${C_RESET}"
println "- CCAAS_TLS_ENABLED: ${C_GREEN}${CCAAS_TLS_ENABLED}${C_RESET}"

FABRIC_CFG_PATH=$PWD/./config/

#User has not provided a name
if [ -z "$CC_NAME" ] || [ "$CC_NAME" = "NA" ]; then
  fatalln "No chaincode name was provided. Valid call example: ./network.sh deployCCAS -ccn basic -ccp ../asset-transfer-basic/chaincode-go "

# User has not provided a path
elif [ -z "$CC_SRC_PATH" ] || [ "$CC_SRC_PATH" = "NA" ]; then
  fatalln "No chaincode path was provided. Valid call example: ./network.sh deployCCAS -ccn basic -ccp ../asset-transfer-basic/chaincode-go "

## Make sure that the path to the chaincode exists
elif [ ! -d "$CC_SRC_PATH" ]; then
  fatalln "Path to chaincode does not exist. Please provide different path."
fi

if [ "$CC_END_POLICY" = "NA" ]; then
  CC_END_POLICY=""
else
  CC_END_POLICY="--signature-policy $CC_END_POLICY"
fi

if [ "$CC_COLL_CONFIG" = "NA" ]; then
  CC_COLL_CONFIG=""
else
  CC_COLL_CONFIG="--collections-config $CC_COLL_CONFIG"
fi

# import utils
. scripts/envVar.sh
. scripts/ccutils.sh

packageChaincode() {
  address="${CC_NAME}.{{.org}}.example.com:{{.port}}"
  prefix=$(basename "$0")
  tempdir=$(mktemp -d -t "$prefix.XXXXXXXX") || error_exit "Error creating temporary directory"
  label=${CC_NAME}_${CC_VERSION}
  mkdir -p "$tempdir/src"

cat > "$tempdir/src/connection.json" <<CONN_EOF
{
  "address": "${address}",
  "dial_timeout": "10s",
  "tls_required": $CCAAS_TLS_ENABLED,
  "client_auth_required": $CCAAS_TLS_ENABLED
}
CONN_EOF

  cat "$tempdir/src/connection.json"

   mkdir -p "$tempdir/pkg"

cat << METADATA-EOF > "$tempdir/pkg/metadata.json"
{
    "type": "external",
    "label": "$label"
}
METADATA-EOF

    tar -C "$tempdir/src" -czf "$tempdir/pkg/code.tar.gz" .
    tar -C "$tempdir/pkg" -czf "${CC_NAME}.tar.gz" metadata.json code.tar.gz
    rm -Rf "$tempdir"

    export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)
  
    successln "Chaincode is packaged  ${address}, id ${PACKAGE_ID}"
}

# Build the docker image 
buildDockerImages

# Package chaincode
packageChaincode

# start the container
startDockerContainer

if [ $ORG_QNTY -gt 1 ] ; then
  ## Install chaincode on peer0.org1 and peer0.org2
  infoln "Installing chaincode on peer0.org1..."
  installChaincode 1

  infoln "Install chaincode on peer0.org2..."
  installChaincode 2

  infoln "Install chaincode on peer0.org3..."
  installChaincode 3

  ## query whether the chaincode is installed
  queryInstalled 1
  queryInstalled 2
  queryInstalled 3

  ## approve the definition for org1
  approveForMyOrg 1

  ## check whether the chaincode definition is ready to be committed
  ## expect org1 to have approved and org2 not to
  checkCommitReadiness 1 "\"org1MSP\": true" "\"org2MSP\": false" "\"org3MSP\": false"
  checkCommitReadiness 2 "\"org1MSP\": true" "\"org2MSP\": false" "\"org3MSP\": false"
  checkCommitReadiness 3 "\"org1MSP\": true" "\"org2MSP\": false" "\"org3MSP\": false"

  ## now approve also for org2
  approveForMyOrg 2

  ## check whether the chaincode definition is ready to be committed
  ## expect them both to have approved
  checkCommitReadiness 1 "\"org1MSP\": true" "\"org2MSP\": true" "\"org3MSP\": false"
  checkCommitReadiness 2 "\"org1MSP\": true" "\"org2MSP\": true" "\"org3MSP\": false"
  checkCommitReadiness 3 "\"org1MSP\": true" "\"org2MSP\": true" "\"org3MSP\": false"

  ## now approve also for org3 
  approveForMyOrg 3

  ## check whether the chaincode definition is ready to be committed
  ## expect them both to have approved
  checkCommitReadiness 1 "\"org1MSP\": true" "\"org2MSP\": true" "\"org3MSP\": true"
  checkCommitReadiness 2 "\"org1MSP\": true" "\"org2MSP\": true" "\"org3MSP\": true"
  checkCommitReadiness 3 "\"org1MSP\": true" "\"org2MSP\": true" "\"org3MSP\": true"

  ## now that we know for sure both orgs have approved, commit the definition
  commitChaincodeDefinition 1 2 3

  ## query on both orgs to see that the definition committed successfully
  queryCommitted 1
  queryCommitted 2
  queryCommitted 3

  ## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
  ## method defined
  if [ "$CC_INIT_FCN" = "NA" ]; then
    infoln "Chaincode initialization is not required"
  else
    chaincodeInvokeInit 1 2 3
  fi
else
  ## Install chaincode on peer0.org
  infoln "Installing chaincode on peer0.org..."
  installChaincode 0

  ## query whether the chaincode is installed
  queryInstalled 0

  ## approve the definition for org
  approveForMyOrg 0

  ## check whether the chaincode definition is ready to be committed
  checkCommitReadiness 0 "\"orgMSP\": true"

  ## now that we know for sure org have approved, commit the definition
  commitChaincodeDefinition 0

  ## query on both orgs to see that the definition committed successfully
  queryCommitted 0

  ## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
  ## method defined
  if [ "$CC_INIT_FCN" = "NA" ]; then
    infoln "Chaincode initialization is not required"
  else
    chaincodeInvokeInit 0
  fi
fi

exit 0
