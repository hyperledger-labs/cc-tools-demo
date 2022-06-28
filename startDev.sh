#!/usr/bin/env bash

# Clear unused images and volumes
docker rmi $(docker images --quiet --filter "dangling=true")
docker volume rm $(docker volume ls -qf dangling=true)

# Script used to start the development environment.
if [ ! -d "chaincode/vendor" ]; then
    cd ./chaincode; go mod vendor; cd ..
fi

while test $# -gt 0; do
  case "$1" in
      -n)
        shift
        NUMORGS=$1
        shift
        ;;
      *)
        errorln "Unknown flag: $1"
        exit(1)
        ;;
  esac
done

cd ./chaincode; go fmt ./...; cd ..
cd ./fabric; ./startDev.sh -n ${NUMORGS}; cd ..
cd ./rest-server; ./startDev.sh -n ${NUMORGS}; cd ..
