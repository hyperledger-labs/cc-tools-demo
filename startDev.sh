#!/usr/bin/env bash

ORG_QNTY=3
SKIP_COLL_GEN=false

while getopts n:c opt; do
    case $opt in
        n)  ORG_QNTY=${OPTARG}
            ;;
        c)  SKIP_COLL_GEN=true
            ;;
    esac
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
        cd ./chaincode; go run . -g --orgs-1; cd ..
    else
        cd ./chaincode; go run . -g --orgs-3; cd ..
    fi
fi

# Clear unused images and volumes
docker rmi $(docker images --quiet --filter "dangling=true")
docker volume rm $(docker volume ls -qf dangling=true)

# Script used to start the development environment.
if [ ! -d "chaincode/vendor" ]; then
    cd ./chaincode; go mod vendor; cd ..
fi
cd ./chaincode; go fmt ./...; cd ..
cd ./fabric; ./startDev.sh -n $ORG_QNTY; cd ..

## This brings up API in Node
# cd ./rest-server; ./startDev.sh -n $ORG_QNTY; cd ..

## This brings up API in Go
if [ $ORG_QNTY == 1 ]
then
    cd ./ccapi; docker-compose -f docker-compose-1org.yaml up -d; cd ..
else
    cd ./ccapi; docker-compose up -d; cd ..
fi