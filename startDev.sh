#!/usr/bin/env bash

NUMORGS=3

while test $# -gt 0; do
  case "$1" in
      -n)
        shift
        NUMORGS=$1
        shift
        ;;
      *)
        errorln "Unknown flag: $1"
        exit 1
        ;;
  esac
done

if [ $NUMORGS != 3 -a $NUMORGS != 1 ]
then
  echo 'ERROR: The number of organizations aloowed is either 3 or 1.'
  echo 'The deafult operation assumes 3 organizations.'
  echo "Usage: "
  echo "  startDev.sh [Flags]"
  echo
  echo "    Flags:"
  echo "    -n <number or orgs> - Specify number of orgs in test network (1 or 3 valid)"
  echo ""
  echo " Examples:"
  echo "   startDev.sh -n 1"
  exit 1
fi

# Clear unused images and volumes
docker rmi -f $(docker images --quiet --filter "dangling=true")
docker volume rm -f $(docker volume ls -qf dangling=true)

# Script used to start the development environment.
if [ ! -d "chaincode/vendor" ]; then
    cd ./chaincode; go mod vendor; cd ..
fi

cd ./chaincode; go fmt ./...; cd ..
cd ./fabric; ./startDev.sh -n ${NUMORGS}; cd ..
cd ./rest-server; ./startDev.sh -n ${NUMORGS}; cd ..