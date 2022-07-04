#!/usr/bin/env bash

# set -x

generate() { 
  if [ "$1" == "1" ]; then
    orgs=(org.example.com)
    mkdir -p ca/org
  else
    orgs=(org1.example.com org2.example.com org3.example.com)
    mkdir -p ca/org1 ca/org2 ca/org3
  fi

  sudo rm -rf crypto-config
  mkdir -p crypto-config
  sudo rm -rf ca/org*/fabric-ca-server.db
  
  source scripts/generateCerts.sh
  generateCerts $1
  if [[ $? -ne 0 ]]; then
    echo "Error at cert generation"
    exit 1
  fi

  # Create rest-certs folder
  REST_CERTS_FOLDER=./crypto-config/rest-certs

  for ORGNAME in ${orgs[@]}; do
    # Copy org rest-certs
    PRIVATE_KEY_PATH=$(ls $PWD/crypto-config/peerOrganizations/${ORGNAME}/users/Admin.${ORGNAME}/msp/keystore/*_sk)
    PRIVATE_KEY=$(basename $PRIVATE_KEY_PATH)
    echo $PRIVATEKEY

    mkdir -p $REST_CERTS_FOLDER/${ORGNAME}
    # Copy admin privkey to rest-certs
    cp $PRIVATE_KEY_PATH $REST_CERTS_FOLDER/${ORGNAME}/admin.key
    # Copy admin cert to rest-certs
    ADMINCERT_PATH=$PWD/crypto-config/peerOrganizations/${ORGNAME}/users/Admin.${ORGNAME}/msp/signcerts/Admin.${ORGNAME}-cert.pem
    cp $ADMINCERT_PATH $REST_CERTS_FOLDER/${ORGNAME}/admin.cert
    # Copy cacert to rest-certs
    CACERT_PATH=$PWD/crypto-config/peerOrganizations/${ORGNAME}/peers/peer0.${ORGNAME}/tls/ca.crt
    cp $CACERT_PATH $REST_CERTS_FOLDER/${ORGNAME}/ca.crt
  done
}

start_network() {
  CHAINCODE_NAME=cc-tools-demo
  ./clearDev.sh
  docker network create ${CHAINCODE_NAME}-net
  # docker-compose -f docker-compose-ca.yaml up -d

  if [ "$1" == "1" ]; then
    orgs=(org)
    ports=(3000)

    docker-compose -f docker-compose-org.yaml pull
    docker-compose -f docker-compose-org.yaml up -d
  else
    orgs=(org1 org2 org3)
    ports=(3000 3001 3002)

    docker-compose pull
    docker-compose up -d
  fi
  sleep 30

  printf '\n\nCreate channel - mainchannel\n'
  curl -k --request POST \
  --url 'https://localhost:3000/api/v1/network/channel/create?channelName=mainchannel' \
  --header 'gofabricversion: 0.9.0' \
  -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
  --form channeltx=@./channel-artifacts/channel.tx

  sleep 10 # wait for channel RAFT to be available
  
  # Join peers to channel
  n=${#orgs[@]}
  for i in $(seq 0 $((n-1))); do
    ORGNAME=${orgs[i]}
    PORTNUMBER=${ports[i]}

    printf "\n\nJoin ${ORGNAME} to channel\n"
    curl -k -X POST \
      "https://localhost:${PORTNUMBER}/api/v1/network/channel/join" \
      -H 'Content-Type: application/json' \
      --header 'gofabricversion: 0.9.0' \
      -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
      -d "{ \"channelName\": \"mainchannel\", 
      \"targets\": [
      {
        \"name\": \"peer0.${ORGNAME}.example.com\",
        \"ip\": \"peer0.${ORGNAME}.example.com\"
      }
    ]
    }"
  done

  # Update anchor peers
  for i in $(seq 0 $((n-1))); do
    ORGNAME=${orgs[i]}
    PORTNUMBER=${ports[i]}
    
    printf "\n\nUpdate anchor peers on ${ORGNAME}\n"
    curl -k -X POST \
    "https://localhost:${PORTNUMBER}/api/v1/network/channel/anchorPeers?channelName=mainchannel" \
    --header 'gofabricversion: 0.9.0' \
    -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
    --form anchorstx=@./channel-artifacts/${ORGNAME}MSPanchors.tx
  done

  # Install chaincode
  for i in $(seq 0 $((n-1))); do
    ORGNAME=${orgs[i]}
    PORTNUMBER=${ports[i]}
    
    printf "\n\nInstall chaincode on ${ORGNAME}\n"
    curl -sS -k -X POST \
    "https://localhost:${PORTNUMBER}/api/v1/network/channel/chaincode/install" \
    -H 'Content-Type: application/json' \
    --header 'gofabricversion: 0.9.0' \
    -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
    -d '{
          "chaincode":        "cc-tools-demo",
          "channelName":      "mainchannel",
          "chaincodeVersion": "0.1"
        }'
  done

  if [ "$1" == "1" ]; then
    instantiateChaincodeData="{
          \"channelName\": \"mainchannel\",
          \"chaincode\": \"cc-tools-demo\",
          \"chaincodeVersion\": \"0.1\",
          \"endorsement\": {
              \"identities\": [
                  {
                      \"role\": {
                          \"name\": \"member\",
                          \"mspId\": \"orgMSP\"
                      }
                  }
              ],
              \"policy\": {
                  \"1-of\": [
                      {
                          \"signed-by\": 0
                      }
                  ]
              }
          }
      }"
  else
    instantiateChaincodeData="{
          \"channelName\": \"mainchannel\",
          \"chaincode\": \"cc-tools-demo\",
          \"chaincodeVersion\": \"0.1\",
          \"endorsement\": {
              \"identities\": [
                  {
                      \"role\": {
                          \"name\": \"member\",
                          \"mspId\": \"org1MSP\"
                      }
                  },
                  {
                      \"role\": {
                          \"name\": \"member\",
                          \"mspId\": \"org2MSP\"
                      }
                  },
                  {
                      \"role\": {
                          \"name\": \"member\",
                          \"mspId\": \"org3MSP\"
                      }
                  }
              ],
              \"policy\": {
                  \"1-of\": [
                      {
                          \"signed-by\": 1
                      },
                      {
                          \"signed-by\": 2
                      }
                  ]
              }
          }
      }"
  fi  

  printf '\n\nInstantiate chaincode\n'
  curl -k -X POST \
    https://localhost:3000/api/v1/network/channel/chaincode/instantiate \
    -H 'Content-Type: application/json' \
    --header 'gofabricversion: 0.9.0' \
    -H 'magicnumber: dfff482c-1df5-42ad-95d4-d8d72b2398be' \
    -d "$instantiateChaincodeData"

  printf '\n\n'
}

create_artifacts() {
  if [ "$1" == "1" ]; then
    orgs=(org)
    CONFIGPATH="./configtx-1org"
  else
    orgs=(org1 org2 org3)
    CONFIGPATH="./configtx-3org"
  fi

  rm -rf channel-artifacts
  mkdir channel-artifacts

  export FABRIC_CFG_PATH=$PWD
  configtxgen -configPath $CONFIGPATH -profile SingleNodeEtcdRaft -channelID sys-channel -outputBlock ./channel-artifacts/genesis.block
  export CHANNEL_NAME=mainchannel && configtxgen -configPath $CONFIGPATH -profile OrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
  
  for org in ${orgs[@]}; do
    configtxgen -configPath $CONFIGPATH -profile OrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/${org}MSPanchors.tx -channelID $CHANNEL_NAME -asOrg ${org}MSP
  done
}

export IMAGE_TAG=1.4
export COMPOSE_PROJECT_NAME=fabric
if [ "$OSTYPE" == "linux-gnu" ] ; then 
  export PATH=$PWD/bin:$PATH
fi

###############################################
######### Script starts here ##################
###############################################

# if [ "$1" == "up" ] ||  [ "$1" == "generate" ]; then
#   COMMAND=$1
#   shift
# fi

while test $# -gt 0; do
  case "$1" in
      -n)
        shift
        NUMORGS=$1
        shift
        ;;
      generate)
        COMMAND=$1
        shift
        ;;
      up)
        COMMAND=$1
        shift
        ;;
      *)
        errorln "Unknown command or flag: $1"
        exit 1
        ;;
  esac
done

case "$COMMAND" in
  generate)
    generate $NUMORGS
    create_artifacts $NUMORGS
    exit 1
    ;;
  up)
    start_network $NUMORGS
    exit 1
    ;;
  *)
    # if [ ! -f "channel-artifacts/channel.tx" ]; then
    #   echo "Certs not found. Generating certs."
    generate $NUMORGS
    create_artifacts $NUMORGS
    # fi
    start_network $NUMORGS
    ;;
esac
