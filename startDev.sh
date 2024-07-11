#!/usr/bin/env bash

ORG_QNTY=3
DEPLOY_CCAAS=false
CCAAS_TLS_ENABLEd=""
SKIP_COLL_GEN=false

while [[ $# -ge 1 ]] ; do
    key="$1"
    case $key in
        -n )
            ORG_QNTY=$2
            shift
            ;;
        -ccaas )
            DEPLOY_CCAAS=$2
            shift
            ;;
        -ccaastls )
            CCAAS_TLS_ENABLED="-ccaastls"
            shift
            ;;
        -c )
            SKIP_COLL_GEN=true
            ;;
  esac
  shift
done

if [ $ORG_QNTY != 3 -a $ORG_QNTY != 1 ]
then
  echo 'WARNING: The number of organizations allowed is either 3 or 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QNTY=3
fi

if [ "$SKIP_COLL_GEN" = false ] ; then
    echo 'Generating collections configuration file...'
    if [ $ORG_QNTY == 1 ]
    then
        cd ./chaincode; go run . -g --orgs orgMSP; cd ..
    else
        cd ./chaincode; go run . -g --orgs org1MSP org2MSP org3MSP; cd ..
    fi
fi

# Clear unused images and volumes
docker rmi $(docker images --quiet --filter "dangling=true")
docker volume rm $(docker volume ls -qf dangling=true)

# Script used to start the development environment.
if [ ! -d "chaincode/vendor" ]; then
    cd ./chaincode; GOWORK=off go mod vendor; cd ..
fi
cd ./chaincode; go fmt ./...; cd ..
cd ./fabric; ./startDev.sh -n $ORG_QNTY -ccaas $DEPLOY_CCAAS $CCAAS_TLS_ENABLED; cd ..

## This brings up API in Go
if [ $ORG_QNTY == 1 ]
then
    cd ./ccapi; docker-compose -f docker-compose-1org.yaml up -d; cd ..
else
    cd ./ccapi; docker-compose up -d; cd ..
fi