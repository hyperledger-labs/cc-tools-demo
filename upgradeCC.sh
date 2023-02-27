#!/usr/bin/env bash

if [[ $# -lt 2 || $# -gt 4 ]]; then
    printf 'Usage: ./upgradeCC.sh <version> <sequence> [-n <org-quantity>]\n'
    exit 1
fi

ORG_QNTY=3

if [[ $# -eq 4 && "$3" == "-n" ]]; then
    ORG_QNTY=$4
fi

if [ $ORG_QNTY != 3 -a $ORG_QNTY != 1 ]
then
  echo 'WARNING: The number of organizations allowed is either 3 or 1.'
  echo 'Defaulting to 3 organizations.'
  ORG_QNTY=3
fi

if [ $ORG_QNTY == 1 ]
then
  CCCG_PATH="../chaincode/collections2-org.json"
else
  CCCG_PATH="../chaincode/collections2.json"
fi

cd ./chaincode; go fmt ./...; cd ..

version=$1
sequence=$2

cd fabric
./network.sh deployCC -ccn cc-tools-demo -ccp ../chaincode -ccl go -ccv $version -ccs $sequence -n $ORG_QNTY -cccg $CCCG_PATH -cci init
cd ..
